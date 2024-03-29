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

package exec

import (
	"fmt"
	"os/exec"

	"github.com/keybase/go-ps"
)

func IsCmdRunning(cmd *exec.Cmd) bool {
	if cmd == nil {
		return false
	}
	if cmd.Process == nil {
		return false
	}
	process, err := ps.FindProcess(cmd.Process.Pid)
	if err != nil {
		panic(err)
	}
	if process == nil && err == nil {
		// not found
		return false
	}
	return true
}

func IsCmdRunningHandler(cmd *exec.Cmd) error {
	if running := IsCmdRunning(cmd); !running {
		return fmt.Errorf("command is not running")
	}
	return nil
}
