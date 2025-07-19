package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/Deadly-Smile/simple-bank/database/utils"
	"github.com/stretchr/testify/require"
)

// CreateRandomUser creates a random user for testing purposes
func CreateRandomUser(t *testing.T) User {
	ctx := context.Background()
	username := utils.RandomString(6)
	user, err := testQueries.CreateUser(ctx, username)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, username, user.Username)
	require.NotZero(t, user.ID)
	require.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	return user
}

// TestCreateUser tests the creation of a user
func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

// TestGetUser tests the retrieval of a user by ID
func TestGetUser(t *testing.T) {
	ctx := context.Background()

	// Create a user first
	createdUser := CreateRandomUser(t)

	// Now get the user by ID
	user, err := testQueries.GetUser(ctx, createdUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, createdUser.ID, user.ID)
	require.Equal(t, createdUser.Username, user.Username)
	require.WithinDuration(t, createdUser.CreatedAt, user.CreatedAt, time.Second)
}

// TestDeleteUser tests the deletion of a user
func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	// Create a user first
	createdUser := CreateRandomUser(t)

	// Now delete the user
	err := testQueries.DeleteUser(ctx, createdUser.ID)
	require.NoError(t, err)

	// After deletion, trying to get the user should return an error
	user, err := testQueries.GetUser(ctx, createdUser.ID)
	require.Error(t, err)
	require.Empty(t, user)
}

// TestUpdateUser tests the update functionality of a user
func TestUpdateUser(t *testing.T) {
	cxt := context.Background()
	createdUser := CreateRandomUser(t)

	arg := UpdateUserParams{
		ID:       createdUser.ID,
		Username: utils.RandomString(6),
	}
	updatedUser, err := testQueries.UpdateUser(cxt, arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, createdUser.ID, updatedUser.ID)
	require.NotEqual(t, createdUser.Username, updatedUser.Username)
	require.WithinDuration(t, createdUser.CreatedAt, updatedUser.CreatedAt, time.Second)
}

// TestListUsers tests the listing of users with pagination
func TestListUsers(t *testing.T) {
	ctx := context.Background()

	// Create multiple users
	for i := 0; i < 10; i++ {
		CreateRandomUser(t)
	}

	// List users with a limit and offset
	arg := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}
	users, err := testQueries.ListUsers(ctx, arg)
	require.NoError(t, err)
	require.Len(t, users, int(arg.Limit))

	for _, user := range users {
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.Username)
		require.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	}
}
