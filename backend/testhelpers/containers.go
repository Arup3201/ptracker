package testhelpers

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:18-alpine",
		postgres.WithInitScripts(
			"../migrations/1_users_table.up.sql",
			"../migrations/2_sessions_table.up.sql",
			"../migrations/3_projects_table.up.sql",
			"../migrations/4_task_status_type.up.sql",
			"../migrations/5_tasks_table.up.sql",
			"../migrations/6_user_role_type.up.sql",
			"../migrations/7_roles_table.up.sql",
			"../migrations/8_assignees_table.up.sql",
			"../migrations/9_comments_table.up.sql",
		),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connString,
	}, nil
}
