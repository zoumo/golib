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
