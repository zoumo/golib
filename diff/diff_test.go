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

import "github.com/stretchr/testify/suite"

type diffTestSuite struct {
	suite.Suite

	impl        Diff
	implColored Diff
}

func (s *diffTestSuite) SetupSuite() {
	s.impl = New()
	s.implColored = New(WithColored())
}

func (s *diffTestSuite) TestDiffUnified() {
	tests := []struct {
		name string
		old  string
		new  string
		want string
	}{
		{
			name: "equal",
			old:  "Hello, World!",
			new:  "Hello, World!",
			want: ``,
		},
		{
			name: "diff",
			old:  "header\nHello, World!\ntail\n",
			new:  "header\nHello, Go!\ntail\n",
			want: "--- old\n+++ new\n@@ -1,3 +1,3 @@\n header\n-Hello, World!\n+Hello, Go!\n tail\n",
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			diff := s.impl.DiffUnified("old", tt.old, "new", tt.new)
			s.Equal(tt.want, diff)
		})
	}
}

func (s *diffTestSuite) TestDiffUnified_Colored() {
	tests := []struct {
		name string
		old  string
		new  string
		want string
	}{
		{
			name: "equal",
			old:  "Hello, World!",
			new:  "Hello, World!",
			want: ``,
		},
		{
			name: "diff",
			old:  "header\nHello, World!\ntail\n",
			new:  "header\nHello, Go!\ntail\n",
			want: "\x1b[32m--- old\x1b[0m\n\x1b[31m+++ new\x1b[0m\n@@ -1,3 +1,3 @@\n header\n\x1b[32m-Hello, World!\x1b[0m\n\x1b[31m+Hello, Go!\x1b[0m\n tail\n",
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			diff := s.implColored.DiffUnified("old", "new", tt.old, tt.new)
			s.Equal(tt.want, diff)
		})
	}
}
