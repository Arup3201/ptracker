package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ptracker/db"
	"github.com/ptracker/testhelpers"
	"github.com/stretchr/testify/assert"
)

const (
	IDPProvider      = "keycloak"
	TestKCRealm      = "ptracker"
	TestUserId       = "b24b8712-6819-47fb-83e5-11eb28280a2f"
	TestUsername     = "test_user"
	TestFirstName    = "Test"
	TestLastName     = "User"
	TestEmail        = "test@example.com"
	TestClientId     = "api"
	TestClientSecret = "cp50avHQeX18cESEraheJvr3RhUBMq2A"
	TestPassword     = "1234"
	TestUserAgent    = "Firefox"
	TestIpAddress    = "127.0.0.1"
	TestDevice       = "HP"
)

func createKCTestUser(t testing.TB, serverUrl string) {
	t.Helper()

	credentials := url.Values{
		"grant_type": []string{"password"},
		"client_id":  []string{"admin-cli"},
		"username":   []string{"admin"},
		"password":   []string{"admin"},
	}
	tokenUrl := serverUrl + "/realms/master/protocol/openid-connect/token"
	res, err := http.PostForm(tokenUrl, credentials)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		var kcError KCError
		json.NewDecoder(res.Body).Decode(&kcError)
		log.Fatalf("keycloak get token error: %v\n", kcError)
	}

	var accessToken struct {
		Value string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&accessToken); err != nil {
		log.Fatal(err)
	}

	type Credential struct {
		Value     string `json:"value"`
		Type      string `json:"type"`
		Temporary bool   `json:"temporary"`
	}
	var user = struct {
		Username      string       `json:"username"`
		Firstname     string       `json:"firstName"`
		Lastname      string       `json:"lastName"`
		Email         string       `json:"email"`
		Enabled       bool         `json:"enabled"`
		EmailVerified bool         `json:"emailVerified"`
		Credentials   []Credential `json:"credentials"`
	}{
		Username:      TestUsername,
		Firstname:     TestFirstName,
		Lastname:      TestLastName,
		Email:         TestEmail,
		Enabled:       true,
		EmailVerified: true,
		Credentials: []Credential{
			{
				Type:      "password",
				Value:     TestPassword,
				Temporary: false,
			},
		},
	}
	payload, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	userUrl := fmt.Sprintf(serverUrl+"/admin/realms/%s/users", TestKCRealm)
	req, err := http.NewRequest(
		"POST",
		userUrl,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken.Value))
	req.Header.Set("Content-Type", "application/json")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusCreated {
		var kcError KCError
		json.NewDecoder(res.Body).Decode(&kcError)
		log.Fatalf("keycloak create user error: %v\n", kcError)
	}
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

	ctx = context.Background()
	kcContainer, err := testhelpers.CreateKeycloakContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	adminClient, err := kcContainer.GetAdminClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Keycloak register user
	createKCTestUser(t, adminClient.ServerURL)

	// Keycloak implicit flow
	token, err := GetToken(adminClient.ServerURL, TestKCRealm, url.Values{
		"grant_type":    []string{"password"},
		"client_id":     []string{TestClientId},
		"client_secret": []string{TestClientSecret},
		"username":      []string{TestUsername},
		"password":      []string{TestPassword},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get user information
	userInfo, err := GetUserInfo(adminClient.ServerURL, TestKCRealm, token.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	// Create user
	user, err := db.CreateUser(userInfo.Subject, IDPProvider,
		userInfo.Username, userInfo.Name, userInfo.Email, userInfo.AvatarUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create session
	cookie, err := GetSessionCookie(token.RefreshExpiresIn, token.AccessToken, token.RefreshToken, user.Id, TestUserAgent, TestIpAddress, TestDevice)
	if err != nil {
		log.Fatal(err)
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
