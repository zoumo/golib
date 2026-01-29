// Copyright 2022 jim.zoumo@gmail.com
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

package log

import (
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"k8s.io/klog/v2/klogr"
)

// cleanup logger for test
func cleanup() {
	singleton = newPlaceHolderLogger()
	Log = singleton
}

func TestSetLogrLogger(t *testing.T) {
	type args struct {
		l logr.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"noop",
			args{
				l: logr.Discard(),
			},
		},
		{
			"klogr",
			args{
				l: klogr.New(),
			},
		},
		{
			"zapr",
			args{
				l: func() logr.Logger {
					zapLog, err := zap.NewDevelopment()
					if err != nil {
						panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
					}
					return zapr.NewLogger(zapLog)
				}(),
			},
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			cleanup()
			Log.Info("before set logger")
			logger1 := Log.WithName("name1").WithName("sub1").WithValues("with", "value")
			logger2 := Log.WithName("name2").WithName("sub2").WithValues("with", "value")
			SetLogrLogger(tt.args.l)
			Log.Info("after set logger")
			logger1.Info("test", "key", "value")
			logger2.Info("test", "key", "value")
			// t.Error("xx")
		})
	}
}
