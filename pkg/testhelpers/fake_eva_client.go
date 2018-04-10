package testhelpers

import (
	"fmt"
	"reflect"

	"strings"

	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api/errors"
	"github.com/cloudfoundry-incubator/go-cf-http-api/pkg/evaclient"
)

type fakeClient struct {
	requests  map[string]interface{}
	responses map[string]evaclient.ApiResponse
	results   map[string]interface{}
}

func NewFakeEVAClient() *fakeClient {
	return &fakeClient{
		requests:  make(map[string]interface{}, 0),
		responses: make(map[string]evaclient.ApiResponse, 0),
		results:   make(map[string]interface{}, 0),
	}
}

func (f *fakeClient) ExpectsGet(path string, responseBody interface{}, response evaclient.ApiResponse) func() error {
	key := "GET " + path
	f.responses[key] = response
	f.results[key] = responseBody

	return func() error {
		_, ok := f.requests[key]
		if !ok {
			return fmt.Errorf("get request to %s not made", path)
		}

		return nil
	}
}

func (f *fakeClient) Get(path string, result interface{}) evaclient.ApiResponse {
	key := "GET " + path
	f.requests[key] = true

	if f.results[key] != nil {
		r := reflect.ValueOf(result).Elem()
		r.Set(reflect.ValueOf(f.results[key]))
	}

	return f.responses[key]
}

func (f *fakeClient) ExpectsPost(path string, requestBody interface{}, response evaclient.ApiResponse) func() error {
	key := "POST " + path
	f.responses[key] = response

	return func() error {
		actual, ok := f.requests[key]
		if !ok {
			return fmt.Errorf("post request to %s not made", path)
		}

		if !reflect.DeepEqual(actual, requestBody) {
			return fmt.Errorf("expected %v to match %v", actual, requestBody)
		}

		return nil
	}
}

func (f *fakeClient) Post(path string, requestBody interface{}) evaclient.ApiResponse {
	key := "POST " + path
	f.requests[key] = requestBody

	return f.responses[key]
}

func (f *fakeClient) ExpectsPut(path string, requestBody interface{}, response evaclient.ApiResponse) func() error {
	key := "PUT " + path
	f.responses[key] = response

	return func() error {
		actual, ok := f.requests[key]
		if !ok {
			return fmt.Errorf("put request to %s not made", path)
		}

		if !reflect.DeepEqual(actual, requestBody) {
			return fmt.Errorf("expected %v to match %v", actual, requestBody)
		}

		return nil
	}
}

func (f *fakeClient) Put(path string, requestBody interface{}) evaclient.ApiResponse {
	key := "PUT " + path
	f.requests[key] = requestBody

	return f.responses[key]
}

func (f *fakeClient) ExpectsDelete(path string, response evaclient.ApiResponse) func() error {
	key := "DELETE " + path
	f.responses[key] = response

	return func() error {
		_, ok := f.requests[key]
		if !ok {
			allPaths := make([]string, 0)
			for k := range f.requests {
				allPaths = append(allPaths, k)
			}
			return fmt.Errorf("delete request to %s not made. requests made to %s", path, strings.Join(allPaths, ", "))
		}

		return nil
	}
}

func (f *fakeClient) Delete(path string) evaclient.ApiResponse {
	key := "DELETE " + path
	f.requests[key] = true

	return f.responses[key]
}

func (f *fakeClient) OkResponse() evaclient.ApiResponse {
	return evaclient.ApiResponse{
		Ok:     true,
		Errors: nil,
	}
}

func (f *fakeClient) ErrorResponse(errsToReturn ...string) evaclient.ApiResponse {
	errs := []errors.ErrorResponse{}
	for _, err := range errsToReturn {
		errs = append(errs, errors.ErrorResponse{Description: err})
	}

	return evaclient.ApiResponse{
		Ok:     false,
		Errors: errs,
	}
}
