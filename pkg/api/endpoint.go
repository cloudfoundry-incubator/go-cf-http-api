package api

import "./auth"

type Endpoint struct {
	Path   string
	Method string
	Auth   *auth.Config
	Handle func(r Request) *Response
}
