/*
Copyright 2016 caicloud authors. All rights reserved.

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

package execd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	ps "github.com/keybase/go-ps"
)

const (
	crashBackoff = 3
)

var (
	// ErrNotRunning ...
	ErrNotRunning = errors.New("execd: process is not running")
)

// D represents an external daemon command being prepared or run.
//
// A daemon is a long-time running background process providing reliable services for others.
// If your program needs graceful shutdown, you can SetGracePeriod or SetGracefulShutDown.
type D struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	Env []string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	Dir string

	// Stdin specifies the process's standard input.
	//
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	//
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	//
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the command
	// over a pipe. In this case, Wait does not complete until the goroutine
	// stops copying, either because it has reached the end of Stdin
	// (EOF or a read error) or because writing to the pipe returned an error.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If either is an *os.File, the corresponding output from the process
	// is connected directly to that file.
	//
	// Otherwise, during the execution of the command a separate goroutine
	// reads from the process over a pipe and delivers that data to the
	// corresponding Writer. In this case, Wait does not complete until the
	// goroutine reaches EOF or encounters an error.
	//
	// If Stdout and Stderr are the same writer, and have a type that can
	// be compared with ==, at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer

	// SysProcAttr holds optional, operating system-specific attributes.
	// Run passes it to os.StartProcess as the os.ProcAttr's Sys field.
	SysProcAttr *syscall.SysProcAttr

	cmd *exec.Cmd

	gracePeriod      time.Duration
	gracefulShutDown func(*exec.Cmd) error

	lookPathErr error
	stopCh      chan struct{}
	errCh       chan error
}

// Daemon returns the D struct to execute the named program with
// the given arguments.
//
// It sets only the Path and Args in the returned structure.
//
// If name contains no path separators, Command uses exec.LookPath to
// resolve name to a complete path if possible. Otherwise it uses name
// directly as Path.
//
// The returned D's Args field is constructed from the command name
// followed by the elements of arg, so arg should not include the
// command name itself. For example, Daemon("echo", "hello").
// Args[0] is always name, not the possibly resolved Path.
func Daemon(name string, arg ...string) *D {
	cmd := &D{
		Path: name,
		Args: append([]string{name}, arg...),
	}
	if filepath.Base(name) == name {
		if lp, err := exec.LookPath(name); err != nil {
			cmd.lookPathErr = err
		} else {
			cmd.Path = lp
		}
	}
	return cmd
}

// DaemonFrom retrieves the necessary information from given exec.Cmd
// and returns a new D struct
func DaemonFrom(c *exec.Cmd) *D {
	cmd := convertFromExec(c)
	return cmd
}

// Command returns the running exec.Cmd struct in D
func (c *D) Command() *exec.Cmd {
	return c.cmd
}

// Pid returns the running process's pid
// if the daemon is not running, it will return ErrNotRunning
func (c *D) Pid() (int, error) {
	if !c.IsRunning() {
		return 0, ErrNotRunning
	}
	return c.cmd.Process.Pid, nil
}

// Signal sends a signal to the daemon pocess.
// if the daemon is not running, it will return an ErrNotRunning
func (c *D) Signal(signal os.Signal) error {
	if !c.IsRunning() {
		return ErrNotRunning
	}
	return c.cmd.Process.Signal(signal)
}

// Name returns the name of daemon
// Generally, it is the first arg in Args
func (c *D) Name() string {
	if len(c.Args) > 0 {
		return c.Args[0]
	}
	return c.Path
}

// SetGracePeriod sets graceful shutdown period time
// It will be invalidated when SetGracefulShutDown is called
//
// If the grace period == 0, execd sends SIGKILL to the daemon process immediately
// If the grace period > 0 , execd sends SIGTERM to the daemon. If the daemon doesn't
// terminate within the grace period, a SIGKILL will be sent and the daemon violently terminated
func (c *D) SetGracePeriod(d time.Duration) {
	c.gracePeriod = d
}

// SetGracefulShutDown sets graceful shutdown handler to override
// the default behavior
func (c *D) SetGracefulShutDown(f func(*exec.Cmd) error) {
	c.gracefulShutDown = f
}

// RunForever starts the specified command and waits for it to complete in another goroutine.
// If there is no error, the daemon will run forever.
//
// In the meantime, It starts a goroutine to keep the backgroud process alive.
// But if the error occurs more than `crashBackOff` times when command is starting,
// it will stop tracking anymore.
func (c *D) RunForever() error {
	if c.lookPathErr != nil {
		return c.lookPathErr
	}
	if c.cmd == nil {
		c.cmd = c.delegate()
	}
	if c.stopCh == nil {
		c.stopCh = make(chan struct{})
	}
	if c.errCh == nil {
		c.errCh = make(chan error)
	}

	err := c.run()
	if err != nil {
		close(c.stopCh)
		close(c.errCh)
		return err
	}
	c.keepalive()
	c.reportError()

	return nil
}

// IsRunning returns true if the daemon is still running background
func (c *D) IsRunning() bool {
	if c.cmd == nil {
		return false
	}
	if c.cmd.Process == nil {
		return false
	}
	process, err := ps.FindProcess(c.cmd.Process.Pid)
	if err != nil {
		panic(err)
	}
	if process == nil && err == nil {
		// not found
		return false
	}
	return true
}

// Stop stops the daemon process
func (c *D) Stop() error {
	if c.stopCh == nil {
		return errors.New("execd: stop must be called after run")
	}

	close(c.stopCh)

	if c.gracefulShutDown != nil && c.IsRunning() {
		return c.gracefulShutDown(c.cmd)
	}

	return c.shutdown()

}

func (c *D) shutdown() error {
	if c.gracePeriod > 0 {
		err := c.cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
		<-time.After(c.gracePeriod)
	}
	if c.IsRunning() {
		err := c.cmd.Process.Kill()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *D) run() error {
	if c.cmd == nil {
		return errors.New("execd: no command")
	}

	if err := c.cmd.Start(); err != nil {
		return err
	}

	go func() {
		// maybe killed
		c.errCh <- c.cmd.Wait()
	}()

	return nil
}

func (c *D) keepalive() {
	go func() {
		tick := time.NewTicker(1 * time.Second)
		defer tick.Stop()
		restartErrTimes := 0
		for {
			select {
			case <-tick.C:
				if !c.IsRunning() {
					c.cmd = c.delegate()
					err := c.run()
					if err != nil {
						if restartErrTimes >= crashBackoff {
							fmt.Printf("execd(%v): too many errors occur when restarting the process, stop the daemon\n", c.Name())
							c.Stop()
							return
						}
						fmt.Printf("execd(%v): error restart command: %v\n", c.Name(), err)
						restartErrTimes++
					}
				}
			case <-c.stopCh:
				return
			}
		}

	}()
}

func (c *D) reportError() {
	go func() {
		for {
			select {
			case err := <-c.errCh:
				if err != nil {
					fmt.Printf("execd(%v): receive an error, %v\n", c.Name(), err)
				}
			case <-c.stopCh:
				return
			}
		}
	}()
}

func (c *D) delegate() *exec.Cmd {
	return convertToExec(c)
}

func convertToExec(c *D) *exec.Cmd {
	cmd := &exec.Cmd{
		Path:        c.Path,
		Args:        c.Args,
		Env:         c.Env,
		Dir:         c.Dir,
		Stdin:       c.Stdin,
		Stderr:      c.Stderr,
		Stdout:      c.Stdout,
		SysProcAttr: c.SysProcAttr,
	}
	return cmd
}

func convertFromExec(c *exec.Cmd) *D {
	cmd := &D{
		Path:        c.Path,
		Args:        c.Args,
		Env:         c.Env,
		Dir:         c.Dir,
		Stdin:       c.Stdin,
		Stderr:      c.Stderr,
		Stdout:      c.Stdout,
		SysProcAttr: c.SysProcAttr,
	}
	return cmd
}
