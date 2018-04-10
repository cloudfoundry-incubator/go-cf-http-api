package api_test

import (
	"net/http"

	. "github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api"
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EndpointMatcher", func() {
	It("returns error if unable to decipher api.Endpoint", func() {
		result, err := BeCompleteEndpoint().Match("something")
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeFalse())
	})

	It("succeeds if the actual Endpoint has all the fields", func() {
		e := Endpoint{
			Path:   "/v1/info",
			Method: http.MethodGet,
			Auth:   auth.None,
			Handle: func(r Request) *Response { return nil },
		}

		result, err := BeCompleteEndpoint().Match(e)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeTrue())
	})

	It("succeeds if the actual Endpoint is a pointer", func() {
		e := &Endpoint{
			Path:   "/v1/info",
			Method: http.MethodGet,
			Auth:   auth.None,
			Handle: func(r Request) *Response { return nil },
		}

		result, err := BeCompleteEndpoint().Match(e)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeTrue())
	})

	It("returns error if path is empty", func() {
		e := Endpoint{
			Path:   "",
			Method: http.MethodGet,
			Auth:   auth.None,
			Handle: func(r Request) *Response { return nil },
		}

		result, err := BeCompleteEndpoint().Match(e)
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeFalse())
	})

	It("returns error if method is empty", func() {
		e := Endpoint{
			Path:   "/v1/info",
			Method: "",
			Auth:   auth.None,
			Handle: func(r Request) *Response { return nil },
		}

		result, err := BeCompleteEndpoint().Match(e)
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeFalse())
	})

	It("returns error if handle is nil", func() {
		e := Endpoint{
			Path:   "/v1/info",
			Method: http.MethodGet,
			Auth:   auth.None,
			Handle: nil,
		}

		result, err := BeCompleteEndpoint().Match(e)
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeFalse())
	})

	It("returns error if Auth is nil", func() {
		e := Endpoint{
			Path:   "/v1/info",
			Method: http.MethodGet,
			Auth:   nil,
			Handle: func(r Request) *Response { return nil },
		}

		result, err := BeCompleteEndpoint().Match(e)
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeFalse())
	})
})
