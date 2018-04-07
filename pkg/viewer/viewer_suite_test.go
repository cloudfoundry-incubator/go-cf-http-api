package viewer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestViewer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Viewer Suite")
}
