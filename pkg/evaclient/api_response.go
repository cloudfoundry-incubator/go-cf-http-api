package evaclient

import "github.com/cloudfoundry-incubator/go-cf-http-api/pkg/api/errors"

type ApiResponse struct {
	Ok     bool
	Errors []errors.ErrorResponse
}
