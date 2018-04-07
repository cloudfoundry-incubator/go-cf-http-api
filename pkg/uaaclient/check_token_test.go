package uaaclient_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"../uaaclient"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAAClient - CheckToken", func() {
	It("checks that the token is valid and returns the user", func() {
		userID := "abc-123"
		mux := http.NewServeMux()
		var receivedRequest *http.Request
		mux.HandleFunc("/check_token", func(w http.ResponseWriter, r *http.Request) {
			receivedRequest = r
			w.WriteHeader(http.StatusOK)
			w.Write(checkTokenResponse(userID, "notifications.write", "admin", "test@example.com"))
		})

		ts := httptest.NewServer(mux)
		defer ts.Close()
		client, err := uaaclient.New(ts.URL, true, "client_id", "client_secret")
		Expect(err).ToNot(HaveOccurred())

		user, err := client.CheckToken("valid-token", context.Background())
		Expect(user).To(Equal(&uaaclient.User{
			ID:       userID,
			Scopes:   []string{"notifications.write"},
			Username: "admin",
			Email:    "test@example.com",
		}))
		Expect(err).ToNot(HaveOccurred())

		Expect(receivedRequest.Method).To(Equal(http.MethodPost))

		username, password, ok := receivedRequest.BasicAuth()
		Expect(username).To(Equal("client_id"))
		Expect(password).To(Equal("client_secret"))
		Expect(ok).To(BeTrue())

		Expect(receivedRequest.Header.Get("Content-Type")).To(Equal("application/x-www-form-urlencoded"))
	})

	It("returns an error if the endpoint returns a non-200 status code", func() {
		mux := http.NewServeMux()
		var receivedRequest *http.Request
		mux.HandleFunc("/check_token", func(w http.ResponseWriter, r *http.Request) {
			receivedRequest = r
			w.WriteHeader(http.StatusBadRequest)
		})

		ts := httptest.NewServer(mux)
		defer ts.Close()
		client, err := uaaclient.New(ts.URL, true, "client_id", "client_secret")
		Expect(err).ToNot(HaveOccurred())

		user, err := client.CheckToken("invalid-token", context.Background())
		Expect(user).To(BeNil())
		Expect(err).To(HaveOccurred())
	})

	It("returns an error if the user cannot be deserialized", func() {
		mux := http.NewServeMux()
		var receivedRequest *http.Request
		mux.HandleFunc("/check_token", func(w http.ResponseWriter, r *http.Request) {
			receivedRequest = r
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("}"))
		})

		ts := httptest.NewServer(mux)
		defer ts.Close()
		client, err := uaaclient.New(ts.URL, true, "client_id", "client_secret")
		Expect(err).ToNot(HaveOccurred())

		user, err := client.CheckToken("valid-token", context.Background())
		Expect(user).To(BeNil())
		Expect(err).To(HaveOccurred())
	})

})

func checkTokenResponse(userID, scope, username, email string) []byte {
	return []byte(fmt.Sprintf(`{
	  "jti": "3dfec2a647184cd9815db7880dc7c7f0",
	  "sub": "59d8c635-e5e8-485e-bfc6-dcbccbe00cc7",
	  "scope": [
		"%s"
	  ],
	  "client_id": "cf",
	  "cid": "cf",
	  "azp": "cf",
	  "grant_type": "password",
	  "user_id": "%s",
	  "origin": "uaa",
	  "user_name": "%s",
	  "email": "%s",
	  "rev_sig": "223f8692",
	  "iat": 1513375807,
	  "exp": 1513376407,
	  "iss": "https://uaa.fry.cf-app.com/oauth/token",
	  "zid": "uaa",
	  "aud": [
		"scim",
		"cloud_controller",
		"password",
		"cf",
		"uaa",
		"openid",
		"doppler",
		"routing.router_groups",
		"network"
	  ]
	}`, scope, userID, username, email))
}
