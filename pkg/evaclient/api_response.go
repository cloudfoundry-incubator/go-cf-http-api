package evaclient

import "../api/errors"

type ApiResponse struct {
	Ok     bool
	Errors []errors.ErrorResponse
}
