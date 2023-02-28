package test

import (
	"context"
	"log"
	"simplebank/db/sqlc"
	lib "simplebank/libs"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateUser(t *testing.T) (sqlc.User, sqlc.CreateUserParams, error ) {
	arg := sqlc.CreateUserParams{
		Username: lib.RandomOwner(),
		HashedPassword: lib.RandomString(6),
		FullName: lib.RandomOwner(),
		Email: lib.RandomString(10),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	log.Println("user:", user, arg)


	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user,arg,  err
}

func TestCreateUser(t *testing.T) {
	CreateUser(t)
}


func TestListingUser(t *testing.T) {
	CreateUser(t)


	users, err := testQueries.ListUsers(context.Background(), sqlc.ListUsersParams{
		Limit: 5,
		Offset: 1,
	})

	
	require.NoError(t, err)
	require.NotEmpty(t, users)

}

func TestGettingUserById(t *testing.T) {
	createdUser, _, _ := CreateUser(t)


	findUser, err := testQueries.GetUser(context.Background(), createdUser.ID)

	
	require.NoError(t, err)
	require.NotEmpty(t, createdUser.Email, findUser.Email)
	require.NotEmpty(t, createdUser.ID, findUser.ID)

}