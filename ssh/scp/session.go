// Copyright 2023 jim.zoumo@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scp

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/zoumo/golib/fileinfo"
	"github.com/zoumo/golib/log"
)

type uploadFn func(ctx context.Context, fullpath string, info os.FileInfo, content io.Reader) error

type walkFn func(fpath string, info os.FileInfo, content io.Reader) (bool, error)

type sessionMode int

const (
	scpWrite sessionMode = iota
	scpRead
)

func newSession(client *ssh.Client, mode sessionMode, readTimeout time.Duration, logger log.Logger) (*session, error) {
	if logger == nil {
		logger = log.Discard()
	}
	s, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return &session{
		Session:     s,
		readTimeout: readTimeout,
		mode:        mode,
		logger:      logger,
	}, err
}

type session struct {
	*ssh.Session

	reader      *reader
	writer      io.WriteCloser
	readTimeout time.Duration
	mode        sessionMode

	logger log.Logger
}

func (s *session) Close() error {
	s.reader.close()
	return s.Session.Close()
}

func (s *session) scpCmd(location string) string {
	cmd := []string{
		"scp",
		"-p",
	}
	cmd = append(cmd, "-r")
	switch s.mode {
	case scpWrite:
		cmd = append(cmd, "-t")
	case scpRead:
		cmd = append(cmd, "-f")
	}
	cmd = append(cmd, location)
	ret := strings.Join(cmd, " ")
	return ret + "\n"
}

func (s *session) init(location string) error {
	location = strings.TrimSuffix(location, "/")
	var err error
	s.writer, err = s.StdinPipe()
	if err != nil {
		return err
	}
	reader, err := s.StdoutPipe()
	if err != nil {
		return err
	}

	s.reader = newReader(reader, s.readTimeout)
	cmd := s.scpCmd(location)
	s.logger.V(5).Info("sending command", "cmd", cmd)
	err = s.Start(cmd)
	if err != nil {
		return err
	}

	go s.reader.readInBackground()
	return nil
}

func (s *session) write(data []byte) error {
	_, err := s.writer.Write(data)
	return err
}

func (s *session) writeStatusOK() error {
	return s.write([]byte{StatusOK})
}

func (s *session) sendMsgAndReadStatus(msg string) error {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	s.logger.V(5).Info("send scp msg", "msg", strings.TrimSpace(msg))
	err := s.write([]byte(msg))
	if err != nil {
		return err
	}
	return s.reader.readStatus()
}

func (s *session) enterDir(info os.FileInfo) error {
	return s.upload(info, nil)
}

func (s *session) endDir() error {
	return s.sendMsgAndReadStatus(string(EndDirToken))
}

func (s *session) Open(location string) (os.FileInfo, io.Reader, error) {
	var info os.FileInfo
	var reader io.Reader
	err := s.Walk(location, func(fpath string, finfo os.FileInfo, freader io.Reader) (bool, error) {
		if fpath == location {
			info = finfo
			reader = freader
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return info, reader, nil
}

func (s *session) Walk(location string, handler walkFn) error {
	if s.mode != scpRead {
		return fmt.Errorf("session mode must be read in downloader")
	}

	err := s.init(location)
	if err != nil {
		return errors.Wrap(err, "failed to init session")
	}

	now := time.Now()

	curentParentPath := path.Dir(location)
	for {
		continued, err := s.walkFn(&curentParentPath, &now, handler)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if !continued {
			break
		}
	}
	return nil
}

func (s *session) walkFn(parentPath *string, modified *time.Time, handler walkFn) (bool, error) {
	err := s.writeStatusOK()
	if err != nil {
		return false, err
	}
	response, err := s.reader.read()
	if err != nil {
		return false, err
	}
	token := response[0]

	s.logger.V(5).Info("receive scp response", "response", strings.TrimSpace(string(response)))

	switch token {
	case FileToken, DirToken:
		shallContinue, err := s.processResponse(parentPath, response, *modified, handler)
		if err != nil {
			return false, err
		}
		if !shallContinue {
			return false, nil
		}
	case EndDirToken:
		*parentPath = strings.TrimRight(path.Dir(*parentPath), "/")
	case TimestampToken:
		timestamp, err := ParseTimeResponse(string(response))
		if err != nil {
			return false, err
		}
		*modified = *timestamp
	case WarningToken:
		msg := strings.TrimSpace(string(response[1:]))
		if strings.Contains(msg, "No such file or directory") {
			return false, os.ErrNotExist
		}
		fallthrough
	case ErrorToken:
		errorMessage := strings.TrimSpace(string(response[1:]))
		return false, fmt.Errorf("receive error from peer: %s", errorMessage)
	default:
		return false, fmt.Errorf("unsupported token: %v, %s", token, response)
	}
	return true, nil
}

func (s *session) processResponse(parentPath *string, response []byte, modified time.Time, handler walkFn) (bool, error) {
	fileInfo, err := ParseFileInfo(string(response), &modified)
	if err != nil {
		return false, err
	}
	var reader io.Reader
	curPath := *parentPath
	curPath = path.Join(curPath, fileInfo.Name())
	if fileInfo.IsDir() {
		*parentPath = path.Join(*parentPath, fileInfo.Name())
	} else {
		err := s.writeStatusOK()
		if err != nil {
			return false, err
		}
		reader, err = s.reader.readFile(fileInfo)
		if err != nil {
			return false, err
		}
	}
	return handler(curPath, fileInfo, reader)
}

func (s *session) Uploader(location string) (uploadFn, io.Closer, error) {
	location = cleanPath(location)
	if s.mode != scpWrite {
		return nil, nil, fmt.Errorf("session mode must be write in uploader")
	}
	err := s.init(location)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to init session")
	}
	err = s.reader.readStatus()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to read status")
	}

	prevParentPath := location

	upload := func(ctx context.Context, fullpath string, info os.FileInfo, reader io.Reader) error {
		fullpath = cleanPath(fullpath)

		if !strings.HasPrefix(fullpath, location) {
			return errors.Errorf("file path must has prefix: %v", location)
		}

		localPrevParentPath := prevParentPath
		parentPath, basename := path.Split(fullpath)
		parentPath = cleanPath(parentPath)

		if parentPath != prevParentPath {
			prevParentPath = parentPath
		}
		err = adjustPath(localPrevParentPath, parentPath, s.enterDir, s.endDir)
		if err != nil {
			return err
		}

		if basename != info.Name() {
			info = fileinfo.NewInfo(basename, info.Size(), info.Mode(), info.ModTime(), info.IsDir())
		}

		if info.IsDir() {
			prevParentPath = path.Join(prevParentPath, info.Name())
			return s.enterDir(info)
		}
		return s.upload(info, reader)
	}
	return upload, s, nil
}

func (s *session) upload(info os.FileInfo, reader io.Reader) error {
	timestampCmd := s.scpTimestampCmd(info)
	err := s.sendMsgAndReadStatus(timestampCmd)
	if err != nil {
		return err
	}
	createCmd := s.scpCreateCmd(info)
	err = s.sendMsgAndReadStatus(createCmd)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		if _, err = io.Copy(s.writer, reader); err != nil {
			return err
		}
		err = s.writeStatusOK()
		if err == nil {
			err = s.reader.readStatus()
		}
	}
	return err
}

// scpTimestampCmd returns scp timestamp command for supplied info
func (s *session) scpTimestampCmd(info os.FileInfo) string {
	unixTimestamp := info.ModTime().Unix()
	return fmt.Sprintf("T%v 0 %v 0\n", unixTimestamp, unixTimestamp)
}

// scpCreateCmd returns scp create command for supplied info
func (s *session) scpCreateCmd(info os.FileInfo) string {
	mode := info.Mode()
	fileType := "C"
	size := info.Size()
	if info.IsDir() {
		fileType = "D"
		size = 0
	}
	fileMode := fmt.Sprintf("%v%04o", fileType, mode.Perm())[:5]
	return fmt.Sprintf("%v %d %s\n", fileMode, size, info.Name())
}

// adjustPath tracks current and previous relative path to adjust accordingly
func adjustPath(prev, current string, enterDir func(info os.FileInfo) error, endDir func() error) error {
	if prev == current {
		return nil
	}
	var prevElements []string
	var currElements []string

	if prev != "" {
		prevElements = strings.Split(prev, "/")
	}
	if prev != "" {
		currElements = strings.Split(current, "/")
	}
	if len(prevElements) < len(currElements) {
		for i := len(prevElements); i < len(currElements); i++ {
			dirInfo := fileinfo.NewInfo(currElements[i], 0, DefaultDirMode, time.Now(), true)
			if err := enterDir(dirInfo); err != nil {
				return err
			}
		}
	}
	downElements := make([]string, 0)
	for i := len(prevElements) - 1; i >= 0; i-- {
		prevElem := prevElements[i]
		currentElem := ""
		if i < len(currElements) {
			currentElem = currElements[i]
		}
		if currentElem == prevElem {
			break
		}
		if currentElem != "" {
			downElements = append(downElements, currentElem)
		}
		if err := endDir(); err != nil {
			return err
		}
	}
	for _, element := range downElements {
		dirInfo := fileinfo.NewInfo(element, 0, DefaultDirMode, time.Now(), true)
		if err := enterDir(dirInfo); err != nil {
			return err
		}
	}
	return nil
}
