package gapi

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vantu-fit/master-go-be/db/mock"
	db "github.com/vantu-fit/master-go-be/db/sqlc"
	"github.com/vantu-fit/master-go-be/pb"
	"github.com/vantu-fit/master-go-be/utils"
	"github.com/vantu-fit/master-go-be/worker"
	mockwk "github.com/vantu-fit/master-go-be/worker/mock"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (e eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	// In case, some value is nil
	arg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	fmt.Println(">>> check password")

	err := utils.CheckPassword(arg.HashedPassword, e.password)
	if err != nil {
		fmt.Println("khong bang")
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	fmt.Println(">>> deep equal")

	if !reflect.DeepEqual(e.arg.CreateUserParams, arg.CreateUserParams) {
		return false
	}

	// call after create fn here
	err = arg.AfterCreate(e.user)
	if err != nil {
		return false
	}

	return true
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v ", e.arg, e.password)
}

func EqCreateUserTxParamsMatcher(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func TestCreateUserAPI(t *testing.T) {

	user, password, err := randomUser()
	require.NoError(t, err)

	testCase := []struct {
		name          string
		body          *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "Ok",
			// muon gui len body
			body: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			// key qua mong doi
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)

				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username:       user.Username,
						FullName:       user.FullName,
						Email:          user.Email,
						HashedPassword: user.HashedPassword,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParamsMatcher(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)

			},
			// kiem tra ket qua
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				creatdeUser := res.User
				require.Equal(t, user.Username, creatdeUser.Username)
				require.Equal(t, user.Email, creatdeUser.Email)
				require.Equal(t, user.FullName, creatdeUser.FullName)
			},
		},
	}

	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			// tao cntroller cho mockdb
			ctrlMockStore := gomock.NewController(t)
			defer ctrlMockStore.Finish()
			// tao mockdb
			store := mockdb.NewMockStore(ctrlMockStore)

			ctrlMockTask := gomock.NewController(t)
			defer ctrlMockTask.Finish()

			taskDitributor := mockwk.NewMockTaskDistributor(ctrlMockTask)

			tc.buildStubs(store, taskDitributor)
			// tao server voi mockdb cho qua trinh test
			server, err := NewServer(store, taskDitributor, config)
			require.NoError(t, err)
			// tao recorder de chua ket qua sau khi gui request de so sanh voi check response
			res, err := server.CreateUser(context.Background(), tc.body)

			tc.checkResponse(t, res, err)
		})
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
		Role:           utils.DepositorRole,
	}, password, nil

}
