module github.com/zoumo/golib

go 1.15

require (
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0
	github.com/jcmturner/gokrb5/v8 v8.4.2
	github.com/keybase/go-ps v0.0.0-20190827175125-91aafc93ba19
	github.com/moby/moby v1.13.1
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.5.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1
	golang.org/x/text v0.3.6
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v0.20.0
	k8s.io/klog/v2 v2.10.0
)

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr => github.com/go-logr/zapr v0.4.0
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.10.0
)
