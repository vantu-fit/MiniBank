package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vantu-fit/master-go-be/db/mock"
	db "github.com/vantu-fit/master-go-be/db/sqlc"
	"github.com/vantu-fit/master-go-be/utils"
)

func TestGetAccountAPI(t *testing.T) {
	
	account := randomAccount()

	testCase := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Ok",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "notfound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "internalserver",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, pgx.ErrTxCommitRollback)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "invalidid",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.All()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)

		server, err := NewServer(store)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()

		url := fmt.Sprintf("/accounts/%d", tc.accountID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		addAuthorization(t, request, server.maker, authorizationTypeBearer, account.Owner, utils.DepositorRole ,  time.Minute)

		server.router.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       int64(utils.RandomInt(1, 1000)),
		Owner:    utils.RandomString(6),
		Balance:  int64(utils.RandomInt(500, 1000)),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, gotAccount, account)

}

func TestListAccount(t *testing.T) {
	req := ListAccountRequest{
		PageID:   5,
		PageSize: 5,
	}
	testCase := []struct {
		name          string
		reqBody       ListAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			reqBody: req,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return([]db.Account{}, nil)

			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {

			},
		},
	}

	for i := range testCase {

		user, _, err := randomUser()
		require.NoError(t, err)

		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server, err := NewServer(store)
			require.NoError(t, err)

			tc.buildStubs(store)

			path := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", tc.reqBody.PageID, tc.reqBody.PageSize)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)

			addAuthorization(t, request, server.maker, authorizationTypeBearer, user.Username, user.Role ,  time.Minute)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}
