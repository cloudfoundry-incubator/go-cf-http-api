package api

import (
	"net/http"

	"../api/auth"
)

type InfoResponse struct {
	URL     string `json:"url"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func InfoEndpoint(response *InfoResponse) *Endpoint {
	return &Endpoint{
		Path:   "/v1/info",
		Method: http.MethodGet,
		Auth:   auth.None,
		Handle: func(r Request) *Response {
			return Ok(response)
		},
	}
}
