package handlers

type HTTPSuccessResponse struct {
	Status  string         `json:"status"`
	Data    map[string]any `json:"data,omitempty"`
	Message string         `json:"message,omitempty"`
}

type ErrorBody struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type HTTPErrorResponse struct {
	Status string    `json:"status"`
	Error  ErrorBody `json:"error"`
}

type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}
