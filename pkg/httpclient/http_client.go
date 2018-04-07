package httpclient

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	client *http.Client
	host   *url.URL
}

func New(host string, skipSSLValidation bool) (*HTTPClient, error) {
	if host == "" {
		return nil, errors.New("host cannot be empty")
	}

	uri, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: skipSSLValidation,
				},
			},
			Timeout: 10 * time.Second,
		},
		host: uri,
	}, nil
}

func (c *HTTPClient) GetRequest(endpoint string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, c.url(endpoint), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *HTTPClient) PostRequest(endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, c.url(endpoint), body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *HTTPClient) PutRequest(endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPut, c.url(endpoint), body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *HTTPClient) DeleteRequest(endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodDelete, c.url(endpoint), body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *HTTPClient) Do(req *http.Request, ctx context.Context) (*http.Response, error) {
	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) url(endpoint string) string {
	u := c.newUrl()

	return u.String() + endpoint
}

func (c *HTTPClient) newUrl() *url.URL {
	newUrl := *c.host
	return &newUrl
}
