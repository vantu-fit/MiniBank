package token

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vantu-fit/master-go-be/utils"
)

func TestPasetoMaker(t *testing.T) {
	username := "user" + utils.RandomString(6)
	duration := time.Minute

	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)

	fmt.Println(token)

	payload , err := maker.VerifyToken(token) 
	require.NoError(t , err)
	
	jsonPayload , err := json.Marshal(payload)
	require.NoError(t , err)

	fmt.Println(string(jsonPayload))

}
