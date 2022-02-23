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
	"sync"

	"github.com/go-logr/logr"
)

var _ logr.Logger = &placeHolder{}

// placeHolder knows how to populate a concrete logr.Logger with
// options and propagate to its children.
// It use the logr.Discard Logger before actual Logger is fulfilled.
type placeHolder struct {
	logr.Logger

	name   string
	values []interface{}

	children []*placeHolder

	once     *sync.Once
	onceDone bool
}

func newPlaceHolderLogger() *placeHolder {
	return &placeHolder{
		Logger:   logr.Discard(),
		children: make([]*placeHolder, 0),
		once:     &sync.Once{},
	}
}

// Propagate switches the logger over to use the actual logger,
// instread of the discard logger, and propagates to all its children.
func (l *placeHolder) Propagate(actual logr.Logger) {
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

// WithName adds a new element to the logger's name.
// Successive calls with WithName continue to append
// suffixes to the logger's name.  It's strongly recommended
// that name segments contain only letters, digits, and hyphens
// (see the package documentation for more information).
func (l *placeHolder) WithName(name string) logr.Logger {
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
func (l *placeHolder) WithValues(kvs ...interface{}) logr.Logger {
	if l.onceDone {
		return l.Logger.WithValues(kvs)
	}

	child := newPlaceHolderLogger()
	child.values = append(child.values, kvs...)
	l.children = append(l.children, child)

	return child
}
