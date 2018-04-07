package errors

type ErrorListResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Description string `json:"description"`
}

func (e *ErrorListResponse) GetErrors() []string {
	list := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		list = append(list, err.Description)
	}

	return list
}
