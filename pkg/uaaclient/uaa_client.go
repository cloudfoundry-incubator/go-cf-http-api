package uaaclient

import (
	"errors"

	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/httpclient"
)

type User struct {
	ID       string   `json:"user_id"`
	Scopes   []string `json:"scope"`
	Email    string   `json:"email"`
	Username string   `json:"user_name"`
}

type UAAClient struct {
	client        *httpclient.HTTPClient
	client_id     string
	client_secret string
}

func New(host string, skipSSLValidation bool, client_id string, client_secret string) (*UAAClient, error) {
	client, err := httpclient.New(host, skipSSLValidation)
	if err != nil {
		return nil, err
	}

	if client_id == "" || client_secret == "" {
		return nil, errors.New("client id and secret must be set")
	}

	return &UAAClient{
		client:        client,
		client_id:     client_id,
		client_secret: client_secret,
	}, nil
}
