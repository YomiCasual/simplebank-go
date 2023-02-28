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

	myPassword := "password"
	
	hashed_password, err := lib.PasswordCrypt().HashPassword(myPassword)

	arg := sqlc.CreateUserParams{
		Username: lib.RandomOwner(),
		HashedPassword: hashed_password ,
		FullName: lib.RandomOwner(),
		Email: lib.RandomEmail(),
	}


	user, err := testQueries.CreateUser(context.Background(), arg)
	log.Println("user:", user, arg)


	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require. NotEqual(t, myPassword, user.HashedPassword)
	require. Equal(t, arg.HashedPassword, user.HashedPassword)

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