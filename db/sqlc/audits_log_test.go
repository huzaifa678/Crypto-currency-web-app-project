package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestCreateAuditLog(t *testing.T) {

	userArgs := CreateUserParams {
		Email: "exam000@example.com",
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	AuditLogArgs := CreateAuditLogParams {
		UserID: user.ID,
		Action: "login",
		IpAddress: sql.NullString{String: "0.0.0.0", Valid: true},
	}

	auditLog, err := testQueries.CreateAuditLog(context.Background(), AuditLogArgs)
	require.NoError(t, err, "Failed to create audit log")
	require.NotEmpty(t, auditLog.ID, "AuditLog ID should not be empty")
}

func TestDeleteAuditLog(t *testing.T) {
	userArgs := CreateUserParams {
		Email: "exam007@example.com",
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	AuditlogArgs := CreateAuditLogParams {
		UserID: user.ID,
		Action: "login",
		IpAddress: sql.NullString{String: "0.0.0.1", Valid: true},
	}

	auditLog, err := testQueries.CreateAuditLog(context.Background(), AuditlogArgs)
	require.NoError(t, err, "Failed to create audit log")
	err = testQueries.DeleteAuditLog(context.Background(), auditLog.ID)
	require.NoError(t, err, "Failed to delete audit log")
}

func TestGetAuditLogByUserId(t *testing.T) {

	userArgs := CreateUserParams {
		Email: "exam117@example.com",
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	AuditLogArgs := CreateAuditLogParams {
		UserID: user.ID,
		Action: "login",
		IpAddress: sql.NullString{String: "0.0.0.2", Valid: true},
	}

	auditLog, err := testQueries.CreateAuditLog(context.Background(), AuditLogArgs)
	require.NoError(t, err, "Failed to create audit log")

	auditLogByUserId, err := testQueries.GetAuditLogsByUserID(context.Background(), auditLog.UserID)
	require.NoError(t, err, "Failed to get audit log by user ID")
	require.NotEmpty(t, auditLogByUserId, "Audit log should not be empty")
	require.Equal(t, auditLog.UserID, auditLogByUserId[0].UserID, "UserID should match")
}