package gapi

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vantu-fit/master-go-be/db/mock"
	db "github.com/vantu-fit/master-go-be/db/sqlc"
	"github.com/vantu-fit/master-go-be/pb"
	"github.com/vantu-fit/master-go-be/token"
	"github.com/vantu-fit/master-go-be/utils"
	"github.com/vantu-fit/master-go-be/worker"
	"google.golang.org/grpc/metadata"
)

type eqUpdateUserParamsMatcher struct {
	arg      db.UpdateUserParams
	password string
	user     db.User
}

func (e eqUpdateUserParamsMatcher) Matches(x interface{}) bool {
	return true
	// In case, some value is nil
	arg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	fmt.Println(">>> check password")

	e.arg.HashedPassword = arg.HashedPassword
	fmt.Println(">>> deep equal")

	if !reflect.DeepEqual(e.arg, arg) {
		return false
	}

	return true
}

func (e eqUpdateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v ", e.arg, e.password)
}

func EqUpdateUserTxParamsMatcher(arg db.UpdateUserParams, password string, user db.User) gomock.Matcher {
	return eqUpdateUserParamsMatcher{arg, password, user}
}

func TestUpdateUserAPI(t *testing.T) {

	user, password, err := randomUser()
	_ = password

	userUpdate := db.User{
		Username: user.Username,
		FullName: "updated" + user.FullName,
		Email:    "updated" + user.Email,
	}

	require.NoError(t, err)

	testCase := []struct {
		name          string
		body          *pb.UpdateUserRequest
		buildContext  func(t *testing.T, maker token.Maker) context.Context
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "UnAuthenticated",
			// muon gui len body
			body: &pb.UpdateUserRequest{},
			// key qua mong doi
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)

			},
			buildContext: func(t *testing.T, maker token.Maker) context.Context {
				return context.Background()
			},
			// kiem tra ket qua
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NotEmpty(t, err)
			},
		},
		{
			name: "OK",
			// muon gui len body
			body: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &userUpdate.FullName,
				Email:    &userUpdate.Email,
			},
			// key qua mong doi

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: pgtype.Text{
						String: userUpdate.FullName,
						Valid:  true,
					},
					Email: pgtype.Text{
						String: userUpdate.Email,
						Valid:  true,
					},
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.User{
						Username:          arg.Username,
						FullName:          arg.FullName.String,
						Email:             arg.Email.String,
						HashedPassword:    user.HashedPassword,
						PasswordChangedAt: user.PasswordChangedAt,
						CreatedAt:         user.CreatedAt,
						IsEmailVerified:   user.IsEmailVerified,
						Role:              user.Role,
					}, nil)

			},
			buildContext: func(t *testing.T, maker token.Maker) context.Context {
				ctx := context.Background()
				accessToken, _, err := maker.CreateToken(user.Username, user.Role , time.Minute*5)
				require.NoError(t, err)
				bearerToken := fmt.Sprintf("%s %s", authorizationTypeBearer, accessToken)
				md := metadata.MD{
					authorizationHeader: []string{
						bearerToken,
					},
				}
				return metadata.NewIncomingContext(ctx, md)
			},

			// kiem tra ket qua
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {

				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, userUpdate.Username, res.User.Username)
				require.Equal(t, userUpdate.Email, res.User.Email)
				require.Equal(t, userUpdate.FullName, res.User.FullName)
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

			redisOpt := asynq.RedisClientOpt{}

			taskDitributor := worker.NewRedisTaskDistrubutor(redisOpt)

			maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
			require.NoError(t, err)

			tc.buildStubs(store)
			// tao server voi mockdb cho qua trinh test
			server, err := NewServer(store, taskDitributor, config)
			require.NoError(t, err)
			// tao recorder de chua ket qua sau khi gui request de so sanh voi check response
			ctx := tc.buildContext(t, maker)
			res, err := server.UpdateUser(ctx, tc.body)

			tc.checkResponse(t, res, err)
		})
	}

}
