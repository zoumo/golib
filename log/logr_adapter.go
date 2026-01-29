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

import "github.com/go-logr/logr"

// FromLogr creates a Logger from any logr.Logger implementation.
// This works with both logr v0.4.0 (interface) and v1.0.0+ (struct).
func FromLogr(l logr.Logger) Logger {
	return &logrAdapter{logger: l}
}

// logrAdapter wraps a logr.Logger (either v0.4.0 interface or v1.0.0+ struct)
// and implements the Logger interface.
// This is a private implementation detail for bridging logr to golib's Logger.
type logrAdapter struct {
	logger logr.Logger
}

func (a *logrAdapter) Enabled() bool {
	return a.logger.Enabled()
}

func (a *logrAdapter) Info(msg string, keysAndValues ...any) {
	a.logger.Info(msg, keysAndValues...)
}

func (a *logrAdapter) Error(err error, msg string, keysAndValues ...any) {
	a.logger.Error(err, msg, keysAndValues...)
}

func (a *logrAdapter) V(level int) Logger {
	return &logrAdapter{logger: a.logger.V(level)}
}

func (a *logrAdapter) WithValues(keysAndValues ...any) Logger {
	return &logrAdapter{logger: a.logger.WithValues(keysAndValues...)}
}

func (a *logrAdapter) WithName(name string) Logger {
	return &logrAdapter{logger: a.logger.WithName(name)}
}

func Discard() Logger {
	return FromLogr(logr.Discard())
}
