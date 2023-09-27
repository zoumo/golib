// Copyright 2023 jim.zoumo@gmail.com
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

package exec

import (
	"bytes"
	"io"
)

// writerWithBuffer warps a writer with buffer
// so you can read bytes from the buffer
type writerWithBuffer struct {
	buffer *bytes.Buffer
	w      io.Writer
}

func newWriterWithBuffer(w io.Writer) io.ReadWriter {
	buf := new(bytes.Buffer)
	mwr := &writerWithBuffer{
		buffer: buf,
		w:      buf,
	}
	if w != nil {
		mwr.w = io.MultiWriter(w, buf)
	}
	return mwr
}

func (mwr *writerWithBuffer) Write(p []byte) (n int, err error) {
	return mwr.w.Write(p)
}

func (mwr *writerWithBuffer) Read(p []byte) (n int, err error) {
	return mwr.buffer.Read(p)
}
