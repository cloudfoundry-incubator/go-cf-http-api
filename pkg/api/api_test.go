package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"path/filepath"

	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api"
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api/auth"
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/testhelpers"
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/uaaclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Api", func() {
	It("provides debug/metrics without authorization", func() {
		port, err := testhelpers.GetOpenPort()
		Expect(err).ToNot(HaveOccurred())

		server := api.New(api.Config{
			UAAClient: testhelpers.NewFakeUAAClient(),
			Port:      port,
		})
		stop := server.Start()
		defer stop()

		err = testhelpers.PollForUp(port)
		Expect(err).ToNot(HaveOccurred())

		resp, err := http.Get("http://localhost:" + port + "/debug/metrics")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()

		Expect(string(body)).ToNot(BeEmpty())

		By("expecting json output")
		var dataMap map[string]interface{}
		err = json.Unmarshal(body, &dataMap)
		Expect(err).ToNot(HaveOccurred())

		Expect(dataMap).To(HaveKey("memstats"))
	})

	It("returns 404 on unhandled endpoints", func() {
		port, err := testhelpers.GetOpenPort()
		Expect(err).ToNot(HaveOccurred())

		server := api.New(api.Config{
			UAAClient: testhelpers.NewFakeUAAClient(),
			Port:      port,
		})
		stop := server.Start()
		defer stop()

		err = testhelpers.PollForUp(port)
		Expect(err).ToNot(HaveOccurred())

		resp, err := http.Get("http://localhost:" + port + "/INVALID-ENDPOINT")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("allows custom Endpoints and defaults the content-type to application/json", func() {
		port, err := testhelpers.GetOpenPort()
		Expect(err).ToNot(HaveOccurred())

		type SomeType struct {
			Sent int `json:"sent"`
		}

		type SomeOutput struct {
			Received int    `json:"received"`
			Param    string `json:"param"`
			UserId   string `json:"userId"`
			PathVar  string `json:"pathVar"`
		}

		uaaClient := testhelpers.NewFakeUAAClient()
		uaaClient.SetUser(&uaaclient.User{ID: "some-id"})

		server := api.New(api.Config{
			UAAClient:  uaaClient,
			Port:       port,
			LogRequest: func(api.Request, api.Response, *api.Endpoint, time.Time, time.Duration) {},
			Endpoints: []*api.Endpoint{
				{
					Method: http.MethodPost,
					Path:   "/some-endpoint/{somePathVariable:[0-9]+}",
					Auth:   auth.LoggedIn,
					Handle: func(r api.Request) *api.Response {
						var t SomeType
						r.Decode(&t)

						param := r.GetParam("test")
						user := r.CurrentUser()

						return api.Ok(SomeOutput{
							Received: t.Sent,
							Param:    param,
							UserId:   user.ID,
							PathVar:  r.GetParam("somePathVariable"),
						})
					},
				},
			},
		})

		stop := server.Start()
		defer stop()

		err = testhelpers.PollForUp(port)
		Expect(err).ToNot(HaveOccurred())

		resp, err := http.Post("http://localhost:"+port+"/some-endpoint/222?test=value", "application/json", bytes.NewReader([]byte(`{
			"sent": 2
		}`)))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(body).To(MatchJSON(`{
			"received": 2,
			"param": "value",
			"userId": "some-id",
			"pathVar": "222"
		}`))
	})

	It("stops the server", func() {
		port, err := testhelpers.GetOpenPort()
		Expect(err).ToNot(HaveOccurred())

		server := api.New(api.Config{
			UAAClient: testhelpers.NewFakeUAAClient(),
			Port:      port,
		})
		stop := server.Start()

		err = testhelpers.PollForUp(port)
		Expect(err).ToNot(HaveOccurred())

		_, err = http.Get("http://localhost:" + port + "/v1/info")
		Expect(err).ToNot(HaveOccurred())

		stop()

		_, err = http.Get("http://localhost:" + port + "/v1/info")
		Expect(err).To(HaveOccurred())
	})

	Context("Auth", func() {
		It("requires user to be logged in with logged in auth type", func() {
			port, err := testhelpers.GetOpenPort()
			Expect(err).ToNot(HaveOccurred())

			uaaClient := testhelpers.NewFakeUAAClient()
			uaaClient.SetUser(&uaaclient.User{ID: "some-id"})

			logRequestCalled := make(chan interface{}, 10)
			fakeLogRequest := func(api.Request, api.Response, *api.Endpoint, time.Time, time.Duration) {
				logRequestCalled <- struct{}{}
			}

			called := make(chan interface{}, 10)
			server := api.New(api.Config{
				UAAClient:  uaaClient,
				Port:       port,
				LogRequest: fakeLogRequest,
				Endpoints: []*api.Endpoint{
					{
						Method: http.MethodGet,
						Path:   "/auth-endpoint",
						Auth:   auth.LoggedIn,
						Handle: func(r api.Request) *api.Response {
							called <- struct{}{}

							return api.Ok(nil)
						},
					},
				},
			})

			stop := server.Start()
			defer stop()

			err = testhelpers.PollForUp(port)
			Expect(err).ToNot(HaveOccurred())

			By("having no error returned by uaa")
			resp, err := http.Get("http://localhost:" + port + "/auth-endpoint")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Eventually(called).Should(Receive())
			Eventually(logRequestCalled).Should(Receive())

			By("having an error returned by uaa")
			uaaClient.SetUser(nil)

			resp, err = http.Get("http://localhost:" + port + "/auth-endpoint")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			Consistently(called).ShouldNot(Receive())
			Eventually(logRequestCalled).Should(Receive())
		})

		It("does not require user to be logged in with none auth type", func() {
			port, err := testhelpers.GetOpenPort()
			Expect(err).ToNot(HaveOccurred())

			uaaClient := testhelpers.NewFakeUAAClient()
			uaaClient.SetUser(&uaaclient.User{ID: "some-id"})

			logRequestCalled := make(chan interface{}, 10)
			fakeLogRequest := func(api.Request, api.Response, *api.Endpoint, time.Time, time.Duration) {
				logRequestCalled <- struct{}{}
			}

			called := make(chan interface{}, 10)
			server := api.New(api.Config{
				UAAClient:  uaaClient,
				Port:       port,
				LogRequest: fakeLogRequest,
				Endpoints: []*api.Endpoint{
					{
						Method: http.MethodGet,
						Path:   "/no-auth-endpoint",
						Auth:   auth.None,
						Handle: func(r api.Request) *api.Response {
							called <- struct{}{}

							return api.Ok(nil)
						},
					},
				},
			})

			stop := server.Start()
			defer stop()

			err = testhelpers.PollForUp(port)
			Expect(err).ToNot(HaveOccurred())

			By("having no error returned by uaa")
			resp, err := http.Get("http://localhost:" + port + "/no-auth-endpoint")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Eventually(called).Should(Receive())
			Eventually(logRequestCalled).Should(Receive())

			By("having an error returned by uaa")
			uaaClient.SetError(errors.New("uaa error"))

			resp, err = http.Get("http://localhost:" + port + "/no-auth-endpoint")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Eventually(called).Should(Receive())
			Eventually(logRequestCalled).Should(Receive())
		})
	})

	It("requires the user to have specific scopes and be logged in with scope auth type", func() {
		port, err := testhelpers.GetOpenPort()
		Expect(err).ToNot(HaveOccurred())

		uaaClient := testhelpers.NewFakeUAAClient()
		uaaClient.SetUser(&uaaclient.User{ID: "some-id", Scopes: []string{"some-scope"}})

		logRequestCalled := make(chan interface{}, 10)
		fakeLogRequest := func(api.Request, api.Response, *api.Endpoint, time.Time, time.Duration) {
			logRequestCalled <- struct{}{}
		}

		called := make(chan interface{}, 10)
		server := api.New(api.Config{
			UAAClient:  uaaClient,
			Port:       port,
			LogRequest: fakeLogRequest,
			Endpoints: []*api.Endpoint{
				{
					Method: http.MethodGet,
					Path:   "/auth-endpoint",
					Auth:   auth.AnyScope("some-scope"),
					Handle: func(r api.Request) *api.Response {
						called <- struct{}{}

						return api.Ok(nil)
					},
				},
			},
		})

		stop := server.Start()
		defer stop()

		err = testhelpers.PollForUp(port)
		Expect(err).ToNot(HaveOccurred())

		By("having no error returned by uaa and the correct scopes")
		resp, err := http.Get("http://localhost:" + port + "/auth-endpoint")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Eventually(called).Should(Receive())
		Eventually(logRequestCalled).Should(Receive())

		By("having no error returned by uaa and the incorrect scopes")
		uaaClient.SetUser(&uaaclient.User{ID: "some-id", Scopes: []string{"bad-scope"}})

		resp, err = http.Get("http://localhost:" + port + "/auth-endpoint")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		Consistently(called).ShouldNot(Receive())
		Eventually(logRequestCalled).Should(Receive())

		By("having no error returned by uaa and no scopes")
		uaaClient.SetUser(&uaaclient.User{ID: "some-id", Scopes: []string{}})

		resp, err = http.Get("http://localhost:" + port + "/auth-endpoint")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		Consistently(called).ShouldNot(Receive())
		Eventually(logRequestCalled).Should(Receive())

		By("having an error returned by uaa")
		uaaClient.SetError(errors.New("uaa error"))

		resp, err = http.Get("http://localhost:" + port + "/auth-endpoint")
		Expect(err).ToNot(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		Consistently(called).ShouldNot(Receive())
		Eventually(logRequestCalled).Should(Receive())
	})

	Context("HTMLTemplate Responses", func() {
		It("renders html templates with a hostname", func() {
			port, err := testhelpers.GetOpenPort()
			Expect(err).ToNot(HaveOccurred())

			uaaClient := testhelpers.NewFakeUAAClient()
			uaaClient.SetUser(&uaaclient.User{ID: "some-id"})

			templatesDir, err := filepath.Abs(".")
			Expect(err).ToNot(HaveOccurred())

			server := api.New(api.Config{
				UAAClient:          uaaClient,
				Hostname:           "example.com",
				TemplatesDirectory: templatesDir,
				Port:               port,
				LogRequest:         func(api.Request, api.Response, *api.Endpoint, time.Time, time.Duration) {},
				Endpoints: []*api.Endpoint{
					{
						Method: http.MethodGet,
						Path:   "/some-endpoint",
						Auth:   auth.None,
						Handle: func(r api.Request) *api.Response {
							return api.Ok(api.HTMLTemplate{Name: "test"})
						},
					},
				},
			})

			stop := server.Start()
			defer stop()

			err = testhelpers.PollForUp(port)
			Expect(err).ToNot(HaveOccurred())

			resp, err := http.Get("http://localhost:" + port + "/some-endpoint")
			Expect(err).ToNot(HaveOccurred())

			Expect(resp.Header.Get("Content-Type")).To(Equal("text/html"))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(body)).To(ContainSubstring("the hostname is example.com"))
		})

		It("returns a server error if the template doesn't exist", func() {
			port, err := testhelpers.GetOpenPort()
			Expect(err).ToNot(HaveOccurred())

			uaaClient := testhelpers.NewFakeUAAClient()
			uaaClient.SetUser(&uaaclient.User{ID: "some-id"})

			templatesDir, err := filepath.Abs(".")
			Expect(err).ToNot(HaveOccurred())

			server := api.New(api.Config{
				UAAClient:          uaaClient,
				Hostname:           "example.com",
				TemplatesDirectory: templatesDir,
				Port:               port,
				LogRequest:         func(api.Request, api.Response, *api.Endpoint, time.Time, time.Duration) {},
				Endpoints: []*api.Endpoint{
					{
						Method: http.MethodGet,
						Path:   "/some-endpoint",
						Auth:   auth.None,
						Handle: func(r api.Request) *api.Response {
							return api.Ok(api.HTMLTemplate{Name: "non-existent"})
						},
					},
				},
			})

			stop := server.Start()
			defer stop()

			err = testhelpers.PollForUp(port)
			Expect(err).ToNot(HaveOccurred())

			resp, err := http.Get("http://localhost:" + port + "/some-endpoint")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("Asset Responses", func() {
		It("renders assets", func() {
			port, err := testhelpers.GetOpenPort()
			Expect(err).ToNot(HaveOccurred())

			uaaClient := testhelpers.NewFakeUAAClient()
			uaaClient.SetUser(&uaaclient.User{ID: "some-id"})

			assetsDir, err := filepath.Abs(".")
			Expect(err).ToNot(HaveOccurred())

			server := api.New(api.Config{
				UAAClient:       uaaClient,
				Hostname:        "example.com",
				AssetsDirectory: assetsDir,
				Port:            port,
				LogRequest:      func(api.Request, api.Response, *api.Endpoint, time.Time, time.Duration) {},
				Endpoints:       []*api.Endpoint{},
			})

			stop := server.Start()
			defer stop()

			err = testhelpers.PollForUp(port)
			Expect(err).ToNot(HaveOccurred())

			resp, err := http.Get("http://localhost:" + port + "/assets/test_asset.png")
			Expect(err).ToNot(HaveOccurred())

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(body).ToNot(BeNil())
			Expect(resp.Header.Get("Content-Type")).To(Equal("image/png"))
		})
	})
})
