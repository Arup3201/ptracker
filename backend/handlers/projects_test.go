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
	"github.com/ptracker/models"
	"github.com/ptracker/testhelpers"
	"github.com/redis/go-redis/v9"
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
	TestUserAgent     = "Mozilla/5.0 (X11; Linux x86_64; rv:144.0) Gecko/20100101 Firefox/144.0"
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
	redis *redis.Client
	kcUrl string
}

func (atc *attacher) attachMiddleware(pattern string, handler HTTPErrorHandler) {
	authMiddleware := Authorize(atc.redis, atc.kcUrl, TestKCRealm)
	atc.mux.Handle(pattern, HTTPErrorHandler(authMiddleware(handler)))
}

type ProjectTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	kcContainer *keycloak.KeycloakContainer
	redis       *redis.Client
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

	redisContainer, err := testhelpers.CreateRedisContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	connString, err := redisContainer.ConnectionString(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	opt, err := redis.ParseURL(connString)
	if err != nil {
		log.Fatal(err)
	}

	suite.redis = redis.NewClient(opt)

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

	suite.cookie, err = GetSessionCookie(suite.redis, token.RefreshExpiresIn, token.AccessToken, token.RefreshToken, user.Id, TestUserAgent, TestIpAddress, TestDevice, TestEncryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	suite.mux = http.NewServeMux()
	atch := &attacher{
		mux:   suite.mux,
		redis: suite.redis,
		kcUrl: adminClient.ServerURL,
	}

	atch.attachMiddleware("POST /api/projects", CreateProject)
	atch.attachMiddleware("GET /api/projects/{id}", GetProject)
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

	t.Run("success response is 200", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak"
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

	t.Run("success response body is correct", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		var responseBody HTTPSuccessResponse[models.CreatedProject]
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, RESPONSE_SUCCESS_STATUS, responseBody.Status)
		assert.Equal(t, "PTracker Go", responseBody.Data.Name)
		assert.Equal(t, "Collaborative project tracking application with Go", *responseBody.Data.Description)
		assert.Equal(t, "Go, React, TypeScript, PostgreSQL, Keycloak", *responseBody.Data.Skills)
	})

	t.Run("error with missing name in payload", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		var responseBody HTTPErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, RESPONSE_ERROR_STATUS, responseBody.Status)
		assert.Equal(t, ERR_INVALID_BODY, responseBody.Error.Id)
		assert.Equal(t, "Project 'name' is missing", responseBody.Error.Message)
	})

	t.Run("error with unknown fields in payload", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak",
			"custom": "It is a custom field"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		var responseBody HTTPErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, RESPONSE_ERROR_STATUS, responseBody.Status)
		assert.Equal(t, ERR_INVALID_BODY, responseBody.Error.Id)
		assert.Equal(t, "Project must have 'name' with optional 'description' and 'skills' fields only", responseBody.Error.Message)
	})
}

func (suite *ProjectTestSuite) TestGetProjectDetails() {
	t := suite.T()

	t.Run("get project details success", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()
		suite.mux.ServeHTTP(res, req)
		if res.Result().StatusCode != http.StatusOK {
			t.Log("project create failed")
			t.Fail()
		}
		var createdProject HTTPSuccessResponse[models.CreatedProject]
		if err := json.NewDecoder(res.Body).Decode(&createdProject); err != nil {
			t.Log("project create decode failed")
			t.Fail()
		}
		projectId := createdProject.Data.Id
		req, err = http.NewRequest("GET", "/api/projects/"+projectId, nil)
		if err != nil {
			t.Log("project get request create failed")
			t.Fail()
		}
		req.AddCookie(suite.cookie)
		res = httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		var projectDetails HTTPSuccessResponse[models.ProjectDetails]
		if err := json.NewDecoder(res.Body).Decode(&projectDetails); err != nil {
			t.Log("project get decode failed")
			t.Fail()
		}
		assert.Equal(t, RESPONSE_SUCCESS_STATUS, projectDetails.Status)
		assert.Equal(t, "PTracker Go", projectDetails.Data.Name)
		assert.Equal(t, "Collaborative project tracking application with Go", *projectDetails.Data.Description)
		assert.Equal(t, "Go, React, TypeScript, PostgreSQL, Keycloak", *projectDetails.Data.Skills)
		assert.Equal(t, models.ROLE_OWNER, projectDetails.Data.Role)
		assert.Equal(t, 0, projectDetails.Data.UnassignedTasks)
		assert.Equal(t, 0, projectDetails.Data.OngoingTasks)
		assert.Equal(t, 0, projectDetails.Data.CompletedTasks)
		assert.Equal(t, 0, projectDetails.Data.AbandonedTasks)
		assert.Equal(t, 0, projectDetails.Data.MemberCount)
	})
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
