// Copyright 2025 The jim.zoumo@gmail.com
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

package diff

type DiffOption interface {
	applyOptions(*DiffOptions)
}

type diffOptionFunc func(*DiffOptions)

func (fn diffOptionFunc) applyOptions(options *DiffOptions) {
	fn(options)
}

func WithColored() DiffOption {
	return diffOptionFunc(func(do *DiffOptions) {
		do.Colored = true
	})
}

type DiffOptions struct {
	Colored bool
}
