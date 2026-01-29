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
	"github.com/go-logr/logr"
)

// SetLogrLogger sets a concrete logging implementation for all deferred Loggers.
// Accepts a logr.Logger (works with both v0.4.0 and v1.0.0+).
func SetLogrLogger(l logr.Logger) {
	singleton.Propagate(FromLogr(l))
}

func SetLogger(l Logger) {
	singleton.Propagate(l)
}

var (
	singleton = newPlaceHolderLogger()
	// Log is the base logger. It delegates to another Logger in
	// place holder. You must call SetLogger to get any actual logging.
	Log Logger = singleton
)
