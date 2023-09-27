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

package krb5

import (
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/pkg/errors"
)

func GetDefaultPrinciplaNameFromCCache(ccache string) (string, error) {
	if ccache == "" {
		return "", errors.New("invalid ccache file path")
	}
	c, err := credentials.LoadCCache(ccache)
	if err != nil {
		return "", errors.Wrap(err, "failed to load kerberos ccache")
	}
	if len(c.DefaultPrincipal.PrincipalName.NameString) == 0 {
		return "", errors.New("empty principalName")
	}
	return c.DefaultPrincipal.PrincipalName.NameString[0], nil
}
