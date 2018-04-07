package errors_test

import (
	. "../errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Errors", func() {
	It("returns a list of descriptions", func() {
		errorList := &ErrorListResponse{Errors: []ErrorResponse{
			{"first error"},
			{"second error"}},
		}

		Expect(errorList.GetErrors()).To(ConsistOf("first error", "second error"))
	})
})
