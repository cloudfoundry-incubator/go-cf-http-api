package httpclient_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/httpclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTPClient", func() {
	It("creates a get request object", func() {
		c, err := httpclient.New("http://localhost", true)
		Expect(err).ToNot(HaveOccurred())

		req, err := c.GetRequest("/endpoint")
		Expect(err).ToNot(HaveOccurred())

		Expect(req.URL.String()).To(Equal("http://localhost/endpoint"))
		Expect(req.Method).To(Equal(http.MethodGet))
	})

	It("creates a post request object", func() {
		c, err := httpclient.New("http://localhost", true)
		Expect(err).ToNot(HaveOccurred())

		req, err := c.PostRequest("/endpoint", strings.NewReader("some-body"))
		Expect(err).ToNot(HaveOccurred())

		Expect(req.URL.String()).To(Equal("http://localhost/endpoint"))
		Expect(req.Method).To(Equal(http.MethodPost))

		output, err := ioutil.ReadAll(req.Body)
		Expect(err).ToNot(HaveOccurred())
		defer req.Body.Close()

		Expect(string(output)).To(Equal("some-body"))
	})

	It("creates a put request object", func() {
		c, err := httpclient.New("http://localhost", true)
		Expect(err).ToNot(HaveOccurred())

		req, err := c.PutRequest("/endpoint", strings.NewReader("some-body"))
		Expect(err).ToNot(HaveOccurred())

		Expect(req.URL.String()).To(Equal("http://localhost/endpoint"))
		Expect(req.Method).To(Equal(http.MethodPut))

		output, err := ioutil.ReadAll(req.Body)
		Expect(err).ToNot(HaveOccurred())
		defer req.Body.Close()

		Expect(string(output)).To(Equal("some-body"))
	})

	It("creates a delete request object", func() {
		c, err := httpclient.New("http://localhost", true)
		Expect(err).ToNot(HaveOccurred())

		req, err := c.DeleteRequest("/endpoint", strings.NewReader("some-body"))
		Expect(err).ToNot(HaveOccurred())

		Expect(req.URL.String()).To(Equal("http://localhost/endpoint"))
		Expect(req.Method).To(Equal(http.MethodDelete))

		output, err := ioutil.ReadAll(req.Body)
		Expect(err).ToNot(HaveOccurred())
		defer req.Body.Close()

		Expect(string(output)).To(Equal("some-body"))
	})

	It("makes a request with a context", func() {
		requests := make(chan *http.Request, 1)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests <- r

			fmt.Fprint(w, `{"some": "json"}`)
		}))

		c, err := httpclient.New(mockServer.URL, true)
		Expect(err).ToNot(HaveOccurred())

		req, err := c.GetRequest("/endpoint")
		Expect(err).ToNot(HaveOccurred())

		resp, err := c.Do(req, context.Background())
		Expect(err).ToNot(HaveOccurred())

		Expect(requests).To(HaveLen(1))
		request := <-requests
		Expect(request.Method).To(Equal(http.MethodGet))

		output, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()
		Expect(output).To(MatchJSON(`{"some": "json"}`))
	})

	It("returns an error if the host is empty", func() {
		_, err := httpclient.New("", true)
		Expect(err).To(MatchError("host cannot be empty"))
	})

	It("returns an error if the url is bad", func() {
		_, err := httpclient.New("://", true)
		Expect(err).To(HaveOccurred())
	})
})
