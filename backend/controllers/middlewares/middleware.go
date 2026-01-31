package middlewares

import (
	"net/http"

	"github.com/ptracker/controllers"
)

type Middleware interface {
	Next(next http.Handler) controllers.HTTPErrorHandler
}
