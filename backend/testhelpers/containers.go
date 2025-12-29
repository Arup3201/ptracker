package testhelpers

import (
	"context"
	"path/filepath"
	"time"

	keycloak "github.com/stillya/testcontainers-keycloak"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
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

func CreateKeycloakContainer(ctx context.Context) (*keycloak.KeycloakContainer, error) {
	keycloakContainer, err := keycloak.Run(context.Background(),
		"quay.io/keycloak/keycloak:26.4.7",
		testcontainers.WithWaitStrategy(wait.ForListeningPort("8080/tcp")),
		keycloak.WithContextPath("/auth"),
		keycloak.WithRealmImportFile("../testdata/ptracker-realm.json"),
		keycloak.WithAdminUsername("admin"),
		keycloak.WithAdminPassword("admin"),
	)
	return keycloakContainer, err
}
