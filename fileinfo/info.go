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
