package repositories

import (
	"context"
	"log"
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

var USER_FIXTURES = []string{}
var PROJECT_FIXTURES = []string{}

func CreatFixtures(db Execer) {
	userRepo := NewUserRepo(db)
	projectRepo := NewProjectRepo(db)

	user_fixture_data := []map[string]string{USER_ONE, USER_TWO}

	ctx := context.Background()
	var displayName, avatarUrl string
	for _, fixture := range user_fixture_data {
		displayName = fixture["display_name"]
		avatarUrl = fixture["avatar_url"]
		userId, err := userRepo.Create(ctx, fixture["idp_provider"], fixture["idp_subject"],
			fixture["username"], fixture["email"],
			&displayName, &avatarUrl)
		if err != nil {
			log.Fatal(err)
		}

		USER_FIXTURES = append(USER_FIXTURES, userId)
	}

	userOneProjects := []map[string]string{PROJECT_ONE, PROJECT_TWO}
	var description, skills string
	for _, fixture := range userOneProjects {
		description = fixture["description"]
		skills = fixture["skills"]
		id, err := projectRepo.Create(ctx,
			fixture["name"],
			&description, &skills,
			USER_FIXTURES[0])
		if err != nil {
			log.Fatal(err)
		}

		PROJECT_FIXTURES = append(PROJECT_FIXTURES, id)
	}

	userTwoProjects := []map[string]string{PROJECT_THREE}
	for _, fixture := range userTwoProjects {
		description = fixture["description"]
		skills = fixture["skills"]
		id, err := projectRepo.Create(ctx,
			fixture["name"],
			&description, &skills,
			USER_FIXTURES[1])
		if err != nil {
			log.Fatal(err)
		}

		PROJECT_FIXTURES = append(PROJECT_FIXTURES, id)
	}
}
