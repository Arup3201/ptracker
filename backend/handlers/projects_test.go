package handlers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ptracker/db"
	"github.com/ptracker/models"
	"github.com/ptracker/testhelpers"
	"github.com/stretchr/testify/assert"
)

const (
	IDPProvider     = "keycloak"
	TestUsername    = "test_user"
	TestDisplayName = "User Test"
	TestEmail       = "test@example.com"
	TestAvatarUrl   = "https://example.com/avatar/test_user.png"
)

func CreateDummySession(t testing.TB, userId, refreshToken string, expires int64) *models.Session {
	t.Helper()

	session, err := db.CreateSession(userId, []byte(refreshToken), "Firefox", "127.0.0.1", "Windows", time.Now().Add(time.Duration(expires)*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	return session
}

func attachMiddlewares(mux *http.ServeMux, pattern string, handler HTTPErrorHandler) {
	mux.Handle(pattern, HTTPErrorHandler(Authorize(handler)))
}

func TestProjectCreate(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = db.ConnectPostgres(pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	user, err := db.CreateUser(IDPSubject, IDPProvider, TestUsername, TestDisplayName, TestEmail, TestAvatarUrl)
	if err != nil {
		log.Fatal(err)
	}

	session := CreateDummySession(t, user.Id)
	cookie := &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    session.Id,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  session.ExpiresAt,
	}

	mux := http.NewServeMux()
	attachMiddlewares(mux, "POST /api/projects", CreateProject)

	t.Run("success - create a project - 200 response", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go"
		`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(cookie)
		res := httptest.NewRecorder()

		mux.ServeHTTP(res, req)

		assert.Equal(t, res.Result().StatusCode, http.StatusOK)
	})
}
