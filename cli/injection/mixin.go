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

package injection

import (
	"github.com/zoumo/golib/log"
)

var _ RequiresLogger = &DefaultInjectionMixin{}
var _ RequiresWorkspace = &DefaultInjectionMixin{}

type DefaultInjectionMixin struct {
	Logger    log.Logger
	Workspace string
}

func NewDefaultInjectionMixin() *DefaultInjectionMixin {
	return &DefaultInjectionMixin{}
}

func (m *DefaultInjectionMixin) InjectLogger(logger log.Logger) {
	m.Logger = logger
}

func (m *DefaultInjectionMixin) InjectWorkspace(ws string) {
	m.Workspace = ws
}
