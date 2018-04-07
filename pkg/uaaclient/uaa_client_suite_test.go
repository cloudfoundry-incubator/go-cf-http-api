package uaaclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUAA(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UAAClient Suite")
}
