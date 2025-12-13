package apierr

type ResourceNotFound struct{}

func (e *ResourceNotFound) Error() string {
	return "Resource not found"
}
