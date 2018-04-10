package evaclient

import (
	"context"
	"encoding/json"

	"fmt"

	"bytes"

	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api/errors"
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/httpclient"
)

type EVAClient struct {
	client    *httpclient.HTTPClient
	userToken string
}

func New(host string, skipSSLValidation bool, userToken string) (*EVAClient, error) {
	client, err := httpclient.New(host, skipSSLValidation)
	if err != nil {
		return nil, err
	}

	return &EVAClient{
		client:    client,
		userToken: userToken,
	}, nil
}

func (e *EVAClient) Get(path string, result interface{}) ApiResponse {
	req, err := e.client.GetRequest(path)
	if err != nil {
		return wrapError(err)
	}

	req.Header.Add("Authorization", e.userToken)

	resp, err := e.client.Do(req, context.Background())
	if err != nil {
		return wrapError(err)
	}
	defer resp.Body.Close()

	ok := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !ok {
		var apiErrors errors.ErrorListResponse
		err := json.NewDecoder(resp.Body).Decode(&apiErrors)
		if err != nil {
			return wrapError(fmt.Errorf("No errors specified. Received status %d", resp.StatusCode))
		}

		return ApiResponse{
			Ok:     false,
			Errors: apiErrors.Errors,
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return wrapError(err)
	}

	return ApiResponse{Ok: ok}
}

func (e *EVAClient) Put(path string, body interface{}) ApiResponse {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return wrapError(err)
	}

	req, err := e.client.PutRequest(path, bytes.NewReader(bodyBytes))
	if err != nil {
		return wrapError(err)
	}

	req.Header.Add("Authorization", e.userToken)

	resp, err := e.client.Do(req, context.Background())
	if err != nil {
		return wrapError(err)
	}
	defer resp.Body.Close()

	ok := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !ok {
		var apiErrors errors.ErrorListResponse
		err := json.NewDecoder(resp.Body).Decode(&apiErrors)
		if err != nil {
			return wrapError(fmt.Errorf("No errors specified. Received status %d", resp.StatusCode))
		}

		return ApiResponse{
			Ok:     false,
			Errors: apiErrors.Errors,
		}
	}

	return ApiResponse{Ok: ok}
}

func (e *EVAClient) Post(path string, body interface{}) ApiResponse {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return wrapError(err)
	}

	req, err := e.client.PostRequest(path, bytes.NewReader(bodyBytes))
	if err != nil {
		return wrapError(err)
	}

	req.Header.Add("Authorization", e.userToken)

	resp, err := e.client.Do(req, context.Background())
	if err != nil {
		return wrapError(err)
	}
	defer resp.Body.Close()

	ok := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !ok {
		var apiErrors errors.ErrorListResponse
		err := json.NewDecoder(resp.Body).Decode(&apiErrors)
		if err != nil {
			return wrapError(fmt.Errorf("No errors specified. Received status %d", resp.StatusCode))
		}

		return ApiResponse{
			Ok:     false,
			Errors: apiErrors.Errors,
		}
	}

	return ApiResponse{Ok: ok}
}

func (e *EVAClient) Delete(path string) ApiResponse {
	req, err := e.client.DeleteRequest(path, bytes.NewReader(nil))
	if err != nil {
		return wrapError(err)
	}

	req.Header.Add("Authorization", e.userToken)

	resp, err := e.client.Do(req, context.Background())
	if err != nil {
		return wrapError(err)
	}
	defer resp.Body.Close()

	ok := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !ok {
		var apiErrors errors.ErrorListResponse
		err := json.NewDecoder(resp.Body).Decode(&apiErrors)
		if err != nil {
			return wrapError(fmt.Errorf("No errors specified. Received status %d", resp.StatusCode))
		}

		return ApiResponse{
			Ok:     false,
			Errors: apiErrors.Errors,
		}
	}

	return ApiResponse{Ok: ok}
}

func wrapError(err error) ApiResponse {
	return ApiResponse{
		Ok: false,
		Errors: []errors.ErrorResponse{
			{Description: err.Error()},
		},
	}
}
