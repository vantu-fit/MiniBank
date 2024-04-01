package db

import (
	"context"
	"database/sql"
	_ "fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vantu-fit/master-go-be/utils"
)

func createRandomUser(t *testing.T) User {
	hashedPassword , err := utils.HashedPassword(utils.RandomString(6))
	require.NoError(t , err)
	arg := CreateUserParams{
		Username:       "user" + utils.RandomString(6),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomString(6),
		Email:          utils.RandomEmail(),
	}
	user , err := testQueries.CreateUser(context.Background() , arg)
	require.NoError(t , err)
	require.NotEmpty(t , user)
	require.Equal(t , arg.Username , user.Username)
	require.Equal(t , arg.FullName , user.FullName)
	require.Equal(t , arg.Email , user.Email)
	require.NotZero(t , user.CreatedAt)
	require.False(t , user.IsEmailVerified)
	require.NotZero(t , user.PasswordChangedAt)
	return user
}

func TestCreateUser(t*testing.T){
	hashedPassword , err := utils.HashedPassword(utils.RandomString(6))
	require.NoError(t , err)
	arg := CreateUserParams{
		Username:       "user" + utils.RandomString(6),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomString(6),
		Email:          utils.RandomEmail(),
	}
	user2 , err := testQueries.CreateUser(context.Background() , arg)
	require.NoError(t , err)
	require.NotEmpty(t , user2)
	require.Equal(t , arg.Username , user2.Username)
	require.Equal(t , arg.FullName , user2.FullName)
	require.Equal(t , arg.Email , user2.Email)
	require.NotZero(t , user2.CreatedAt)
	require.False(t , user2.IsEmailVerified)
	require.NotZero(t , user2.PasswordChangedAt)

}

func TestGetUser(t *testing.T)  {
	user1 := createRandomUser(t)

	user2 , err := testQueries.GetUser(context.Background() , user1.Username)
	require.NoError(t , err)
	require.NotEmpty(t , user2)
	require.Equal(t , user1.Username , user2.Username)
	require.Equal(t , user1.FullName , user2.FullName)
	require.Equal(t , user1.Email , user2.Email)
	require.Equal(t , user1.CreatedAt , user2.CreatedAt)
	require.Equal(t , user1.IsEmailVerified , user2.IsEmailVerified)
	require.Equal(t , user1.PasswordChangedAt , user2.PasswordChangedAt)
	require.WithinDuration(t , user1.CreatedAt , user2.CreatedAt , time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t) 



	arg := UpdateUserParams{
		HashedPassword: sql.NullString{ String: utils.RandomString(10) , Valid: true},
		PasswordChangedAt: sql.NullTime{ Time: time.Now() , Valid: true},
		FullName: sql.NullString{ String: utils.RandomString(6) , Valid: true},
		Email: sql.NullString{ String: utils.RandomEmail() , Valid: true},
		IsEmailVerified: sql.NullBool{ Bool: true , Valid: true},
		Username: user1.Username,
	}
	user2 , err := testQueries.UpdateUser(context.Background() , arg)
	require.NoError(t , err)
	require.NotEmpty(t , user2)
	require.Equal(t , user1.Username , user2.Username)
	require.Equal(t , arg.FullName.String , user2.FullName)
	require.Equal(t , arg.Email.String , user2.Email)
	require.Equal(t , arg.IsEmailVerified.Bool , user2.IsEmailVerified)
	require.Equal(t , arg.HashedPassword.String , user2.HashedPassword)
	require.WithinDuration(t , user2.PasswordChangedAt , arg.PasswordChangedAt.Time , time.Second)

}


