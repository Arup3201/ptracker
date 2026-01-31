package middlewares

import (
	"log"
	"maps"
	"net/http"
	"net/http/httptest"

	"github.com/ptracker/controllers"
)

type loggingMiddleware struct{}

func (m *loggingMiddleware) Next(next http.Handler) controllers.HTTPErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		log.Printf("%s %s - %d", r.Method, r.RequestURI, rec.Result().StatusCode)

		maps.Copy(w.Header(), rec.Header())
		w.WriteHeader(rec.Result().StatusCode)
		w.Write(rec.Body.Bytes())

		return nil
	}
}
