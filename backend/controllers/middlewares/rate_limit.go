package middlewares

import (
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/ptracker/controllers"
	"github.com/ptracker/interfaces"
	"github.com/ptracker/utils"
)

type rateLimitMiddleware struct {
	service interfaces.LimiterService
}

func NewRateLimitMiddleware(service interfaces.LimiterService) Middleware {
	return &rateLimitMiddleware{
		service: service,
	}
}

func (m *rateLimitMiddleware) Next(next http.Handler) controllers.HTTPErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		userId, err := utils.GetUserId(r)
		if err != nil {
			return fmt.Errorf("get user id: %w", err)
		}

		allow, err := m.service.IsAllowed(r.Context(), userId)
		if err != nil {
			return fmt.Errorf("limiter service is allowed: %w", err)
		}

		if allow {
			rec := httptest.NewRecorder()

			next.ServeHTTP(rec, r)

			tokens, err := m.service.GetTokens(r.Context(), userId)
			if err != nil {
				log.Printf("[ERROR] limiter service get tokens: %s\n", err)
				maps.Copy(w.Header(), rec.Header())
				w.WriteHeader(rec.Result().StatusCode)
				w.Write(rec.Body.Bytes())
				return nil
			}

			capacity := m.service.GetCapacity(r.Context())

			maps.Copy(w.Header(), rec.Header())
			w.Header().Add("X-Ratelimit-Remaining", strconv.Itoa(tokens))
			w.Header().Add("X-Ratelimit-Limit", strconv.Itoa(capacity))
			w.WriteHeader(rec.Result().StatusCode)
			w.Write(rec.Body.Bytes())
		} else {
			retry, err := m.service.GetRetryTime(r.Context(), userId)
			if err != nil {
				log.Printf("[ERROR] limiter service get retry time: %s\n", err)
				return nil
			}
			w.Header().Add("X-Ratelimit-Retry-After", strconv.Itoa(retry))
			w.WriteHeader(http.StatusTooManyRequests)
		}
		return nil
	}
}
