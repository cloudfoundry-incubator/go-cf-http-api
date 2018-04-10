package api

import "github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api/auth"

type Endpoint struct {
	Path   string
	Method string
	Auth   *auth.Config
	Handle func(r Request) *Response
}
