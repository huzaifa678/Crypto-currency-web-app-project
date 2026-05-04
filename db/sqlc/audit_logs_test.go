package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateAuditLog(t *testing.T) {

	email := createRandomEmailForAudits()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:        email,
		PasswordHash: "9009909",
		Role:         "user",
		IsVerified:   true,
	}

	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	AuditLogArgs := CreateAuditLogParams{
		Username:  user.Username,
		UserEmail: user.Email,
		Action:    "login",
		IpAddress: "0.0.0.0",
	}

	auditLog, err := testStore.CreateAuditLog(context.Background(), AuditLogArgs)
	require.NoError(t, err, "Failed to create audit log")
	require.NotEmpty(t, auditLog.ID, "AuditLog ID should not be empty")
}

func TestDeleteAuditLog(t *testing.T) {

	email := createRandomEmailForAudits()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:        email,
		PasswordHash: "9009909",
		Role:         "user",
		IsVerified:   true,
	}

	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	AuditlogArgs := CreateAuditLogParams{
		Username:  user.Username,
		UserEmail: user.Email,
		Action:    "login",
		IpAddress: "0.0.0.1",
	}

	auditLog, err := testStore.CreateAuditLog(context.Background(), AuditlogArgs)
	require.NoError(t, err, "Failed to create audit log")
	err = testStore.DeleteAuditLog(context.Background(), auditLog.ID)
	require.NoError(t, err, "Failed to delete audit log")
}

func TestGetAuditLogByUserId(t *testing.T) {

	email := createRandomEmailForAudits()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:        email,
		PasswordHash: "9009909",
		Role:         "user",
		IsVerified:   true,
	}

	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	AuditLogArgs := CreateAuditLogParams{
		Username:  user.Username,
		UserEmail: user.Email,
		Action:    "login",
		IpAddress: "0.0.0.2",
	}

	auditLog, err := testStore.CreateAuditLog(context.Background(), AuditLogArgs)
	require.NoError(t, err, "Failed to create audit log")

	auditLogByUserId, err := testStore.GetAuditLogsByUserEmail(context.Background(), auditLog.UserEmail)
	require.NoError(t, err, "Failed to get audit log by user ID")
	require.NotEmpty(t, auditLogByUserId, "Audit log should not be empty")
	require.Equal(t, auditLog.UserEmail, auditLogByUserId[0].UserEmail, "UserID should match")
}

func createRandomEmailForAudits() string {
	return fmt.Sprintf("audits-%s@example.com", uuid.New().String())
}
