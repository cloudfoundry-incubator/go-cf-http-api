package evaclient_test

import (
	"../evaclient"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"../api"
	"../api/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EVAClient", func() {
	It("returns an error if the internal client cannot be created", func() {
		_, err := evaclient.New("", true, "")
		Expect(err).To(HaveOccurred())
	})

	Context("Get", func() {
		It("calls the endpoint, returns a response and the result", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/info", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				json.NewEncoder(w).Encode(api.InfoResponse{
					Commit: "commit-sha",
				})
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Get("/v1/info", &info)

			Expect(resp.Ok).To(BeTrue())
			Expect(info.Commit).To(Equal("commit-sha"))

			Expect(receivedRequest.Method).To(Equal(http.MethodGet))
			Expect(receivedRequest.Header.Get("Authorization")).To(Equal("bearer token"))
			Expect(receivedRequest.URL.Path).To(Equal("/v1/info"))
		})

		It("deserializes Api Errors on a non-2xx status", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/publishers", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(errors.ErrorListResponse{
					Errors: []errors.ErrorResponse{
						{Description: "validation failed"},
					},
				})
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Get("/v1/publishers", &info)

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).To(Equal([]errors.ErrorResponse{
				{Description: "validation failed"},
			}))

			Expect(receivedRequest.Method).To(Equal(http.MethodGet))
			Expect(receivedRequest.Header.Get("Authorization")).To(Equal("bearer token"))
			Expect(receivedRequest.URL.Path).To(Equal("/v1/publishers"))
		})

		It("returns a generic error on a non-2xx status without a body", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/publishers", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				w.WriteHeader(http.StatusUnprocessableEntity)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Get("/v1/publishers", &info)

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).To(Equal([]errors.ErrorResponse{
				{Description: "No errors specified. Received status 422"},
			}))
		})

		It("returns an error when it can't reach the server", func() {
			client, err := evaclient.New("http://localhost:12345", true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Get("/v1/publishers", &info)

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).ToNot(BeNil())
		})
	})

	Context("Post", func() {
		It("calls the endpoint, returns a response and the result", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			var receivedBody string
			mux.HandleFunc("/v1/info", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				bytes, err := ioutil.ReadAll(receivedRequest.Body)
				Expect(err).ToNot(HaveOccurred())

				receivedBody = string(bytes)

				w.WriteHeader(http.StatusOK)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			info := api.InfoResponse{
				URL:     "api-url",
				Version: "api-version",
				Commit:  "commit-sha",
			}
			resp := client.Post("/v1/info", &info)

			Expect(resp.Ok).To(BeTrue())

			Expect(receivedBody).To(MatchJSON(`{
				"url": "api-url",
				"version": "api-version",
				"commit": "commit-sha"
			}`))
			Expect(receivedRequest.Method).To(Equal(http.MethodPost))
			Expect(receivedRequest.Header.Get("Authorization")).To(Equal("bearer token"))
			Expect(receivedRequest.URL.Path).To(Equal("/v1/info"))
		})

		It("deserializes Api Errors on a non-2xx status", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/publishers", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(errors.ErrorListResponse{
					Errors: []errors.ErrorResponse{
						{Description: "validation failed"},
					},
				})
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Post("/v1/publishers", &info)

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).To(Equal([]errors.ErrorResponse{
				{Description: "validation failed"},
			}))

			Expect(receivedRequest.Method).To(Equal(http.MethodPost))
			Expect(receivedRequest.Header.Get("Authorization")).To(Equal("bearer token"))
			Expect(receivedRequest.URL.Path).To(Equal("/v1/publishers"))
		})

		It("returns a generic error on a non-2xx status without a body", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/publishers", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				w.WriteHeader(http.StatusUnprocessableEntity)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Post("/v1/publishers", &info)

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).To(Equal([]errors.ErrorResponse{
				{Description: "No errors specified. Received status 422"},
			}))
		})

		It("returns an error when it can't reach the server", func() {
			client, err := evaclient.New("http://localhost:12345", true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Post("/v1/publishers", &info)

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).ToNot(BeNil())
		})
	})

	Context("Put", func() {
		It("calls the endpoint, returns a response", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			var receivedBody string
			mux.HandleFunc("/v1/targets/1234", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				bytes, err := ioutil.ReadAll(receivedRequest.Body)
				Expect(err).ToNot(HaveOccurred())

				receivedBody = string(bytes)

				w.WriteHeader(http.StatusOK)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			type serializeMe struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			t := serializeMe{
				ID:   "1234",
				Name: "newTargetName",
			}
			resp := client.Put("/v1/targets/1234", &t)

			Expect(resp.Ok).To(BeTrue())

			Expect(receivedBody).To(MatchJSON(`{ "id": "1234", "name": "newTargetName"}`))
			Expect(receivedRequest.Method).To(Equal(http.MethodPut))
			Expect(receivedRequest.Header.Get("Authorization")).To(Equal("bearer token"))
			Expect(receivedRequest.URL.Path).To(Equal("/v1/targets/1234"))
		})

		It("deserializes Api Errors on a non-2xx status", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/targets", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(errors.ErrorListResponse{
					Errors: []errors.ErrorResponse{
						{Description: "validation failed"},
					},
				})
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Put("/v1/targets", &info)

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).To(Equal([]errors.ErrorResponse{
				{Description: "validation failed"},
			}))

			Expect(receivedRequest.Method).To(Equal(http.MethodPut))
			Expect(receivedRequest.Header.Get("Authorization")).To(Equal("bearer token"))
			Expect(receivedRequest.URL.Path).To(Equal("/v1/targets"))
		})

		It("returns a generic error on a non-2xx status without a body", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/targets", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				w.WriteHeader(http.StatusUnprocessableEntity)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Put("/v1/targets", &info)

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).To(Equal([]errors.ErrorResponse{
				{Description: "No errors specified. Received status 422"},
			}))
		})

		It("returns an error when it can't reach the server", func() {
			client, err := evaclient.New("http://localhost:12345", true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			var info api.InfoResponse
			resp := client.Put("/v1/publishers", &info)

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).ToNot(BeNil())
		})
	})

	Context("Delete", func() {
		It("calls the endpoint, returns a response", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/targets/12345", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				w.WriteHeader(http.StatusOK)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			resp := client.Delete("/v1/targets/12345")

			Expect(resp.Ok).To(BeTrue())

			Expect(receivedRequest.Method).To(Equal(http.MethodDelete))
			Expect(receivedRequest.Header.Get("Authorization")).To(Equal("bearer token"))
			Expect(receivedRequest.URL.Path).To(Equal("/v1/targets/12345"))
		})

		It("deserializes Api Errors on a non-2xx status", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/targets/123", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(errors.ErrorListResponse{
					Errors: []errors.ErrorResponse{
						{Description: "delete failed"},
					},
				})
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			resp := client.Delete("/v1/targets/123")

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).To(Equal([]errors.ErrorResponse{
				{Description: "delete failed"},
			}))

			Expect(receivedRequest.Method).To(Equal(http.MethodDelete))
			Expect(receivedRequest.Header.Get("Authorization")).To(Equal("bearer token"))
			Expect(receivedRequest.URL.Path).To(Equal("/v1/targets/123"))
		})

		It("returns a generic error on a non-2xx status without a body", func() {
			mux := http.NewServeMux()
			var receivedRequest *http.Request
			mux.HandleFunc("/v1/targets/123", func(w http.ResponseWriter, r *http.Request) {
				receivedRequest = r

				w.WriteHeader(http.StatusUnprocessableEntity)
			})

			ts := httptest.NewServer(mux)
			defer ts.Close()

			client, err := evaclient.New(ts.URL, true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			resp := client.Delete("/v1/targets/123")

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).To(Equal([]errors.ErrorResponse{
				{Description: "No errors specified. Received status 422"},
			}))
		})

		It("returns an error when it can't reach the server", func() {
			client, err := evaclient.New("http://localhost:12345", true, "bearer token")
			Expect(err).ToNot(HaveOccurred())

			resp := client.Delete("/v1/targets/123")

			Expect(resp.Ok).To(BeFalse())
			Expect(resp.Errors).ToNot(BeNil())
		})
	})
})
