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
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
	"time"
)

var (
	ErrExitedInRunForever = errors.New("exec: command should not exit in RunForever")
)

type argsHolder struct {
	name string
	args []string
}

func (c *argsHolder) Copy() *argsHolder {
	copy := *c
	return &copy
}

type ioHolder struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *ioHolder) SetIO(in io.Reader, out, err io.Writer) {
	c.stdin = in
	c.stdout = out
	c.stderr = err
}

func (c *ioHolder) GetIO() (in io.Reader, out, err io.Writer) {
	return c.stdin, c.stdout, c.stderr
}

// Cmd represents an external command being prepared or run basically.
// It also can combine several existing Command into a pipeline, just like
// running in shell: echo "3\n2\n1" | sort
//
// A Cmd cannot be reused after calling its Run, Output or CombinedOutput
// methods.
type Cmd struct {
	ctx        context.Context
	argsHolder *argsHolder
	ioHolder   *ioHolder

	cmdMutator func(name string, args []string) (string, []string)

	runtimeCmd *exec.Cmd
	preCmd     *Cmd

	started  bool
	finished bool
}

// Standard Command api follow os/exec.Command
func Command(name string, args ...string) *Cmd {
	return &Cmd{
		argsHolder: &argsHolder{
			name: name,
			args: args,
		},
	}
}

// Standard CommandContext api follow os/exec.CommandContext
func CommandContext(ctx context.Context, name string, args ...string) *Cmd {
	return &Cmd{
		ctx: ctx,
		argsHolder: &argsHolder{
			name: name,
			args: args,
		},
	}
}

// SetCmdMutator set a mutator function to mutator the runtime command's name and args
func (c *Cmd) SetCmdMutator(f func(name string, args []string) (string, []string)) {
	c.cmdMutator = f
}

func (c *Cmd) copy() *Cmd {
	newCmd := &Cmd{
		ctx:        c.ctx,
		argsHolder: c.argsHolder.Copy(),
		ioHolder:   c.ioHolder,
		cmdMutator: c.cmdMutator,
	}
	if c.preCmd != nil {
		newCmd.preCmd = c.preCmd.copy()
	}
	return newCmd
}

// Pipe creates a new command with given args and connects this command's
// standard output to new command's standard input
func (c *Cmd) Pipe(name string, args ...string) *Cmd {
	nextCmd := &Cmd{
		ctx:    c.ctx,
		preCmd: c,
		argsHolder: &argsHolder{
			name: name,
			args: args,
		},
		ioHolder:   c.ioHolder,
		cmdMutator: c.cmdMutator,
	}
	return nextCmd
}

// SetIO sets standard input/output/err output for command
func (c *Cmd) SetIO(in io.Reader, out, err io.Writer) {
	if c.ioHolder == nil {
		c.ioHolder = &ioHolder{}
	}
	c.ioHolder.SetIO(in, out, err)
}

func (c *Cmd) getIO() (in io.Reader, out, err io.Writer) {
	if c.ioHolder == nil {
		return nil, nil, nil
	}
	return c.ioHolder.GetIO()
}

func (c *Cmd) ensureCmd() {
	if c.runtimeCmd == nil {
		name := c.argsHolder.name
		args := c.argsHolder.args
		if c.cmdMutator != nil {
			name, args = c.cmdMutator(name, args)
		}
		if c.ctx != nil {
			c.runtimeCmd = exec.CommandContext(c.ctx, name, args...)
		} else {
			c.runtimeCmd = exec.Command(name, args...)
		}
		// reset std input/output for safety
		c.runtimeCmd.Stdin = nil
		c.runtimeCmd.Stdout = nil
		c.runtimeCmd.Stderr = nil
	}
}

// Command returns a runnable *exec.Cmd
func (c *Cmd) Command() *exec.Cmd {
	c.ensureCmd()
	return c.runtimeCmd
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

func (c *Cmd) setDefultProbe(startup *Probe) *Probe {
	if startup == nil {
		startup = &Probe{}
	}

	if startup.Handler == nil {
		startup.Handler = IsCmdRunningHandler
	}
	if startup.PeriodSeconds == 0 {
		startup.PeriodSeconds = 1
	}
	if startup.FailureThreshold == 0 {
		startup.FailureThreshold = 3
	}
	if startup.SuccessThreshold == 0 {
		startup.SuccessThreshold = 2
	}
	return startup
}

func (c *Cmd) RunForever(startup *Probe) error {
	err := c.Start()
	if err != nil {
		return err
	}

	done := make(chan struct{})
	errC := make(chan error)

	go func() {
		// wait for command exit
		errC <- c.Wait()
	}()

	startup = c.setDefultProbe(startup)
	worker := newWorker(c.Command(), startup, time.Now(), done)

	select {
	case err := <-errC:
		close(done) // stop worder
		if err != nil {
			return err
		}
		return ErrExitedInRunForever
	case err := <-worker.run():
		return err
	}
}

// Start starts the specified command but does not wait for it to complete.
//
// The Wait method will return the exit code and release associated resources
// once the command exits.
func (c *Cmd) Start() error {
	if c.started {
		return errors.New("exec: already started")
	}
	defer func() {
		c.started = true
	}()
	err := c.beforeStart()
	if err != nil {
		return err
	}
	err = c.runtimeCmd.Start()
	if err != nil {
		return err
	}
	if c.preCmd != nil {
		return c.preCmd.Start()
	}
	return nil
}

// beforeStart ensure runtime command firstly and setup its Stdout and Stderr.
// It also pipes pre command's stdout to this command's stdin, and use the same
// stderr to collect error message.
func (c *Cmd) beforeStart() error {
	c.ensureCmd()
	stdin, stdout, stderr := c.getIO()

	// setup stdin for first command, so that we can read input from it
	if stdin != nil && c.preCmd == nil {
		c.runtimeCmd.Stdin = stdin
	}
	// setup stdout and stderr for last command
	// the pre command's stdout and stderr will be set by pipe
	if c.runtimeCmd.Stdout == nil {
		c.runtimeCmd.Stdout = newWriterWithBuffer(stdout)
	}
	if c.runtimeCmd.Stderr == nil {
		c.runtimeCmd.Stderr = newWriterWithBuffer(stderr)
	}

	if c.preCmd != nil {
		preCmd := c.preCmd.Command()
		var err error
		// pre's output connect to cmd's input
		c.runtimeCmd.Stdin, err = preCmd.StdoutPipe()
		if err != nil {
			return err
		}
		// pre's error connect to cmd's error
		preCmd.Stderr = c.runtimeCmd.Stderr
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
		return errors.New("exec: not started")
	}
	if c.finished {
		return errors.New("exec: cmd finished")
	}

	defer func() {
		c.finished = true
	}()

	if c.preCmd != nil {
		if err := c.preCmd.Wait(); err != nil {
			return err
		}
	}
	err := c.runtimeCmd.Wait()
	return err
}

// CombinedOutput runs the command and returns its combined standard
// output and standard error.
func (c *Cmd) CombinedOutput() ([]byte, error) {
	err := c.Run()
	if err != nil {
		if eerr, ok := err.(*exec.ExitError); ok {
			stderr, _ := c.ReadStderr()
			eerr.Stderr = stderr
			return stderr, eerr
		}
		return nil, err
	}

	stdout, _ := c.ReadStdout()
	stderr, _ := c.ReadStderr()

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
			stderr, _ := c.ReadStderr()
			eerr.Stderr = stderr
			return nil, eerr
		}
		return nil, err
	}
	return c.ReadStdout()
}

// ReadStdout reads all bytes from command's standard output
// The command must have been finished by Wait.
func (c *Cmd) ReadStdout() ([]byte, error) {
	if !c.finished {
		return nil, errors.New("exec: not finished")
	}
	if c.runtimeCmd.Stdout != nil {
		if reader, ok := c.runtimeCmd.Stdout.(io.Reader); ok {
			msg, err := ioutil.ReadAll(reader)
			return bytes.TrimSpace(msg), err
		}
	}
	return nil, nil
}

// ReadStderr reads all bytes from command's standard error
// The command must have been finished by Wait.
func (c *Cmd) ReadStderr() ([]byte, error) {
	if !c.finished {
		return nil, errors.New("exec: not finished")
	}
	if c.runtimeCmd.Stderr != nil {
		if reader, ok := c.runtimeCmd.Stderr.(io.Reader); ok {
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
		for _, arg := range args {
			if arg == "" {
				continue
			}
			newCmd.argsHolder.args = append(newCmd.argsHolder.args, arg)
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
		for _, arg := range args {
			if arg == "" {
				continue
			}
			newCmd.argsHolder.args = append(newCmd.argsHolder.args, arg)
		}
		return newCmd.CombinedOutput()
	}
}
