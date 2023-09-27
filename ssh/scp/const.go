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

import "os"

const (
	//FileToken crate file token
	FileToken = 'C'
	//DirToken create directory token
	DirToken = 'D'
	//TimestampToken timestamp token
	TimestampToken = 'T'
	//EndDirToken end of dir token
	EndDirToken = 'E'
	//WarningToken warning token
	WarningToken = 0x1
	//ErrorToken error token
	ErrorToken = 0x2
	// StatusOK
	StatusOK = 0x0
)

const (
	//DefaultPort default SSH port
	DefaultPort = 22
)

const (
	//DefaultDirMode folder mode default
	DefaultDirMode = os.ModeDir | 0755
)
