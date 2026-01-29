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

// Logger is a logging interface owned by golib.
// It is designed to be compatible with both logr v0.4.0 (interface) and v1.0.0+ (struct).
type Logger interface {
	// Enabled tests whether this Logger is enabled.
	Enabled() bool

	// Info logs a non-error message with the given key/value pairs as context.
	Info(msg string, keysAndValues ...any)

	// Error logs an error, with the given message and key/value pairs as context.
	Error(err error, msg string, keysAndValues ...any)

	// V returns a new Logger instance for a specific verbosity level.
	V(level int) Logger

	// WithValues returns a new Logger instance with additional key/value pairs.
	WithValues(keysAndValues ...any) Logger

	// WithName returns a new Logger instance with the specified name element added.
	WithName(name string) Logger
}
