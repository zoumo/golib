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
	"reflect"
	"testing"

	"github.com/zoumo/golib/exec"
)

func TestQueryEscape(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"", "", `''`},
		{"", `a`, `'a'`},
		{"", `"a"`, `'"a"'`},
		{"", "'3\n2\n1'", "''\\''3\n2\n1'\\'''"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QueryEscape(tt.arg); got != tt.want {
				t.Errorf("QueryEscape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShell_Run(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *exec.Cmd
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

func TestShell_Output(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *exec.Cmd
		want    []byte
		wantErr bool
	}{
		{"", Command("echo"), nil, false},
		{"", Command("echo '2\n1' | sort"), []byte("1\n2"), false},
		{"", Command("echo '2\n1\n2' | sort | uniq"), []byte("1\n2"), false},
		{"", Command("echo", "'2\n1'").Pipe("sort"), []byte("1\n2"), false},
		{"", Command("echo", "2\n1").Pipe("sort"), nil, true},
		{"", Command("echox", "123").Pipe("sort"), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.Output()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.CombinedOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.CombinedOutput() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestShell_CombinedOutput(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *exec.Cmd
		want    []byte
		wantErr bool
	}{
		{"", Command("echo"), nil, false},
		{"", Command("echo '2\n1' | sort"), []byte("1\n2"), false},
		{"", Command("echo '2\n1\n2' | sort | uniq"), []byte("1\n2"), false},
		{"", Command("echo", "'2\n1'").Pipe("sort"), []byte("1\n2"), false},
		{"", Command("echo", "2\n1").Pipe("sort"), []byte("/bin/bash: line 1: 1: command not found"), true},
		{"", Command("echox", "123").Pipe("sort"), []byte("/bin/bash: echox: command not found"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.CombinedOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.CombinedOutput() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestShell_OutputClosure(t *testing.T) {
	tests := []struct {
		name    string
		cmd     func(...string) ([]byte, error)
		args    []string
		want    []byte
		wantErr bool
	}{
		{"", Command("echo").OutputClosure(), []string{"123"}, []byte("123"), false},
		{"", Command("echo '2\n1'").Pipe("sort").OutputClosure(), []string{""}, []byte("1\n2"), false},
		{"", Command("echo '2\n1'").Pipe("sort").OutputClosure(), []string{"-r"}, []byte("2\n1"), false},
		{"", Command("echo '2\n1' | sort").OutputClosure(), []string{"-r"}, []byte("2\n1"), false},
		{"", Command("echox").OutputClosure(), []string{"-r"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.OutputClosure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.OutputClosure() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestShell_OutputClosureManyTimes(t *testing.T) {
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
				{[]string{"'echo 123'"}, []byte("123"), false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, f := range tt.fields {
				got, err := tt.cmd(f.args...)
				if (err != nil) != f.wantErr {
					t.Errorf("Cmd.OutputClosure() error = %v, wantErr %v", err, f.wantErr)
					return
				}
				if !reflect.DeepEqual(got, f.want) {
					t.Errorf("Cmd.OutputClosure() = %v, want %v", string(got), string(f.want))
				}
			}
		})
	}
}

func TestShell_CombinedOutputClosure(t *testing.T) {
	tests := []struct {
		name    string
		cmd     func(...string) ([]byte, error)
		args    []string
		want    []byte
		wantErr bool
	}{
		{"", Command("echo").CombinedOutputClosure(), []string{"123"}, []byte("123"), false},
		{"", Command("echo '2\n1'").Pipe("sort").CombinedOutputClosure(), []string{""}, []byte("1\n2"), false},
		{"", Command("echo '2\n1'").Pipe("sort").CombinedOutputClosure(), []string{"-r"}, []byte("2\n1"), false},
		{"", Command("echo '2\n1' | sort").CombinedOutputClosure(), []string{"-r"}, []byte("2\n1"), false},
		{"", Command("echox").CombinedOutputClosure(), []string{"-r"}, []byte("/bin/bash: echox: command not found"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.OutputClosure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cmd.OutputClosure() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
