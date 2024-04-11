package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	_ "time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vantu-fit/master-go-be/db/mock"
	db "github.com/vantu-fit/master-go-be/db/sqlc"
	"github.com/vantu-fit/master-go-be/utils"
)

// custome mathcher
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// In case, some value is nil
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(arg.HashedPassword, e.password)
	if err != nil {
		fmt.Println("khong bang")
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v ", e.arg, e.password)
}

func EqCreateUserParamsMatcher(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {

	user, password, err := randomUser()
	require.NoError(t, err)

	arg := db.CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}

	testCase := []struct {
		name          string
		arg           createUsertRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			// muon gui len body
			arg: createUsertRequest{
				Username: arg.Username,
				Password: password,
				FullName: arg.FullName,
				Email:    arg.Email,
			},
			// key qua mong doi
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParamsMatcher(arg, password)).Times(1).Return(user, nil)
			},
			// kiem tra ket qua
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotUser createUserReponse
				err = json.Unmarshal(data, &gotUser)
				require.NoError(t, err)
				fmt.Println(gotUser)
				require.Equal(t, gotUser, createUserReponse{
					Username:          user.Username,
					FullName:          user.FullName,
					Email:             user.Email,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
				})

			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		// tao cntroller cho mockdb
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		// tao mockdb
		store := mockdb.NewMockStore(ctrl)
		// ket qua mong doi
		tc.buildStubs(store)
		// tao server voi mockdb cho qua trinh test
		server, err := NewServer(store)
		require.NoError(t, err)
		// tao recorder de chua ket qua sau khi gui request de so sanh voi check response
		recorder := httptest.NewRecorder()
		// tao url de check aiu
		url := "/users"
		// tao body de gui len request
		jsonBody, err := json.Marshal(tc.arg)
		require.NoError(t, err)

		body := bytes.NewReader(jsonBody)
		// gui request
		request, err := http.NewRequest(http.MethodPost, url, body)
		require.NoError(t, err)
		// luu ket qua vao recorder
		server.router.ServeHTTP(recorder, request)
		// check ket qua
		tc.checkResponse(t, recorder)
	}

}

func randomUser() (db.User, string, error) {
	password := utils.RandomString(6)
	hashedPassword, err := utils.HashedPassword(password)
	if err != nil {
		return db.User{}, password, err
	}
	return db.User{
		Username:       "user" + utils.RandomString(6),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomString(10),
		Email:          utils.RandomEmail(),
	}, password, nil

}

// func TestDuration(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	store := mockdb.NewMockStore(ctrl)
// 	server, err := NewServer(store)
// 	require.NoError(t, err)
// 	fmt.Println("duration" + server.config.AccessTokenDuration)
// 	fmt.Println("symmestrictkey : "+ server.config.TokenSymmetricKey)

// 	duration, err := time.ParseDuration(server.config.AccessTokenDuration)
// 	require.NoError(t, err)
// 	fmt.Println(duration)
// }
