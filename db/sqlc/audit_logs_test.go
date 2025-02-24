package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
)


func TestCreateAuditLog(t *testing.T) {

	email := createRandomEmailForAudits()

	userArgs := CreateUserParams {
		Username: utils.RandomString(29),
		Email: email,
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	AuditLogArgs := CreateAuditLogParams {
		UserEmail: user.Email,
		Action: "login",
		IpAddress: sql.NullString{String: "0.0.0.0", Valid: true},
	}

	auditLog, err := testQueries.CreateAuditLog(context.Background(), AuditLogArgs)
	require.NoError(t, err, "Failed to create audit log")
	require.NotEmpty(t, auditLog.ID, "AuditLog ID should not be empty")
}

func TestDeleteAuditLog(t *testing.T) {

	email := createRandomEmailForAudits()

	userArgs := CreateUserParams {
		Username: utils.RandomString(5),
		Email: email,
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	AuditlogArgs := CreateAuditLogParams {
		UserEmail: user.Email,
		Action: "login",
		IpAddress: sql.NullString{String: "0.0.0.1", Valid: true},
	}

	auditLog, err := testQueries.CreateAuditLog(context.Background(), AuditlogArgs)
	require.NoError(t, err, "Failed to create audit log")
	err = testQueries.DeleteAuditLog(context.Background(), auditLog.ID)
	require.NoError(t, err, "Failed to delete audit log")
}

func TestGetAuditLogByUserId(t *testing.T) {

	email := createRandomEmailForAudits()

	userArgs := CreateUserParams {
		Username: utils.RandomString(8),
		Email: email,
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	AuditLogArgs := CreateAuditLogParams {
		UserEmail: user.Email,
		Action: "login",
		IpAddress: sql.NullString{String: "0.0.0.2", Valid: true},
	}

	auditLog, err := testQueries.CreateAuditLog(context.Background(), AuditLogArgs)
	require.NoError(t, err, "Failed to create audit log")

	auditLogByUserId, err := testQueries.GetAuditLogsByUserEmail(context.Background(), auditLog.UserEmail)
	require.NoError(t, err, "Failed to get audit log by user ID")
	require.NotEmpty(t, auditLogByUserId, "Audit log should not be empty")
	require.Equal(t, auditLog.UserEmail, auditLogByUserId[0].UserEmail, "UserID should match")
}

func createRandomEmailForAudits() string {
	return fmt.Sprintf("audits-%s@example.com", uuid.New().String())
}