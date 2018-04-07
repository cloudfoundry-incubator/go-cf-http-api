package api_test

import (
	"net/http"

	"../api"
	"../uaaclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("API - Info", func() {
	It("returns the api version", func() {
		info := &api.InfoResponse{
			URL:     "api-url",
			Version: "api-version",
			Commit:  "commit-hash",
		}

		resp := api.InfoEndpoint(info).Handle(&api.FakeRequest{
			User:   uaaclient.User{ID: "user-1"},
			Params: map[string]string{"user_guid": "user-1"},
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(resp.Body).To(Equal(info))
	})

	It("is a complete endpoint", func() {
		Expect(api.InfoEndpoint(nil)).To(api.BeCompleteEndpoint())
	})
})
