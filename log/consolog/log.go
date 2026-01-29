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

package consolog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/term"

	"github.com/zoumo/golib/log"
)

const (
	resetColor = "\033[0m"
)

// These constants identify the log levels in order of increasing severity.
// A message written to a high-severity log file is also written to each
// lower-severity log file.
const (
	infoLog int = iota
	errorLog
)

var (
	verbose         = 0
	enableColor     = true
	consoleColorMap = map[string]string{
		"blue":   "\033[34m",
		"green":  "\033[32m",
		"red":    "\033[31m",
		"yellow": "\033[33m",
		"strong": "\033[1m",
	}
)

func init() {
	enableColor = term.IsTerminal(int(os.Stdout.Fd()))
}

func New() log.Logger {
	return &logger{
		level:       0,
		enableColor: enableColor,
		prefix:      "",
		values:      nil,
	}
}

// InitFlags is for explicitly initializing the flags.
func InitFlags(flagset *pflag.FlagSet) {
	pflag.IntVar(&verbose, "v", verbose, "number for the log level verbosity")
	pflag.BoolVar(&enableColor, "color", enableColor, "enable color logging")
}

type logger struct {
	level       int
	enableColor bool
	prefix      string
	values      []any
}

func copySlice(in []any) []any {
	out := make([]any, len(in))
	copy(out, in)
	return out
}

func (l *logger) clone() *logger {
	return &logger{
		level:       l.level,
		enableColor: l.enableColor,
		prefix:      l.prefix,
		values:      copySlice(l.values),
	}
}

func (l *logger) getColor(color string) string {
	if !l.enableColor {
		return ""
	}
	consoleColor, ok := consoleColorMap[color]
	if !ok {
		return ""
	}
	return consoleColor
}

func (l *logger) Enabled() bool {
	return verbose >= l.level
}

func (l *logger) V(level int) log.Logger {
	new := l.clone()
	new.level = level
	return new
}

func (l *logger) WithName(name string) log.Logger {
	new := l.clone()
	if len(l.prefix) > 0 {
		new.prefix = l.prefix + "/"
	}
	new.prefix += name
	return new
}

func (l *logger) WithValues(kvList ...any) log.Logger {
	new := l.clone()
	new.values = append(new.values, kvList...)
	return new
}

func (l *logger) Info(msg string, keysAndValues ...any) {
	if !l.Enabled() {
		return
	}
	trimmed := trimDuplicates(l.values, keysAndValues)
	kvList := []any{}
	for i := range trimmed {
		kvList = append(kvList, trimmed[i]...)
	}
	l.print(infoLog, msg, kvList)
}

func (l *logger) Error(err error, msg string, keysAndValues ...any) {
	if !l.Enabled() {
		return
	}
	trimmed := trimDuplicates(l.values, keysAndValues)
	kvList := []any{}
	for i := range trimmed {
		kvList = append(kvList, trimmed[i]...)
	}
	var loggableErr any
	if err != nil {
		loggableErr = err.Error()
	}
	kvList = append(kvList, "ERROR", loggableErr)
	l.print(errorLog, msg, kvList)
}

func (l *logger) print(level int, msg string, kvList []any) {
	buf := &bytes.Buffer{}
	l.printTime(level, buf)

	if len(l.prefix) > 0 {
		buf.WriteString(" ")
		l.printPrefix(buf)
	}
	buf.WriteString(" ")
	l.printMsg(buf, msg)
	buf.WriteString("\n")
	l.printKV(buf, kvList...)

	fmt.Print(buf.String())
}

func (l *logger) printTime(level int, buf io.Writer) {
	reset := resetColor
	var color string
	if level == infoLog {
		color = l.getColor("blue")
	} else {
		color = l.getColor("red")
	}
	if color == "" {
		reset = ""
	}

	fmt.Fprintf(buf, "%s==> [%s]%s", color, time.Now().Format(time.RFC3339), reset) //nolint
}

func (l *logger) printPrefix(buf io.Writer) {
	reset := resetColor
	green := l.getColor("green")
	if green == "" {
		reset = ""
	}
	buf.Write([]byte(green + l.prefix + reset)) //nolint
}

func (l *logger) printMsg(buf io.Writer, msg string) {
	reset := resetColor
	strong := l.getColor("strong")
	if strong == "" {
		reset = ""
	}
	buf.Write([]byte(strong + msg + reset)) //nolint
}

func (l *logger) printKV(buf io.Writer, kvList ...any) {
	reset := resetColor
	color := l.getColor("yellow")
	if color == "" {
		reset = ""
	}
	keyMaxLen := 0
	keys := make([]string, 0, len(kvList))
	vals := make(map[string]any, len(kvList))
	for i := 0; i < len(kvList); i += 2 {
		k, ok := kvList[i].(string)
		if !ok {
			panic(fmt.Sprintf("key is not a string: %s", pretty(kvList[i])))
		}
		var v any
		if i+1 < len(kvList) {
			v = kvList[i+1]
		}
		keys = append(keys, k)
		vals[k] = v
		if len(k) > keyMaxLen {
			keyMaxLen = len(k)
		}
	}
	sort.Strings(keys)
	// nolint
	for _, k := range keys {
		v := vals[k]
		buf.Write([]byte("    "))
		format := fmt.Sprintf("%%s%%-%ds%%s", keyMaxLen)
		fmt.Fprintf(buf, format, color, k, reset)
		buf.Write([]byte(" = "))
		buf.Write([]byte(pretty(v)))
		buf.Write([]byte("\n"))
	}
}

// trimDuplicates will deduplicates elements provided in multiple KV tuple
// slices, whilst maintaining the distinction between where the items are
// contained.
func trimDuplicates(kvLists ...[]any) [][]any {
	// maintain a map of all seen keys
	seenKeys := map[any]struct{}{}
	// build the same number of output slices as inputs
	outs := make([][]any, len(kvLists))
	// iterate over the input slices backwards, as 'later' kv specifications
	// of the same key will take precedence over earlier ones
	for i := len(kvLists) - 1; i >= 0; i-- {
		// initialize this output slice
		outs[i] = []any{}
		// obtain a reference to the kvList we are processing
		kvList := kvLists[i]

		// start iterating at len(kvList) - 2 (i.e. the 2nd last item) for
		// slices that have an even number of elements.
		// We add (len(kvList) % 2) here to handle the case where there is an
		// odd number of elements in a kvList.
		// If there is an odd number, then the last element in the slice will
		// have the value 'null'.
		for i2 := len(kvList) - 2 + (len(kvList) % 2); i2 >= 0; i2 -= 2 {
			k := kvList[i2]
			// if we have already seen this key, do not include it again
			if _, ok := seenKeys[k]; ok {
				continue
			}
			// make a note that we've observed a new key
			seenKeys[k] = struct{}{}
			// attempt to obtain the value of the key
			var v any
			// i2+1 should only ever be out of bounds if we handling the first
			// iteration over a slice with an odd number of elements
			if i2+1 < len(kvList) {
				v = kvList[i2+1]
			}
			// add this KV tuple to the *start* of the output list to maintain
			// the original order as we are iterating over the slice backwards
			outs[i] = append([]any{k, v}, outs[i]...)
		}
	}
	return outs
}

func pretty(value any) string {
	if err, ok := value.(error); ok {
		if _, ok := value.(json.Marshaler); !ok {
			value = err.Error()
		}
	}
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(value)                     //nolint
	return strings.TrimSpace(buffer.String()) //nolint
}
