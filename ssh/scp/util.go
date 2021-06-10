package scp

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/zoumo/golib/fileinfo"
)

//ParseFileInfo returns new info from SCP response
func ParseFileInfo(resp string, modified *time.Time) (os.FileInfo, error) {
	elements := strings.SplitN(resp, " ", 3)
	if len(elements) != 3 {
		return nil, fmt.Errorf("invalid scp response: %v", resp)
	}
	isDir := strings.HasPrefix(elements[0], "D")
	modeStr := elements[0][1:]
	mode, err := strconv.ParseInt(modeStr, 8, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid file mode: %v", modeStr)
	}
	sizeStr := elements[1]
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid file size: %v", modeStr)
	}
	name := strings.Trim(elements[2], "\r\n")
	if modified == nil {
		now := time.Now()
		modified = &now
	}
	return fileinfo.NewInfo(name, size, os.FileMode(mode), *modified, isDir), nil
}

//ParseTimeResponse parases respons time
func ParseTimeResponse(response string) (*time.Time, error) {
	elements := strings.SplitN(response, " ", 4)
	if len(elements) != 4 {
		return nil, fmt.Errorf("invalid timestamp response: %v", response)
	}
	timestampStr := elements[0][1:]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid timestamp: %v", timestampStr)
	}
	msecLiteral := elements[1]
	nsec, _ := strconv.ParseInt(msecLiteral, 10, 64)
	ts := time.Unix(timestamp, nsec*1000)
	return &ts, nil
}
