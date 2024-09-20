package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/longln/simplebank/db/mock"
	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/utils"
	"github.com/stretchr/testify/require"
)


func TestGetAccount(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()	// ensure which function expected to called were called
	
	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
	GetAccount(gomock.Any(), gomock.Eq(account.ID)).
	Times(1).
	Return(account, nil)
	

	server := NewServer(store)
	// test mock API, we don't have to start server, we can record by httptest and compare result
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}


func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := []struct {
		name	string
		accountID	int64
		buildStubs	func(store *mockdb.MockStore)
		checkResponse	func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// TODO: add test data
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore){

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){

			},
		},

		{

		},
	}

	for _, tc := range(testCases) {
		
	}
}


func randomAccount() db.Account {
	return db.Account{
		ID: int64(utils.RandomInt(1, 1000)),
		Owner: utils.RandomOwner(),
		Balance: utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}


func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}