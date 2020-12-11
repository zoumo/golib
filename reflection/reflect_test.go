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
