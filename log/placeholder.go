// Copyright 2022 jim.zoumo@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package log

import (
	"sync"
)

var _ Logger = &placeHolder{}

// placeHolder knows how to populate a concrete Logger with
// options and propagate to its children.
// It uses the noop Logger before actual Logger is fulfilled.
type placeHolder struct {
	Logger

	name   string
	values []any

	children []*placeHolder

	once     *sync.Once
	onceDone bool
}

func newPlaceHolderLogger() *placeHolder {
	return &placeHolder{
		Logger:   Discard(),
		children: make([]*placeHolder, 0),
		once:     &sync.Once{},
	}
}

// Propagate switches the logger over to use the actual logger,
// instead of the noop logger, and propagates to all its children.
func (l *placeHolder) Propagate(actual Logger) {
	l.once.Do(func() {
		logger := actual
		if len(l.name) > 0 {
			logger = logger.WithName(l.name)
		}
		if len(l.values) > 0 {
			logger = logger.WithValues(l.values...)
		}

		l.Logger = logger

		for _, child := range l.children {
			child.Propagate(logger)
		}
		l.onceDone = true
	})
}

// Enabled tests whether this Logger is enabled.
func (l *placeHolder) Enabled() bool {
	return l.Logger.Enabled()
}

// Info logs a non-error message with the given key/value pairs as context.
func (l *placeHolder) Info(msg string, keysAndValues ...interface{}) {
	l.Logger.Info(msg, keysAndValues...)
}

// Error logs an error, with the given message and key/value pairs as context.
func (l *placeHolder) Error(err error, msg string, keysAndValues ...interface{}) {
	l.Logger.Error(err, msg, keysAndValues...)
}

// V returns a new Logger instance for a specific verbosity level.
func (l *placeHolder) V(level int) Logger {
	if l.onceDone {
		return l.Logger.V(level)
	}

	child := newPlaceHolderLogger()
	child.name = l.name
	child.values = append(child.values, l.values...)
	l.children = append(l.children, child)
	return child
}

// WithName adds a new element to the logger's name.
// Successive calls with WithName continue to append
// suffixes to the logger's name.  It's strongly recommended
// that name segments contain only letters, digits, and hyphens
// (see the package documentation for more information).
func (l *placeHolder) WithName(name string) Logger {
	if l.onceDone {
		return l.Logger.WithName(name)
	}

	child := newPlaceHolderLogger()
	child.name = name
	l.children = append(l.children, child)

	return child
}

// WithValues adds some key-value pairs of context to a logger.
// See Info for documentation on how key/value pairs work.
func (l *placeHolder) WithValues(kvs ...interface{}) Logger {
	if l.onceDone {
		return l.Logger.WithValues(kvs)
	}

	child := newPlaceHolderLogger()
	child.values = append(child.values, kvs...)
	l.children = append(l.children, child)

	return child
}
