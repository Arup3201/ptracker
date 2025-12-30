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
	keycloak "github.com/stillya/testcontainers-keycloak"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	IDPProvider       = "keycloak"
	TestKCRealm       = "ptracker"
	TestUsername      = "test_user"
	TestFirstName     = "Test"
	TestLastName      = "User"
	TestEmail         = "test@example.com"
	TestClientId      = "api"
	TestClientSecret  = "cp50avHQeX18cESEraheJvr3RhUBMq2A"
	TestPassword      = "1234"
	TestUserAgent     = "Firefox"
	TestIpAddress     = "127.0.0.1"
	TestDevice        = "HP"
	TestEncryptionKey = "ab9befcad6859b8d0e6740255bfd6e6f"
)

func createKCTestUser(serverUrl string) {
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

type attacher struct {
	mux   *http.ServeMux
	kcUrl string
}

func (atc *attacher) attachMiddleware(pattern string, handler HTTPErrorHandler) {
	authMiddleware := Authorize(atc.kcUrl, TestKCRealm)
	atc.mux.Handle(pattern, HTTPErrorHandler(authMiddleware(handler)))
}

type ProjectTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	kcContainer *keycloak.KeycloakContainer
	cookie      *http.Cookie
	mux         *http.ServeMux
	ctx         context.Context
}

func (suite *ProjectTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = db.ConnectPostgres(pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	kcContainer, err := testhelpers.CreateKeycloakContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.kcContainer = kcContainer

	adminClient, err := kcContainer.GetAdminClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	createKCTestUser(adminClient.ServerURL)

	// Keycloak implicit flow, scope=openid otherwise 403 error in /userinfo
	token, err := GetToken(adminClient.ServerURL, TestKCRealm, url.Values{
		"grant_type":    []string{"password"},
		"client_id":     []string{TestClientId},
		"client_secret": []string{TestClientSecret},
		"username":      []string{TestUsername},
		"password":      []string{TestPassword},
		"scope":         []string{"openid email profile"}, // IMPORTANT!
	})
	if err != nil {
		log.Fatal(err)
	}

	userInfo, err := GetUserInfo(adminClient.ServerURL, TestKCRealm, token.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	user, err := db.CreateUser(userInfo.Subject, IDPProvider,
		userInfo.Username, userInfo.Name, userInfo.Email, userInfo.AvatarUrl)
	if err != nil {
		log.Fatal(err)
	}

	suite.cookie, err = GetSessionCookie(token.RefreshExpiresIn, token.AccessToken, token.RefreshToken, user.Id, TestUserAgent, TestIpAddress, TestDevice, TestEncryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	suite.mux = http.NewServeMux()
	atch := &attacher{
		mux:   suite.mux,
		kcUrl: adminClient.ServerURL,
	}

	atch.attachMiddleware("POST /api/projects", CreateProject)
}

func (suite *ProjectTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}

	if err := suite.kcContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}
}

func (suite *ProjectTestSuite) TestCreateProject() {
	t := suite.T()

	t.Run("success - create a project - 200 response", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		assert.Equal(t, res.Result().StatusCode, http.StatusOK)
	})
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
