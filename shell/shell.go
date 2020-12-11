/*
Copyright 2019 Jim Zhang (jim.zoumo@gmail.com). All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shell

import (
	"context"
	"fmt"
	"strings"

	"github.com/zoumo/golib/exec"
)

var (
	entrypoint = "/bin/bash"
)

func shellCmdMutator(name string, args []string) (string, []string) {
	return entrypoint, []string{"-c", strings.Join(append([]string{name}, args...), " ")}
}

// QueryEscape escapes the string so it can be safely placed
// inside a shell command query.
func QueryEscape(arg string) string {
	return fmt.Sprintf("'%s'", strings.Replace(arg, "'", "'\\''", -1))
}

// Command returns a new command with args
// Running shell command
func Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.SetCmdMutator(shellCmdMutator)
	return cmd
}

func CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SetCmdMutator(shellCmdMutator)
	return cmd
}
