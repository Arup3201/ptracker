package models

import (
	"database/sql"
	"fmt"
)

type Role struct {
	ProjectId string
	UserId    string
	Role      string
}

type RoleStore struct {
	DB        *sql.DB
	UserId    string
	ProjectId string
}

func (rs *RoleStore) Get() (*Role, error) {
	var role Role
	err := rs.DB.QueryRow(
		"SELECT "+
			"project_id, user_id, role "+
			"FROM roles "+
			"WHERE user_id=($1) AND project_id=($2)",
		rs.UserId, rs.ProjectId,
	).Scan(&role.ProjectId, &role.UserId, &role.Role)
	if err != nil {
		return nil, fmt.Errorf("postgres query user role: %w", err)
	}

	return &role, nil
}

func (rs *RoleStore) CanAccess() (bool, error) {
	var userRole string
	err := rs.DB.QueryRow(
		"SELECT "+
			"role "+
			"FROM roles "+
			"WHERE user_id=($1) AND project_id=($2)",
		rs.UserId, rs.ProjectId,
	).Scan(&userRole)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("postgres query user role: %w", err)
	}

	if userRole == ROLE_MEMBER || userRole == ROLE_OWNER {
		return true, nil
	}

	return false, nil
}

func (rs *RoleStore) CanEdit() (bool, error) {
	var userRole string
	err := rs.DB.QueryRow(
		"SELECT "+
			"role "+
			"FROM roles "+
			"WHERE user_id=($1) AND project_id=($2)",
		rs.UserId, rs.ProjectId,
	).Scan(&userRole)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("postgres query user role: %w", err)
	}

	if userRole == ROLE_OWNER {
		return true, nil
	}

	return false, nil
}
