package errors_test

import (
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Errors", func() {
	It("returns a list of descriptions", func() {
		errorList := &errors.ErrorListResponse{Errors: []errors.ErrorResponse{
			{"first error"},
			{"second error"}},
		}

		Expect(errorList.GetErrors()).To(ConsistOf("first error", "second error"))
	})
})
