package middlewares

import (
	"net/http"

	"github.com/ptracker/internal/controllers"
)

type Middleware interface {
	Handler(next http.Handler) controllers.HTTPErrorHandler
}
