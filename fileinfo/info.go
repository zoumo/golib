package fileinfo

import (
	"os"
	"time"
)

// NewInfo returns a ew file Info
func NewInfo(name string, size int64, mode os.FileMode, mtime time.Time, isDir bool) os.FileInfo {
	return &Info{
		name:    name,
		size:    size,
		mode:    mode,
		modTime: mtime,
		isDir:   isDir,
	}
}

// Info represents a file info
type Info struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

// Name returns a name
func (i *Info) Name() string {
	return i.name
}

// Size returns file size
func (i *Info) Size() int64 {
	return i.size
}

// Mode returns file mode
func (i *Info) Mode() os.FileMode {
	return i.mode
}

// ModTime returns modification time
func (i *Info) ModTime() time.Time {
	return i.modTime
}

// IsDir returns true if resoruce is directory
func (i *Info) IsDir() bool {
	return i.isDir
}

// Sys returns sys object
func (i *Info) Sys() interface{} {
	return nil
}
