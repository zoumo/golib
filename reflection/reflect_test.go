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

package reflection

import (
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"
)

func init() {
	config.DefaultReporterConfig.NoColor = true
}

func TestPorterSuit(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Porter suit")
}

//nolint
var (
	Describe       = ginkgo.Describe
	DescribeTable  = table.DescribeTable
	Entry          = table.Entry
	Context        = ginkgo.Context
	It             = ginkgo.It
	BeforeEach     = ginkgo.BeforeEach
	AfterEach      = ginkgo.AfterEach
	JustBeforeEach = ginkgo.JustBeforeEach
	JustAfterEach  = ginkgo.JustAfterEach
	BeforeSuite    = ginkgo.BeforeSuite
	AfterSuit      = ginkgo.AfterSuite
	Fail           = ginkgo.Fail
	Skip           = ginkgo.Skip
	Expect         = gomega.Expect
	Equal          = gomega.Equal
	BeNil          = gomega.BeNil
	BeTrue         = gomega.BeTrue
)

var _ = BeforeEach(func() {
})
