package apierr

import "errors"

var ErrResourceNotFound = errors.New("resource not found")
var ErrDuplicate = errors.New("duplicate entry in the database")
var ErrInvalidValue = errors.New("invalid value")
var ErrForbidden = errors.New("forbidden action")
