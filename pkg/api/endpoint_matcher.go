package api

import (
	"fmt"

	"github.com/onsi/gomega/types"
)

type endpointMatcher struct{}

func BeCompleteEndpoint() types.GomegaMatcher {
	return &endpointMatcher{}
}

func (m *endpointMatcher) Match(actual interface{}) (bool, error) {
	endpoint, err := getEndpoint(actual)
	if err != nil {
		return false, err
	}

	if endpoint.Path == "" {
		return false, fmt.Errorf("Path cannot be empty")
	}

	if endpoint.Method == "" {
		return false, fmt.Errorf("Method cannot be empty")
	}

	if endpoint.Handle == nil {
		return false, fmt.Errorf("Handle cannot be nil")
	}

	if endpoint.Auth == nil {
		return false, fmt.Errorf("Auth cannot be nil")
	}

	return true, nil
}

func (m *endpointMatcher) FailureMessage(actual interface{}) string {
	return "Expected endpoint to contain all required fields"
}

func (m *endpointMatcher) NegatedFailureMessage(actual interface{}) string {
	return "Expected endpoint to not contain all the required fields"
}

func getEndpoint(actual interface{}) (Endpoint, error) {
	var (
		endpoint Endpoint
		ok       bool
	)

	endpoint, ok = actual.(Endpoint)
	if !ok {
		e, ok := actual.(*Endpoint)
		if !ok {
			return Endpoint{}, fmt.Errorf("EndpointMatcher expects a api.Endpoint")
		}
		endpoint = *e
	}

	return endpoint, nil
}
