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
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

var (
	entrypoint = "/bin/bash"
)

// QueryEscape escapes the string so it can be safely placed
// inside a shell command query.
func QueryEscape(arg string) string {
	return fmt.Sprintf("'%s'", strings.Replace(arg, "'", "'\\''", -1))
}

// Cmd represents an external command being prepared or run basically.
// It also can combine several existing Command into a pipeline, just like
// running in shell: echo "3\n2\n1" | sort
//
// A Cmd cannot be reused after calling its Run, Output or CombinedOutput
// methods.
type Cmd struct {
	started  bool
	finished bool
	args     []string
	cmd      *exec.Cmd
	pre      *Cmd
}

// Command returns a new command with args
func Command(args ...interface{}) *Cmd {
	c := new(Cmd)
	c.addArgs(args...)
	return c
}

func (c *Cmd) addArgs(args ...interface{}) {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			if c.cmd != nil {
				// add args to command directly
				c.cmd.Args = append(c.cmd.Args, v)
			} else {
				c.args = append(c.args, v)
			}
		case *exec.Cmd:
			if len(args) > 1 {
				panic("invalid argument, only one argument is allowed when *exec.Cmd is in args")
			}
			c.cmd = v
		case *Cmd:
			if len(args) > 1 {
				panic("invalid argument, only one argument is allowed when *Cmd is in args")
			}
			c.args = v.args
			c.cmd = v.cmd
			c.pre = v.pre
		default:
			panic(fmt.Sprintf("invalid argument %v", v))
		}
	}

}

func (c *Cmd) ensureCmd() {
	if c.cmd == nil {
		c.cmd = exec.Command(entrypoint, "-c", cmdLine(c.args))
	}
	// clean args after ensuring
	c.args = nil
}

func (c *Cmd) copy() *Cmd {
	newCmd := &Cmd{
		args: c.args,
		cmd:  c.cmd,
	}
	if c.pre != nil {
		newCmd.pre = c.pre.copy()
	}
	return newCmd
}

// Pipe creates a new command with given args and connects this command's
// standard output to new command's standard input
// The args can be string, *exec.Cmd or *Cmd. Only one argument is allowed when
// *exec.Cmd or *Cmd exists.
func (c *Cmd) Pipe(args ...interface{}) *Cmd {
	if len(args) == 0 {
		return c
	}
	newCmd := &Cmd{
		pre: c,
	}
	newCmd.addArgs(args...)
	return newCmd
}

// Command returns a runnable *exec.Cmd
func (c *Cmd) Command() *exec.Cmd {
	c.ensureCmd()
	return c.cmd
}

func (c *Cmd) beforeStart() error {
	c.ensureCmd()
	if c.cmd.Stdout == nil {
		// the last Cmd
		c.cmd.Stdout = new(bytes.Buffer)
	}
	if c.cmd.Stderr == nil {
		// the last Cmd
		c.cmd.Stderr = new(bytes.Buffer)
	}
	if c.pre != nil {
		preCmd := c.pre.Command()
		var err error
		// pre's output connect to cmd's input
		c.cmd.Stdin, err = preCmd.StdoutPipe()
		if err != nil {
			return err
		}
		// pre's error connect to cmd's error
		preCmd.Stderr = c.cmd.Stderr
	}
	return nil
}

// Run starts the specified command and waits for it to complete.
//
// The returned error is nil if the command runs, you can safely get output from
// Stdout() or error message from Stderr()
//
// If the command starts but does not complete successfully, the error is of
// type *ExitError. Other error types may be returned for other situations.
//
// If the calling goroutine has locked the operating system thread
// with runtime.LockOSThread and modified any inheritable OS-level
// thread state (for example, Linux or Plan 9 name spaces), the new
// process will inherit the caller's thread state./
func (c *Cmd) Run() error {
	err := c.Start()
	if err != nil {
		return err
	}
	return c.Wait()
}

// Start starts the specified command but does not wait for it to complete.
//
// The Wait method will return the exit code and release associated resources
// once the command exits.
func (c *Cmd) Start() error {
	if c.started {
		return errors.New("shell: already started")
	}
	defer func() {
		c.started = true
	}()
	err := c.beforeStart()
	if err != nil {
		return err
	}
	err = c.cmd.Start()
	if err != nil {
		return err
	}
	if c.pre != nil {
		return c.pre.Start()
	}
	return nil
}

// Wait waits for the command to exit and waits for any copying to
// stdin or copying from stdout or stderr to complete.
//
// The command must have been started by Start.
//
// The returned error is nil if the command runs, you can safely get output from
// Stdout() or error message from Stderr()
//
// If the command fails to run or doesn't complete successfully, the
// error is of type *ExitError. Other error types may be
// returned for I/O problems.
//
// If any of c.Stdin, c.Stdout or c.Stderr are not an *os.File, Wait also waits
// for the respective I/O loop copying to or from the process to complete.
//
// Wait releases any resources associated with the Cmd.
func (c *Cmd) Wait() error {
	if !c.started {
		return errors.New("shell: not started")
	}
	defer func() {
		c.finished = true
	}()

	if c.pre != nil {
		if err := c.pre.Wait(); err != nil {
			return err
		}
	}
	err := c.cmd.Wait()
	return err
}

// CombinedOutput runs the command and returns its combined standard
// output and standard error.
func (c *Cmd) CombinedOutput() ([]byte, error) {
	err := c.Run()
	if err != nil {
		if eerr, ok := err.(*exec.ExitError); ok {
			stderr, _ := c.Stderr()
			eerr.Stderr = stderr
			return stderr, eerr
		}
		return nil, err
	}

	stdout, _ := c.Stdout()
	stderr, _ := c.Stderr()

	merged := bytes.Buffer{}
	merged.Write(stdout)
	merged.Write(stderr)
	return merged.Bytes(), nil
}

// Output runs the command and returns its standard output.
// Any returned error will usually be of type *ExitError.
func (c *Cmd) Output() ([]byte, error) {
	err := c.Run()
	if err != nil {
		if eerr, ok := err.(*exec.ExitError); ok {
			stderr, _ := c.Stderr()
			eerr.Stderr = stderr
			return nil, eerr
		}
		return nil, err
	}
	return c.Stdout()
}

// Stdout reads all bytes from command's standard output
// The command must have been finished by Wait.
func (c *Cmd) Stdout() ([]byte, error) {
	if !c.finished {
		return nil, errors.New("shell: not finished")
	}
	if c.cmd.Stdout != nil {
		if reader, ok := c.cmd.Stdout.(io.Reader); ok {
			msg, err := ioutil.ReadAll(reader)
			return bytes.TrimSpace(msg), err
		}
	}
	return nil, nil
}

// Stderr reads all bytes from command's standard error
// The command must have been finished by Wait.
func (c *Cmd) Stderr() ([]byte, error) {
	if !c.finished {
		return nil, errors.New("shell: not finished")
	}
	if c.cmd.Stderr != nil {
		if reader, ok := c.cmd.Stderr.(io.Reader); ok {
			msg, err := ioutil.ReadAll(reader)
			return bytes.TrimSpace(msg), err
		}
	}
	return nil, nil
}

// OutputClosure returns a function closure allowing you to call
// this command latter. The closure runs the command and reads all
// bytes from standard output
//
// echo := Command("echo").OutputClosure()
// echo("123")
// echo("321")
func (c *Cmd) OutputClosure() func(...string) ([]byte, error) {
	return func(args ...string) ([]byte, error) {
		newCmd := c.copy()
		for _, a := range args {
			newCmd.addArgs(a)
		}
		return newCmd.Output()
	}
}

// CombinedOutputClosure returns function closure allowing you to call
// this command latter. The closure runs the command and reads all
// bytes from combined standard output and standard error
func (c *Cmd) CombinedOutputClosure() func(...string) ([]byte, error) {
	return func(args ...string) ([]byte, error) {
		newCmd := c.copy()
		for _, a := range args {
			newCmd.addArgs(a)
		}
		return newCmd.CombinedOutput()
	}
}

func cmdLine(args []string) string {
	return strings.Join(args, " ")
}
