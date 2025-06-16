// Copyright 2025 The jim.zoumo@gmail.com
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

package diff

import (
	"fmt"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

const (
	greenColor = "\x1b[32m"
	redColor   = "\x1b[31m"
	resetColor = "\x1b[0m"
)

// Diff is an interface for computing diffs.
type Diff interface {
	DiffUnified(name1, name2, text1, text2 string) string
}

func New(options ...DiffOption) Diff {
	d := &diffImpl{}
	for _, option := range options {
		option.applyOptions(&d.options)
	}
	return d
}

type diffImpl struct {
	options DiffOptions
}

func (d *diffImpl) DiffUnified(name1, name2, text1, text2 string) string {
	edits := myers.ComputeEdits(span.URIFromPath(""), text1, text2)
	unified := gotextdiff.ToUnified(name1, name2, text1, edits)

	diffStr := fmt.Sprint(unified)

	if !d.options.Colored {
		return diffStr
	}

	lines := strings.Split(diffStr, "\n")
	// print color
	for i, line := range lines {
		if len(line) > 1 {
			if line[0] == '+' {
				// print green line in terminal
				lines[i] = fmt.Sprintf("%s%s%s", redColor, line, resetColor)
			} else if line[0] == '-' {
				// print red line in terminal
				lines[i] = fmt.Sprintf("%s%s%s", greenColor, line, resetColor)
			}
		}
	}
	return strings.Join(lines, "\n")
}
