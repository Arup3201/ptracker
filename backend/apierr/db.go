package apierr

import "errors"

var ErrResourceNotFound = errors.New("resource not found")
var ErrDuplicate = errors.New("duplicate entry in the database")
