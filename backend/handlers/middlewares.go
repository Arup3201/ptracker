package handlers

import (
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		log.Printf("%s %s - %d", r.Method, r.RequestURI, rec.Result().StatusCode)

		maps.Copy(w.Header(), rec.Header())
		w.WriteHeader(rec.Result().StatusCode)
		w.Write(rec.Body.Bytes())
	})
}
