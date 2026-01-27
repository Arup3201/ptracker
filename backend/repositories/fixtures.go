package repositories

import "github.com/ptracker/domain"

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

var USER_FIXTURES = []domain.User{}
var PROJECT_FIXTURES = []domain.Project{}
