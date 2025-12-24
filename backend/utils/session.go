package utils

import (
	"fmt"
	"net/http"
	"slices"
)

func GetSessionIdFromCookie(cookies []*http.Cookie, sessionCookieName string) (string, error) {
	ind := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
		return cookie.Name == sessionCookieName
	})
	if ind == -1 {
		return "", fmt.Errorf("session cookie missing")
	}

	sessionId := cookies[ind].Value

	return sessionId, nil
}

func GetUserId(req *http.Request) (string, error) {
	ctx := req.Context()
	userId, ok := ctx.Value("user_id").(string)
	if !ok || userId == "" {
		return "", fmt.Errorf("get user id: empty context")
	}
	return userId, nil
}
