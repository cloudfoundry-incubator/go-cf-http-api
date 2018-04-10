package uaaclient_test

import (
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/uaaclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAAClient", func() {
	It("returns an error if the internal client cannot be created", func() {
		_, err := uaaclient.New("", true, "client_id", "client_secret")
		Expect(err).To(HaveOccurred())
	})

	It("returns an error if client_id and/or client_secret are blank", func() {
		_, err := uaaclient.New("http://some-url.com", true, "", "some-secret")
		Expect(err).To(HaveOccurred())

		_, err = uaaclient.New("http://some-url.com", true, "some-client", "")
		Expect(err).To(HaveOccurred())

		_, err = uaaclient.New("http://some-url.com", true, "", "")
		Expect(err).To(HaveOccurred())
	})
})
