package controller_fixtures

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/ptracker/interfaces"
)

type ControllerFixtures struct {
	ctx     context.Context
	store   interfaces.Store
	Handler http.Handler
}

func NewControllerFixtures(ctx context.Context, store interfaces.Store) *ControllerFixtures {
	fixtures := &ControllerFixtures{
		ctx:   ctx,
		store: store,
	}

	return fixtures
}

// RequestAs creates a request with userID in context
func (cf *ControllerFixtures) RequestAs(userID string) *AuthenticatedClient {
	return &AuthenticatedClient{
		handler: cf.Handler,
		userID:  userID,
	}
}

type AuthenticatedClient struct {
	handler http.Handler
	userID  string
}

func (ac *AuthenticatedClient) Get(path string) *httptest.ResponseRecorder {
	return ac.do("GET", path, nil)
}

func (ac *AuthenticatedClient) Post(path string, body interface{}) *httptest.ResponseRecorder {
	return ac.do("POST", path, body)
}

func (ac *AuthenticatedClient) Put(path string, body interface{}) *httptest.ResponseRecorder {
	return ac.do("PUT", path, body)
}

func (ac *AuthenticatedClient) Delete(path string) *httptest.ResponseRecorder {
	return ac.do("DELETE", path, nil)
}

func (ac *AuthenticatedClient) do(method, path string, body any) *httptest.ResponseRecorder {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Inject userID into request context
	ctx := context.WithValue(req.Context(), "user_id", ac.userID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	ac.handler.ServeHTTP(rec, req)

	return rec
}
