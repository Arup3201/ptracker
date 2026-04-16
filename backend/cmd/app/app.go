package app

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/ptracker/api"
	"github.com/ptracker/auth"
	"github.com/ptracker/auth/manual"
	"github.com/ptracker/core"
	"github.com/ptracker/core/assignees"
	"github.com/ptracker/core/comments"
	"github.com/ptracker/core/members"
	"github.com/ptracker/core/projects"
	"github.com/ptracker/core/requests"
	"github.com/ptracker/core/tasks"
	"github.com/ptracker/core/users"
	"github.com/redis/go-redis/v9"
	"github.com/resend/resend-go/v3"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type patternWithHandler struct {
	method, pattern string
	handler         api.HTTPErrorHandler
}

type App struct {
	// CORS allowed origins
	AllowedCrossOrigins []string

	config  *Config
	handler *http.ServeMux
}

func NewApp(
	// API paths start with Prefix: /v1, /v2, /api, /v1/api
	prefix string,

	config *Config,
	db *gorm.DB,
	redis *redis.Client,
	privateKey *rsa.PrivateKey,
	frontendVerifyUrl string,
) *App {
	handler := http.NewServeMux()

	memberRepo := members.NewMemberRepository(db)
	joinRepo := requests.NewJoinRepository(db)
	accountRepo := manual.NewManualAccountRepository(db)
	userRepo := users.NewUserRepository(db)
	assigneeRepo := assignees.NewAssigneeRepository(db)
	commentRepo := comments.NewCommentRepository(db)
	projectRepo := projects.NewProjectRepository(db)
	taskRepo := tasks.NewTaskRepository(db)
	txManager := core.NewTxManager(db)
	tokenStore := auth.NewTokenStore(redis)

	memberService := members.NewMemberService(memberRepo)
	joinService := requests.NewJoinRequestService(
		txManager,
		joinRepo,
		memberRepo,
	)
	userService := users.NewUserService(userRepo)
	assigneeService := assignees.NewAssigneeService(
		memberRepo,
		assigneeRepo)
	commentService := comments.NewCommentService(
		commentRepo,
		memberRepo)
	projectService := projects.NewProjectService(
		txManager,
		projectRepo,
		memberRepo)
	taskService := tasks.NewTaskService(
		taskRepo,
		memberRepo,
		assigneeRepo,
	)
	registerService := manual.NewRegisterService(txManager, accountRepo, userRepo)
	tokenService := auth.NewTokenService(tokenStore, TOKEN_ISSUER, privateKey)
	emailService := manual.NewEmailService(accountRepo)

	resendClient := resend.NewClient(config.ResendApiKey)
	authApi := api.NewAuthApi(
		registerService,
		tokenService,
		emailService,
		userService,
		resendClient,
		frontendVerifyUrl,
	)
	projectApi := api.NewProjectApi(
		projectService,
		userService,
		memberService,
		joinService,
	)
	taskApi := api.NewTaskApi(
		taskService,
		assigneeService,
		commentService,
	)

	patternWithHandlers := []patternWithHandler{
		// Auth
		{
			method:  "POST",
			pattern: "/auth/register",
			handler: authApi.Register,
		},
		{
			method:  "POST",
			pattern: "/auth/login",
			handler: authApi.Login,
		},
		{
			method:  "POST",
			pattern: "/auth/refresh",
			handler: authApi.Refresh,
		},
		{
			method:  "POST",
			pattern: "/auth/logout",
			handler: authApi.Logout,
		},

		// List APIs
		{
			method:  "GET",
			pattern: "/projects",
			handler: projectApi.ListMyProjects,
		},
		{
			method:  "GET",
			pattern: "/dashboard/projects/created",
			handler: projectApi.ListRecentlyCreated,
		},
		{
			method:  "GET",
			pattern: "/dashboard/projects/joined",
			handler: projectApi.ListRecentlyJoined,
		},
		{
			method:  "GET",
			pattern: "/projects/{project_id}/tasks",
			handler: taskApi.List,
		},
		{
			method:  "GET",
			pattern: "/dashboard/tasks/assigned",
			handler: taskApi.ListAssignedTasks,
		},
		{
			method:  "GET",
			pattern: "/dashboard/tasks/unassigned",
			handler: taskApi.ListUnassignedTasks,
		},
		{
			method:  "GET",
			pattern: "/projects/{id}/members",
			handler: projectApi.ListMembers,
		},
		{
			method:  "GET",
			pattern: "/projects/{id}/join-requests",
			handler: projectApi.ListJoinRequests,
		},
		{
			method:  "GET",
			pattern: "/projects/{project_id}/tasks/{task_id}/comments",
			handler: taskApi.ListComments,
		},
		{
			method:  "GET",
			pattern: "/public/projects",
			handler: projectApi.ListPublic,
		},
		// Get Single Instance APIs
		{
			method:  "GET",
			pattern: "/projects/{id}",
			handler: projectApi.Get,
		},
		{
			method:  "GET",
			pattern: "/projects/{project_id}/tasks/{task_id}",
			handler: taskApi.Get,
		},
		// Create APIs
		{
			method:  "POST",
			pattern: "/projects",
			handler: projectApi.Create,
		},
		{
			method:  "POST",
			pattern: "/projects/{project_id}/tasks",
			handler: taskApi.Create,
		},
		{
			method:  "POST",
			pattern: "/projects/{id}/join-requests",
			handler: projectApi.AddJoinRequest,
		},
		{
			method:  "POST",
			pattern: "/projects/{project_id}/tasks/{task_id}/comments",
			handler: taskApi.AddComment,
		},
		// Update Instance APIs
		{
			method:  "PUT",
			pattern: "/projects/{project_id}/tasks/{task_id}",
			handler: taskApi.Update,
		},
		{
			method:  "PUT",
			pattern: "/projects/{id}/join-requests",
			handler: projectApi.RespondToJoinRequest,
		},
	}

	for _, h := range patternWithHandlers {
		handler.Handle(
			h.method+" "+prefix+h.pattern,
			api.HTTPErrorHandler(h.handler),
		)
	}

	return &App{
		config:  config,
		handler: handler,
	}
}

func (app *App) Start() error {

	cors := cors.New(cors.Options{
		AllowedOrigins:   app.AllowedCrossOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	})

	// server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", app.config.Host, app.config.Port),
		Handler:      cors.Handler(app.handler),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("[INFO] server starting at %s:%s\n", app.config.Host, app.config.Port)

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server listen and serve: %w", err)
	}

	return nil
}
