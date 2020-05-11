/*
Copyright 2020 Jim Zhang (jim.zoumo@gmail.com). All rights reserved.

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

package exec

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
	"reflect"
	"testing"
)

func TestCommand(t *testing.T) {

	tests := []struct {
		name    string
		argName string
		args    []string
		want    *Cmd
	}{
		{
			"",
			"echo",
			[]string{"123"},
			&Cmd{argsHolder: &argsHolder{name: "echo", args: []string{"123"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Command(tt.argName, tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandContext(t *testing.T) {

	tests := []struct {
		name    string
		ctx     context.Context
		argName string
		args    []string
		want    *Cmd
	}{
		{
			"",
			context.TODO(),
			"echo",
			[]string{"123"},
			&Cmd{ctx: context.TODO(), argsHolder: &argsHolder{name: "echo", args: []string{"123"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommandContext(tt.ctx, tt.argName, tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCmd_Command(t *testing.T) {
	tests := []struct {
		name string
		cmd  *Cmd
		want *exec.Cmd
	}{
		{"", Command("echo", "123"), exec.Command("echo", "123")},
		{"", CommandContext(context.TODO(), "echo", "123"), exec.CommandContext(context.TODO(), "echo", "123")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cmd.Command(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.Command() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCmd_Pipe(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *Cmd
		pipeName string
		pipeArgs []string
		want     *Cmd
	}{
		{
			"",
			Command("echo", "123"),
			"sort",
			[]string{"test"},
			&Cmd{argsHolder: &argsHolder{name: "sort", args: []string{"test"}}, preCmd: Command("echo", "123")},
		},
		{
			"",
			CommandContext(context.TODO(), "echo", "123"),
			"sort",
			[]string{"test"},
			&Cmd{ctx: context.TODO(), argsHolder: &argsHolder{name: "sort", args: []string{"test"}}, preCmd: CommandContext(context.TODO(), "echo", "123")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cmd.Pipe(tt.pipeName, tt.pipeArgs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.Pipe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCmd_Run(t *testing.T) {

	tests := []struct {
		name    string
		cmd     *Cmd
		wantErr bool
	}{
		{"", Command("echo", "123").Pipe("sort"), false},
		{"", Command("echox", "123").Pipe("sort"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cmd.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Cmd.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCmd_Output(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *Cmd
		want    []byte
		wantErr bool
	}{
		{"", Command("echo"), nil, false},
		{"", Command("echo", "2\n1").Pipe("sort"), []byte("1\n2"), false},
		{"", Command("echo", "2\n1\n2").Pipe("sort").Pipe("uniq"), []byte("1\n2"), false},
		{"invalidOption", Command("echo", "2\n1").Pipe("sort", "-x"), nil, true},
		{"lookuperr", Command("echox", "123").Pipe("sort"), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.Output()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.Output() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.Output() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestCmd_CombinedOutput(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *Cmd
		want    []byte
		wantErr bool
	}{
		{"", Command("echo"), nil, false},
		{"", Command("echo", "2\n1").Pipe("sort"), []byte("1\n2"), false},
		{"", Command("echo", "2\n1\n2").Pipe("sort").Pipe("uniq"), []byte("1\n2"), false},
		{"invalidOption", Command("echo", "2\n1").Pipe("sort", "-x"), []byte("sort: invalid option -- 'x'\nTry 'sort --help' for more information."), true},
		{"lookuperr", Command("echox", "123").Pipe("sort"), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.CombinedOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.CombinedOutput() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestCmd_OutputClosure(t *testing.T) {
	tests := []struct {
		name    string
		cmd     func(...string) ([]byte, error)
		args    []string
		want    []byte
		wantErr bool
	}{
		{"", Command("echo").OutputClosure(), []string{"123"}, []byte("123"), false},
		{"", Command("echo", "2\n1").Pipe("sort").OutputClosure(), nil, []byte("1\n2"), false},
		{"", Command("echo", "2\n1").Pipe("sort").OutputClosure(), []string{""}, []byte("1\n2"), false},
		{"", Command("echo", "2\n1").Pipe("sort").OutputClosure(), []string{"-r"}, []byte("2\n1"), false},
		{"", Command("echox").OutputClosure(), []string{"-r"}, nil, true},
		{"", Command("echo", "2\n1").Pipe("sort").OutputClosure(), []string{"-x"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd(tt.args...)
			if (err != nil) != tt.wantErr {
				if eerr, ok := err.(*exec.ExitError); ok {
					t.Errorf("Cmd.OutputClosure() error = %v, stderr = %v, wantErr %v", eerr, string(eerr.Stderr), tt.wantErr)
				} else {
					t.Errorf("Cmd.OutputClosure() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.OutputClosure() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestCmd_OutputClosureManyTimes(t *testing.T) {
	type fields struct {
		args    []string
		want    []byte
		wantErr bool
	}

	tests := []struct {
		name   string
		cmd    func(...string) ([]byte, error)
		fields []fields
	}{
		{
			"",
			Command("echo").OutputClosure(),
			[]fields{
				{[]string{"123"}, []byte("123"), false},
				{[]string{"234"}, []byte("234"), false},
			},
		},
		{
			"",
			Command("bash", "-c").OutputClosure(),
			[]fields{
				{[]string{"xxx"}, nil, true},
				{[]string{"echo 123"}, []byte("123"), false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, f := range tt.fields {
				got, err := tt.cmd(f.args...)
				if (err != nil) != f.wantErr {
					if eerr, ok := err.(*exec.ExitError); ok {
						t.Errorf("Cmd.OutputClosure() error = %v, stderr = %v, wantErr %v", eerr, string(eerr.Stderr), f.wantErr)
					} else {
						t.Errorf("Cmd.OutputClosure() error = %v, wantErr %v", err, f.wantErr)
					}
					return
				}
				if !reflect.DeepEqual(got, f.want) {
					t.Errorf("Cmd.OutputClosure() = %v, want %v", string(got), string(f.want))
				}
			}
		})
	}
}

func TestCmd_CombinedOutputClosure(t *testing.T) {
	tests := []struct {
		name    string
		cmd     func(...string) ([]byte, error)
		args    []string
		want    []byte
		wantErr bool
	}{
		{"", Command("echo").CombinedOutputClosure(), []string{"123"}, []byte("123"), false},
		{"", Command("echo", "2\n1").Pipe("sort").CombinedOutputClosure(), nil, []byte("1\n2"), false},
		{"", Command("echo", "2\n1").Pipe("sort").CombinedOutputClosure(), []string{""}, []byte("1\n2"), false},
		{"", Command("echo", "2\n1").Pipe("sort").CombinedOutputClosure(), []string{"-r"}, []byte("2\n1"), false},
		{"lokkuperr", Command("echox").CombinedOutputClosure(), []string{"-r"}, nil, true},
		{"invalidOption", Command("echo", "2\n1").Pipe("sort").CombinedOutputClosure(), []string{"-x"}, []byte("sort: invalid option -- 'x'\nTry 'sort --help' for more information."), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd(tt.args...)
			if (err != nil) != tt.wantErr {
				if eerr, ok := err.(*exec.ExitError); ok {
					t.Errorf("Cmd.OutputClosure() error = %v, stderr = %v, wantErr %v", eerr, string(eerr.Stderr), tt.wantErr)
				} else {
					t.Errorf("Cmd.OutputClosure() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.OutputClosure() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestCmd_CombinedOutputClosureManyTimes(t *testing.T) {
	type fields struct {
		args    []string
		want    []byte
		wantErr bool
	}

	tests := []struct {
		name   string
		cmd    func(...string) ([]byte, error)
		fields []fields
	}{
		{
			"",
			Command("echo").CombinedOutputClosure(),
			[]fields{
				{[]string{"123"}, []byte("123"), false},
				{[]string{"234"}, []byte("234"), false},
			},
		},
		{
			"",
			Command("bash", "-c").CombinedOutputClosure(),
			[]fields{
				{[]string{"xxx"}, []byte("bash: xxx: command not found"), true},
				{[]string{"echo 123"}, []byte("123"), false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, f := range tt.fields {
				got, err := tt.cmd(f.args...)
				if (err != nil) != f.wantErr {
					if eerr, ok := err.(*exec.ExitError); ok {
						t.Errorf("Cmd.CombinedOutputClosure() error = %v, stderr = %v, wantErr %v", eerr, string(eerr.Stderr), f.wantErr)
					} else {
						t.Errorf("Cmd.CombinedOutputClosure() error = %v, wantErr %v", err, f.wantErr)
					}
				}
				if !reflect.DeepEqual(got, f.want) {
					t.Errorf("Cmd.CombinedOutputClosure() = %v, want %v", string(got), string(f.want))
				}
			}
		})
	}
}

func TestCmd_SetStdin(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *Cmd
		in      io.ReadWriter
		want    []byte
		wantErr bool
	}{
		{"", Command("sort"), bytes.NewBuffer([]byte("2\n1")), []byte("1\n2"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cmd.SetIO(tt.in, nil, nil)
			got, err := tt.cmd.Output()
			if (err != nil) != tt.wantErr {
				if eerr, ok := err.(*exec.ExitError); ok {
					t.Errorf("Cmd.SetStdin() error = %v, stderr = %v, wantErr %v", eerr, string(eerr.Stderr), tt.wantErr)
				} else {
					t.Errorf("Cmd.SetStdin() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(string(got), string(tt.want)) {
				t.Errorf("Cmd.SetStdin() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestCmd_SetStdout(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *Cmd
		out     io.ReadWriter
		want    []byte
		wantErr bool
	}{
		{"", Command("echo", "123"), new(bytes.Buffer), []byte("123"), false},
		{"", Command("echo", "2\n1").Pipe("sort"), new(bytes.Buffer), []byte("1\n2"), false},
		{"", Command("echo", "2\n1").Pipe("sort", "-x"), new(bytes.Buffer), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cmd.SetIO(nil, tt.out, nil)
			got, err := tt.cmd.Output()
			if (err != nil) != tt.wantErr {
				if eerr, ok := err.(*exec.ExitError); ok {
					t.Errorf("Cmd.SetStdout() error = %v, stderr = %v, wantErr %v", eerr, string(eerr.Stderr), tt.wantErr)
				} else {
					t.Errorf("Cmd.SetStdout() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(string(got), string(tt.want)) {
				t.Errorf("Cmd.SetStdout() = %v, want %v", string(got), string(tt.want))
			}

			got, err = ioutil.ReadAll(tt.out)
			if err != nil {
				t.Errorf("Failed to read bytes from tt.out")
				return
			}
			if !reflect.DeepEqual(string(bytes.TrimSpace(got)), string(tt.want)) {
				t.Errorf("Cmd.SetStdout() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestCmd_SetStderr(t *testing.T) {
	tests := []struct {
		name         string
		cmd          *Cmd
		out          io.ReadWriter
		wantCombined []byte
		wantStderr   []byte
		wantErr      bool
	}{
		{
			"",
			Command("echo", "123"),
			new(bytes.Buffer),
			[]byte("123"),
			nil,
			false,
		},
		{
			"",
			Command("echo", "2\n1").Pipe("sort"),
			new(bytes.Buffer),
			[]byte("1\n2"),
			nil,
			false,
		},
		{
			"",
			Command("echo", "2\n1").Pipe("sort", "-x"),
			new(bytes.Buffer),
			[]byte("sort: invalid option -- 'x'\nTry 'sort --help' for more information."),
			[]byte("sort: invalid option -- 'x'\nTry 'sort --help' for more information."),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cmd.SetIO(nil, nil, tt.out)
			got, err := tt.cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				if eerr, ok := err.(*exec.ExitError); ok {
					t.Errorf("Cmd.SetStderr() error = %v, stderr = %v, wantErr %v", eerr, string(eerr.Stderr), tt.wantErr)
				} else {
					t.Errorf("Cmd.SetStderr() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(string(got), string(tt.wantCombined)) {
				t.Errorf("Cmd.SetStderr() = %v, wantCombined %v", string(got), string(tt.wantCombined))
			}

			got, err = ioutil.ReadAll(tt.out)
			if err != nil {
				t.Errorf("Failed to read bytes from tt.out")
				return
			}
			if !reflect.DeepEqual(string(bytes.TrimSpace(got)), string(tt.wantStderr)) {
				t.Errorf("Cmd.SetStderr() = %v, wantStderr %v", string(got), string(tt.wantStderr))
			}
		})
	}
}

func TestCmd_RunForever(t *testing.T) {
	tests := []struct {
		name      string
		cmd       *Cmd
		startup   *Probe
		wantErr   bool
		errString string
	}{
		{
			"invalidOption",
			Command("sort", "-x"),
			nil,
			true,
			"exit status 2",
		},
		{
			"exitInRunForever",
			Command("sleep", "1"),
			&Probe{
				InitialDelaySeconds: 2,
			},
			true,
			ErrExitedInRunForever.Error(),
		},
		{
			"",
			Command("sleep", "1"),
			&Probe{
				SuccessThreshold: 3,
			},
			true,
			ErrExitedInRunForever.Error(),
		},
		{
			"",
			Command("sleep", "5"),
			&Probe{
				SuccessThreshold: 1,
				FailureThreshold: 1,
			},
			false,
			"",
		},
		{
			"",
			Command("sleep", "5"),
			&Probe{
				Handler: func(*exec.Cmd) error {
					return errors.New("failed run forever")
				},
				SuccessThreshold: 1,
				FailureThreshold: 1,
			},
			true,
			"probe failed: failed run forever",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.RunForever(tt.startup)
			if (err != nil) != tt.wantErr || (err != nil && err.Error() != tt.errString) {
				t.Errorf("Cmd.RunForever() error = %v, wantErr %v, wantErrStr %v", err, tt.wantErr, tt.errString)
			}
		})
	}
}
