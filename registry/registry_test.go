/*
Copyright 2018 Jim Zhang (jim.zoumo@gmail.com). All rights reserved.

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

package registry

import (
	"reflect"
	"testing"
)

func Test_registry_Register(t *testing.T) {
	r := New(nil)
	type args struct {
		name string
		v    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{"test", 1}, false},
		{"", args{"test", 2}, true},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Register(tt.args.name, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("registry.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registry_Get(t *testing.T) {
	r := New(nil)
	r.Register("test", 1)

	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{"", args{"test"}, 1, true},
		{"", args{"test2"}, nil, false},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := r.Get(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("registry.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("registry.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
