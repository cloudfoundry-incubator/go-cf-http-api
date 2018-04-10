package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/logger"
)

func TestApi(t *testing.T) {
	logger.Init(GinkgoWriter, GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
