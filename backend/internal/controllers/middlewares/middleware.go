package middlewares

import (
	"net/http"

	"github.com/ptracker/internal/controllers"
)

type Middleware interface {
	Next(next http.Handler) controllers.HTTPErrorHandler
}
