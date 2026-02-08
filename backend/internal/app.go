package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ptracker/internal/controllers"
	"github.com/ptracker/internal/controllers/middlewares"
	"github.com/ptracker/internal/infra"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/services"
	"github.com/rs/cors"
)

type app struct {
	_prefix string

	config  *Config
	handler *http.ServeMux

	authMiddleware      middlewares.Middleware
	rateLimitMiddleware middlewares.Middleware

	authController    interfaces.AuthController
	projectController interfaces.ProjectController
	taskController    interfaces.TaskController
	publicController  interfaces.PublicController

	notificationHandler http.Handler
}

func NewApp(config *Config,
	db interfaces.Execer,
	inMemory interfaces.InMemory,
	rateLimiter interfaces.RateLimiter,
	notifier *infra.WSNotifier) *app {

	handler := http.NewServeMux()

	app := &app{
		config:  config,
		handler: handler,
	}

	store := services.NewStorage(db, inMemory, rateLimiter)

	authService := services.NewAuthService(store,
		config.KeycloakURL,
		config.KeycloakRealm,
		config.KeycloakClientId,
		config.KeycloakClientSecret,
		config.KeycloakRedirectURI,
		config.EncryptionKey)
	projectService := services.NewProjectService(store, notifier)
	taskService := services.NewTaskService(store, notifier)
	publicService := services.NewPublicService(store)
	limitService := services.NewLimiterService(store)

	app.authMiddleware = middlewares.NewAuthMiddleware(authService)
	app.rateLimitMiddleware = middlewares.NewRateLimitMiddleware(limitService)

	app.authController = controllers.NewAuthController(authService, config.HomeURL)
	app.projectController = controllers.NewProjectController(projectService)
	app.taskController = controllers.NewTaskController(taskService)
	app.publicController = controllers.NewPublicController(publicService)

	// TODO: restrict origins
	app.notificationHandler = infra.NewNotificationHandler([]string{}, notifier)

	return app
}

func (a *app) attachMiddleware(method, pattern string, handler controllers.HTTPErrorHandler) {
	a.handler.Handle(method+" "+a._prefix+pattern, controllers.HTTPErrorHandler(
		a.authMiddleware.Handler(
			handler,
		),
	),
	)
}

func (a *app) AttachRoutes(prefix string) *app {
	a._prefix = prefix

	a.attachMiddleware("GET", "/auth/login", a.authController.Login)
	a.attachMiddleware("GET", "/auth/callback", a.authController.Callback)
	a.attachMiddleware("POST", "/auth/refresh", a.authController.Refresh)
	a.attachMiddleware("GET", "/auth/me", a.authController.Me)
	a.attachMiddleware("POST", "/auth/logout", a.authController.Logout)

	a.attachMiddleware("POST", "/projects", a.rateLimitMiddleware.Handler(
		controllers.HTTPErrorHandler(a.projectController.Create),
	))

	a.attachMiddleware("GET", "/dashboard/projects/created", a.projectController.ListRecentlyCreatedProjects)
	a.attachMiddleware("GET", "/dashboard/projects/joined", a.projectController.ListRecentlyJoinedProjects)
	a.attachMiddleware("GET", "/dashboard/tasks/assigned", a.taskController.ListAssignedTasks)
	a.attachMiddleware("GET", "/dashboard/tasks/unassigned", a.taskController.ListUnassignedTasks)

	a.attachMiddleware("GET", "/projects", a.projectController.List)
	a.attachMiddleware("GET", "/projects/{id}", a.projectController.Get)
	a.attachMiddleware("POST", "/projects/{id}/join-requests", a.publicController.JoinProject)
	a.attachMiddleware("GET", "/projects/{id}/join-requests", a.projectController.ListJoinRequests)
	a.attachMiddleware("PUT", "/projects/{id}/join-requests", a.projectController.RespondToJoinRequests)
	a.attachMiddleware("GET", "/projects/{id}/members", a.projectController.ListMembers)

	a.attachMiddleware("GET", "/projects/{project_id}/tasks", a.taskController.List)
	a.attachMiddleware("POST", "/projects/{project_id}/tasks", a.taskController.Create)
	a.attachMiddleware("GET", "/projects/{project_id}/tasks/{task_id}", a.taskController.Get)
	a.attachMiddleware("PUT", "/projects/{project_id}/tasks/{task_id}", a.taskController.Update)

	a.attachMiddleware("POST", "/projects/{project_id}/tasks/{task_id}/comments", a.taskController.AddComment)
	a.attachMiddleware("GET", "/projects/{project_id}/tasks/{task_id}/comments", a.taskController.ListComments)

	a.attachMiddleware("GET", "/public/projects", a.publicController.ListProjects)
	a.attachMiddleware("GET", "/public/projects/{id}", a.publicController.GetProject)

	// TODO: Add proper authentication method
	a.handler.Handle("GET /ws", a.notificationHandler)

	return a
}

func (a *app) Start() error {

	// TODO: some problem with not implementing http.Hijacker interface
	// when using with websocket
	// logging := middlewares.NewLoggingMiddleware()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{a.config.HomeURL},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	})

	// server
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", a.config.ServerHost, a.config.ServerPort),
		// Handler:      logging.Handler(cors.Handler(a.handler)),
		Handler:      cors.Handler(a.handler),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("[INFO] server starting at %s:%s\n", a.config.ServerHost, a.config.ServerPort)

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server listen and serve: %w", err)
	}

	return nil
}
