package token

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vantu-fit/master-go-be/utils"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := "user" + utils.RandomString(6)
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, time.Now(), payload.ExpiredAt, time.Minute)
	require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)

}

func TestExpriedToken(t *testing.T) { 
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := "user" + utils.RandomString(6)
	duration := -1 * time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Equal(t , err , ErrorExpiredToken)
	require.Empty(t, payload)
	
	
}

func TestInvalidToken (t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	
	username := "user" + utils.RandomString(6)
	duration :=  time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(strings.Join((strings.Split(token, ""))[0:10], ""))
	require.Error(t, err )
	require.Equal(t , err,  ErrInvalidToken)
	require.Empty(t, payload)
	
}
