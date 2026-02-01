package internal

import (
	"fmt"
	"os"

	"github.com/ptracker/internal/constants"
)

type Config struct {
	ServerHost           string
	ServerPort           string
	DbHost               string
	DbPort               string
	DbUser               string
	DbPass               string
	DbName               string
	KeycloakURL          string
	KeycloakRealm        string
	KeycloakClientId     string
	KeycloakClientSecret string
	KeycloakRedirectURI  string
	EncryptionKey        string
	HomeURL              string
}

func (c *Config) Load() error {
	c.ServerHost = os.Getenv(constants.ENV_SERVER_HOST)
	if c.ServerHost == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_SERVER_HOST)
	}
	c.ServerPort = os.Getenv(constants.ENV_SERVER_PORT)
	if c.ServerPort == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_SERVER_PORT)
	}
	c.DbHost = os.Getenv(constants.ENV_DB_HOST)
	if c.DbHost == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_DB_HOST)
	}
	c.DbUser = os.Getenv(constants.ENV_DB_USER)
	if c.DbUser == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_DB_USER)
	}
	c.DbPort = os.Getenv(constants.ENV_DB_PORT)
	if c.DbPort == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_DB_PORT)
	}
	c.DbPass = os.Getenv(constants.ENV_DB_PASS)
	if c.DbPass == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_DB_PASS)
	}
	c.DbName = os.Getenv(constants.ENV_DB_NAME)
	if c.DbName == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_DB_NAME)
	}
	c.KeycloakURL = os.Getenv(constants.ENV_KEYCLOAK_URL)
	if c.KeycloakURL == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_KEYCLOAK_URL)
	}
	c.KeycloakRealm = os.Getenv(constants.ENV_KEYCLOAK_REALM)
	if c.KeycloakRealm == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_KEYCLOAK_REALM)
	}
	c.KeycloakClientId = os.Getenv(constants.ENV_KEYCLOAK_CLIENT_ID)
	if c.KeycloakClientId == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_KEYCLOAK_CLIENT_ID)
	}
	c.KeycloakClientSecret = os.Getenv(constants.ENV_KEYCLOAK_CLIENT_SECRET)
	if c.KeycloakClientSecret == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_KEYCLOAK_CLIENT_SECRET)
	}
	c.KeycloakRedirectURI = os.Getenv(constants.ENV_KEYCLOAK_REDIRECT_URI)
	if c.KeycloakRedirectURI == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_KEYCLOAK_REDIRECT_URI)
	}
	c.HomeURL = os.Getenv(constants.ENV_HOME_URL)
	if c.HomeURL == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_HOME_URL)
	}
	c.EncryptionKey = os.Getenv(constants.ENV_ENCRYPTION_SECRET)
	if c.EncryptionKey == "" {
		return fmt.Errorf("environment variable '%s' missing", constants.ENV_ENCRYPTION_SECRET)
	}

	return nil
}
