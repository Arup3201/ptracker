package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ptracker/db"
	"github.com/ptracker/utils"
	"github.com/redis/go-redis/v9"
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

// HTTP Error ID
const (
	ERR_UNAUTHORIZED       = "unauthorized"
	ERR_INVALID_QUERY      = "invalid_query"
	ERR_INVALID_BODY       = "invalid_body"
	ERR_ACCESS_DENIED      = "access_denied"
	ERR_RESOURCE_NOT_FOUND = "resource_not_found"
	ERR_SERVER_ERROR       = "server_error"
)

type HTTPErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (fn HTTPErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		fmt.Printf("[ERROR] %s\n", err)

		if httpError, ok := err.(*HTTPError); ok {
			w.WriteHeader(httpError.Code)
			json.NewEncoder(w).Encode(HTTPErrorResponse{
				Status: "error",
				Error: ErrorBody{
					Id:      httpError.ErrId,
					Message: httpError.Message,
				},
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(HTTPErrorResponse{
				Status: "error",
				Error: ErrorBody{
					Id:      ERR_SERVER_ERROR,
					Message: "Something unexpected happened, please try again later.",
				},
			})
		}
		return
	}
}

type AuthMiddleware func(http.Handler) HTTPErrorHandler

func Authorize(redis *redis.Client,
	keycloakUrl, keycloakRealm string) AuthMiddleware {
	return func(next http.Handler) HTTPErrorHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			public := []string{"/auth/login", "/auth/callback", "/auth/refresh"}
			for _, endpoint := range public {
				if strings.Contains(strings.TrimPrefix(r.URL.Path, "/api"), endpoint) {
					next.ServeHTTP(w, r)
					return nil
				}
			}

			sessionId, err := utils.GetSessionIdFromCookie(r.Cookies(), SESSION_COOKIE_NAME)
			if err != nil {
				return &HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "User is not authorized",
					Err:     fmt.Errorf("authorize: session cookie not found"),
				}
			}

			tokenKey := utils.GetAccessTokenKey(sessionId)
			accessToken, err := redis.Get(context.Background(), tokenKey).Result()
			if err != nil {
				return &HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "User is not authorized",
					Err:     fmt.Errorf("authorize: access token expired or missing"),
				}
			}

			sub, err := verifyAccessToken(keycloakUrl, keycloakRealm, accessToken)
			if err != nil {
				return &HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "User is not authorized",
					Err:     fmt.Errorf("authorize: %w", err),
				}
			}

			user, err := db.GetUserBySub(sub)

			ctx := context.WithValue(r.Context(), "user_id", user.Id)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
			return nil
		}
	}
}

func verifyAccessToken(kcUrl, kcRealm, accessToken string) (string, error) {
	jwksKeySet, err := jwk.Fetch(context.TODO(), fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", kcUrl, kcRealm))
	if err != nil {
		return "", fmt.Errorf("verify access token: %w", err)
	}

	token, err := jwt.Parse([]byte(accessToken), jwt.WithKeySet(jwksKeySet), jwt.WithValidate(true))
	if err != nil {
		return "", fmt.Errorf("verify access token: %w", err)
	}

	if token == nil {
		return "", fmt.Errorf("parsed token is null")
	}

	return token.Subject(), nil
}

type tokenBucket struct {
	ctx         context.Context
	capacity    int // maximum tokens in the bucket - handles traffic bursts
	rate        int // refill rate per second
	redisClient *redis.Client
	tbFunc      *redis.Script
	retry       int64
}

const RedisLuaScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])

local now = redis.call("TIME")
local now_sec = tonumber(now[1])
local now_usec = tonumber(now[2])
local now_ts = now_sec + now_usec / 1e6

local tokens = tonumber(redis.call("HGET", key, "tokens"))
if tokens == nil then
	tokens = capacity -- start full
end

local last_refill = tonumber(redis.call("HGET", key, "last_refill"))
if last_refill == nil then
	last_refill = now_ts
end

-- refill
local elapsed = now_ts - last_refill
local tokens_to_add = math.floor(elapsed) * rate

if tokens_to_add > 0 then
	tokens = math.min(capacity, tokens+tokens_to_add)
	redis.call("HSET", key, "last_refill", now_ts)
end

-- deny if not enough tokens
if tokens < 1 then
    local retry_after = math.ceil(((1 - (elapsed*rate)) / rate)*1000)
    return {0, retry_after}
end

-- consume
tokens = tokens - 1
redis.call("HSET", key, "tokens", tokens)

return {1, 0}
`

func CreateTokenBucket(rdc *redis.Client, cap, rate int) *tokenBucket {
	function := redis.NewScript(RedisLuaScript)

	return &tokenBucket{
		ctx:         context.Background(),
		capacity:    cap,
		rate:        rate,
		redisClient: rdc,
		tbFunc:      function,
	}
}

func (tb *tokenBucket) GetToken(key string) (string, error) {
	value, err := tb.redisClient.HGet(tb.ctx, key, "tokens").Result()
	if err != nil {
		return "", fmt.Errorf("get token: %w", err)
	}

	return value, nil
}

func (tb *tokenBucket) AllowRequest(key string) (bool, error) {
	value, err := tb.tbFunc.Run(tb.ctx, tb.redisClient, []string{key}, tb.capacity, tb.rate).Slice()
	if err != nil {
		return false, fmt.Errorf("allow request run lua script: %w", err)
	}

	if value[0].(int64) == 0 {
		tb.retry = value[1].(int64)
		return false, nil
	}

	return true, nil
}

type RateLimiter func(HTTPErrorHandler) HTTPErrorHandler

func TokenBucketRateLimiter(rdc *redis.Client, capacity, rate int) RateLimiter {
	bucket := CreateTokenBucket(rdc, capacity, rate)
	return func(next HTTPErrorHandler) HTTPErrorHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			userId, err := utils.GetUserId(r)
			if err != nil {
				return fmt.Errorf("rate limiter get user id: %w", err)
			}

			redisKey := utils.GetBucketKey(userId)
			allow, err := bucket.AllowRequest(redisKey)
			if err != nil {
				return fmt.Errorf("rate limiter allow request: %w", err)
			}

			if allow {
				rec := httptest.NewRecorder()

				next.ServeHTTP(rec, r)

				tokens, err := bucket.GetToken(redisKey)
				if err != nil {
					// redis error, continue...
					log.Printf("[ERROR] rate limiter: %s\n", err)
				}

				maps.Copy(w.Header(), rec.Header())
				w.Header().Add("X-Ratelimit-Remaining", tokens)
				w.Header().Add("X-Ratelimit-Limit", strconv.Itoa(bucket.capacity))
				w.WriteHeader(rec.Result().StatusCode)
				w.Write(rec.Body.Bytes())
			} else {
				w.Header().Add("X-Ratelimit-Retry-After", strconv.Itoa(int(bucket.retry))+" ms")
				w.WriteHeader(http.StatusTooManyRequests)
			}
			return nil
		}
	}
}
