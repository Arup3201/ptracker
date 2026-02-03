package interfaces

import "net/http"

type AuthController interface {
	Login(w http.ResponseWriter, r *http.Request) error
	Callback(w http.ResponseWriter, r *http.Request) error
	Refresh(w http.ResponseWriter, r *http.Request) error
	Me(w http.ResponseWriter, r *http.Request) error
	Logout(w http.ResponseWriter, r *http.Request) error
}

type ProjectController interface {
	Create(w http.ResponseWriter, r *http.Request) error
	List(w http.ResponseWriter, r *http.Request) error
	Get(w http.ResponseWriter, r *http.Request) error
	ListJoinRequests(w http.ResponseWriter, r *http.Request) error
	RespondToJoinRequests(w http.ResponseWriter, r *http.Request) error
	ListMembers(w http.ResponseWriter, r *http.Request) error
}

type TaskController interface {
	List(w http.ResponseWriter, r *http.Request) error
	Create(w http.ResponseWriter, r *http.Request) error
	Get(w http.ResponseWriter, r *http.Request) error
	Update(w http.ResponseWriter, r *http.Request) error
}

type PublicController interface {
	ListProjects(w http.ResponseWriter, r *http.Request) error
	GetProject(w http.ResponseWriter, r *http.Request) error
	JoinProject(w http.ResponseWriter, r *http.Request) error
}
