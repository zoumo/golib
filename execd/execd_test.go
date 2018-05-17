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
	"fmt"
	"os"
	"testing"
	"time"

	ps "github.com/keybase/go-ps"
	"github.com/moby/moby/pkg/reexec"
)

func init() {
	reexec.Register("execd-test-run", func() {
		var i int
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("===> execd-test: ", i)
			i++
		}
	})

	reexec.Register("execd-test-stop", func() {
		var i int
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("===> execd-test: ", i)
			i++
		}
	})

}

func TestRun(t *testing.T) {

	if reexec.Init() {
		os.Exit(0)
	}

	cmd := DaemonFrom(reexec.Command("execd-test-run"))
	err := cmd.RunForever()
	if err != nil {
		t.Error(err)
	}

	timer := time.NewTimer(3 * time.Second)
	killer := time.NewTimer(500 * time.Millisecond)

LOOP:
	for {
		select {
		case <-timer.C:
			p, err := ps.FindProcess(cmd.Command().Process.Pid)
			if p == nil && err == nil {
				t.Error("no process")
			}
			break LOOP
		case <-killer.C:
			cmd.Command().Process.Kill()
		}
	}

	cmd.Stop()
	time.Sleep(time.Second)

}

func TestStop(t *testing.T) {
	if reexec.Init() {
		os.Exit(0)
	}

	cmd := DaemonFrom(reexec.Command("execd-test-stop"))
	err := cmd.RunForever()
	if err != nil {
		t.Error(err)
	}
	cmd.SetGracePeriod(500 * time.Millisecond)

	cmd.Stop()
	<-time.After(1 * time.Second)

	if cmd.IsRunning() {
		t.Error("still running")
	}
}

func TestCrasLoopBackoff(t *testing.T) {
	cmd := &D{
		Path: "/not-found-path",
		Args: []string{"execd-test-crash"},
	}

	cmd.keepalive()
	cmd.reportError()
	<-time.After(5 * time.Second)

	if cmd.IsRunning() {
		t.Error("still running")
	}
	cmd.Stop()
}
