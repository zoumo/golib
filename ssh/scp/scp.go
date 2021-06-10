package scp

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"golang.org/x/crypto/ssh"

	"github.com/zoumo/golib/fileinfo"
	"github.com/zoumo/golib/log"
)

type stateFn func(string) (os.FileInfo, error)

type SCP struct {
	client *ssh.Client
	fs     afero.Fs
	logger logr.Logger
}

func New(client *ssh.Client, fs afero.Fs, logger logr.Logger) *SCP {
	if logger == nil {
		logger = log.NewNullLogger()
	}
	return &SCP{
		client: client,
		fs:     fs,
		logger: logger,
	}
}

func (s *SCP) Stat(dst string) (os.FileInfo, error) {
	info, _, err := s.Open(dst)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *SCP) Open(dst string) (os.FileInfo, io.Reader, error) {
	if err := validateSCPPath(dst); err != nil {
		return nil, nil, err
	}
	session, err := newSession(s.client, scpRead, 0, s.logger)
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()

	return session.Open(dst)
}

func (s *SCP) beforeCopy(source, target string, sourceStat, targetStat stateFn, notExistHandler func(target string) error) error {
	if err := validateSCPPath(source); err != nil {
		return err
	}
	if err := validateSCPPath(target); err != nil {
		return err
	}
	sourceInfo, err := sourceStat(source)
	if err != nil {
		return errors.Wrap(err, "failed to get source file stat")
	}

	targetInfo, err := targetStat(target)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to get target file stat")
	}

	if os.IsNotExist(err) {
		if err := notExistHandler(target); err != nil {
			return err
		}
	} else if sourceInfo.IsDir() != targetInfo.IsDir() {
		return errors.Errorf(
			"source(paht=%s,isDir=%v) and target(path=%s,isDir=%v) path are not the same type",
			source,
			sourceInfo.IsDir(),
			target,
			targetInfo.IsDir(),
		)
	}
	return nil
}

// Download downloads files from remote to local
// It is different from linux scp, if local path exists, it must be the same type with remote path
// If you are downloading a regular file, the local path must contain file name otherwise scp will
// use the last element of path as its file name
func (s *SCP) Download(ctx context.Context, remote, local string) error {
	local = cleanPath(local)
	remote = cleanPath(remote)

	err := s.beforeCopy(remote, local, s.Stat, s.fs.Stat, func(target string) error {
		// mkdir for local path's dir
		if err := s.fs.MkdirAll(path.Dir(target), DefaultDirMode); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	session, err := newSession(s.client, scpRead, 0, s.logger)
	if err != nil {
		return err
	}
	defer session.Close()

	localName := path.Base(local)

	err = session.Walk(remote, func(fpath string, finfo os.FileInfo, reader io.Reader) (bool, error) {
		var fullpath string

		if fpath == remote {
			finfo = fileinfo.NewInfo(localName, finfo.Size(), finfo.Mode(), finfo.ModTime(), finfo.IsDir())
			fullpath = local
		} else {
			rel, err := filepath.Rel(remote, fpath)
			if err != nil {
				return false, err
			}
			fullpath = path.Join(local, rel)
		}

		s.logger.V(3).Info("scp download", "from", fpath, "to", fullpath, "isDir", finfo.IsDir())

		if finfo.IsDir() {
			if err := s.fs.MkdirAll(fullpath, finfo.Mode().Perm()); err != nil {
				return false, err
			}
			return true, nil
		}

		if err := s.writeFile(fullpath, finfo, reader); err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *SCP) writeFile(path string, info os.FileInfo, content io.Reader) error {
	if err := afero.WriteReader(s.fs, path, content); err != nil {
		return err
	}
	if err := s.fs.Chmod(path, info.Mode()); err != nil {
		return err
	}
	if err := s.fs.Chtimes(path, info.ModTime(), info.ModTime()); err != nil {
		return err
	}
	return nil
}

// Upload uploads files from local to remote
// It is different from linux scp, if remote path exists, it must be the same type with local path
// If you are uploading a regular file, the remote path must contain file name otherwise scp will
// use the last element of path as its file name
func (s *SCP) Upload(ctx context.Context, local, remote string) error {
	local = cleanPath(local)
	remote = cleanPath(remote)

	err := s.beforeCopy(local, remote, s.fs.Stat, s.Stat, func(target string) error {
		// create remote dir
		session, err := s.client.NewSession()
		if err != nil {
			return err
		}
		defer session.Close()
		msg, err := session.CombinedOutput(fmt.Sprintf("mkdir -p %s", path.Dir(target)))
		if err != nil {
			return errors.Wrapf(err, "create remote dir failed, receive msg: %v", string(msg))
		}
		return nil
	})

	if err != nil {
		return err
	}

	session, err := newSession(s.client, scpWrite, 0, s.logger)
	if err != nil {
		return err
	}
	defer session.Close()

	upload, closer, err := session.Uploader(path.Dir(remote))
	if err != nil {
		return err
	}
	defer closer.Close()

	err = afero.Walk(s.fs, local, func(fpath string, finfo os.FileInfo, perr error) error {
		if perr != nil {
			return perr
		}
		if fileinfo.IsSymlink(finfo) {
			// get real file info follow symlink to get true size
			realInfo, err := s.fs.Stat(fpath)
			if err != nil {
				return err
			}
			if realInfo.IsDir() {
				// [by design] dir under symbolic link will be ignored,
				// it is difficult to avoid loops.
				s.logger.V(3).Info("ignore dir behind symbolic link", "path", fpath)
				return nil
			}
			// use real file's info
			finfo = realInfo
		}

		rel, err := filepath.Rel(local, fpath)
		if err != nil {
			return err
		}
		fullpath := cleanPath(path.Join(remote, rel))

		var content io.Reader
		if !finfo.IsDir() {
			f, err := s.fs.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()
			content = f
		}

		s.logger.V(3).Info("scp upload", "from", fpath, "to", fullpath, "isDir", finfo.IsDir())
		return upload(ctx, fullpath, finfo, content)
	})

	return err
}

func cleanPath(p string) string {
	p = path.Clean(p)
	p = strings.TrimRight(p, "/")
	return p
}

func validateSCPPath(fpath string) error {
	fpath = path.Clean(fpath)
	if fpath == "/" {
		return errors.New("can not use root(/) as file path")
	}
	return nil
}
