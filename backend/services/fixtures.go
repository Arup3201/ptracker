package services

import (
	"database/sql"
	"log"

	"github.com/ptracker/models"
)

var USER_ONE = map[string]string{
	"idp_provider": "keycloak",
	"idp_subject":  "f6e1d9a0-7b3c-4d5e-8f2a-1c9b8e7d6f5a",
	"kc_realm":     "ptracker",
	"username":     "test_user",
	"display_name": "Test User",
	"email":        "test@example.com",
}
var USER_TWO = map[string]string{
	"idp_provider": "keycloak",
	"idp_subject":  "f6e1d9a0-7b3c-4d5e-8f2a-1c9b8e7d8d0a",
	"kc_realm":     "ptracker",
	"username":     "test_user_1",
	"display_name": "Test User 1",
	"email":        "test1@example.com",
}

var PROJECT_ONE = map[string]string{
	"name":        "Project A",
	"description": "Description for Project A",
	"skills":      "C, C++, Python",
}
var PROJECT_TWO = map[string]string{
	"name":        "Project B",
	"description": "Description for Project B",
	"skills":      "Java",
}
var PROJECT_THREE = map[string]string{
	"name":        "Project C",
	"description": "Description for Project C",
	"skills":      "Kotlin, Android",
}

var USER_FIXTURES = []models.User{}
var PROJECT_FIXTURES = []models.Project{}

func CreatFixtures(conn *sql.DB) {
	user_fixture_data := []map[string]string{USER_ONE, USER_TWO}
	userStore := &models.UserStore{
		DB: conn,
	}
	for _, fixture := range user_fixture_data {
		userId, err := userStore.Create(fixture["idp_provider"], fixture["idp_subject"],
			fixture["kc_realm"], fixture["username"], fixture["display_name"],
			fixture["email"])
		if err != nil {
			log.Fatal(err)
		}

		user, err := userStore.Get(userId)

		USER_FIXTURES = append(USER_FIXTURES, user)
	}

	projectStore := &models.ProjectStore{
		DB:     conn,
		UserId: USER_FIXTURES[0].Id,
	}
	userOneProjects := []map[string]string{PROJECT_ONE, PROJECT_TWO}
	for _, fixture := range userOneProjects {
		id, err := projectStore.Create(fixture["name"],
			fixture["description"],
			fixture["skills"])
		if err != nil {
			log.Fatal(err)
		}

		project, err := projectStore.Get(id)
		if err != nil {
			log.Fatal(err)
		}

		PROJECT_FIXTURES = append(PROJECT_FIXTURES, *project)
	}

	projectStore = &models.ProjectStore{
		DB:     conn,
		UserId: USER_FIXTURES[1].Id,
	}
	userTwoProjects := []map[string]string{PROJECT_THREE}
	for _, fixture := range userTwoProjects {
		id, err := projectStore.Create(fixture["name"],
			fixture["description"],
			fixture["skills"])
		if err != nil {
			log.Fatal(err)
		}

		project, err := projectStore.Get(id)
		if err != nil {
			log.Fatal(err)
		}

		PROJECT_FIXTURES = append(PROJECT_FIXTURES, *project)
	}
}
