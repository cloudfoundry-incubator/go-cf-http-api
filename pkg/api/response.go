package api

import (
	"net/http"

	"../api/errors"
)

type Response struct {
	StatusCode int
	Body       interface{}
	Header     http.Header
}

func Ok(body interface{}) *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}
}

func Created() *Response {
	return &Response{
		StatusCode: http.StatusCreated,
		Body:       nil,
	}
}

func NoContent() *Response {
	return &Response{
		StatusCode: http.StatusNoContent,
		Body:       nil,
	}
}

func BadRequest(err error) *Response {
	return &Response{
		StatusCode: http.StatusBadRequest,
		Body:       wrapError(err),
	}
}

func Deleted() *Response {
	return &Response{
		StatusCode: http.StatusNoContent,
		Body:       nil,
	}
}

func ValidationError(errs []errors.ErrorResponse) *Response {
	return &Response{
		StatusCode: http.StatusUnprocessableEntity,
		Body:       errors.ErrorListResponse{Errors: errs},
	}
}

func UnprocessableEntity(err error) *Response {
	return &Response{
		StatusCode: http.StatusUnprocessableEntity,
		Body:       wrapError(err),
	}
}

func ServerError(err error) *Response {
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body:       wrapError(err),
	}
}

func Unauthorized() *Response {
	return &Response{
		StatusCode: http.StatusUnauthorized,
		Body:       nil,
	}
}

func NotFound(err error) *Response {
	return &Response{
		StatusCode: http.StatusNotFound,
		Body:       wrapError(err),
	}
}

func Forbidden(err error) *Response {
	return &Response{
		StatusCode: http.StatusForbidden,
		Body:       wrapError(err),
	}
}

func wrapError(err error) errors.ErrorListResponse {
	return errors.ErrorListResponse{
		Errors: []errors.ErrorResponse{
			{Description: err.Error()},
		},
	}
}
