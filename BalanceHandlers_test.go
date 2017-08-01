package main

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"encoding/json"
	"github.com/GlynOwenHanmer/GOHMoney"
	"github.com/GlynOwenHanmer/GOHMoneyDB"
	"net/http"
	"bytes"
	"time"
	"fmt"
	"github.com/lib/pq"
)

func Test_AccountBalances(t *testing.T) {
	req := httptest.NewRequest("GET", "/account/1/balances", nil)
	w := httptest.NewRecorder()
	NewRouter().ServeHTTP(w, req)
	resp := w.Result()
	expectedCode := http.StatusOK
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	balances := GOHMoneyDB.Balances{}
	err = json.Unmarshal(body, &balances)
	if err != nil {
		t.Errorf("Error unmarshalling response body to Balances\nError: %s\nBody: %s", err.Error(), body)
	}
	minBalancesLength := 91
	if len(balances) < minBalancesLength {
		t.Fatalf("Expected balances min length %d, got length %d", minBalancesLength, len(balances))
	}
	expectedId := uint(1)
	actualId := balances[0].Id
	if expectedId != actualId {
		t.Errorf(`Unexpected Id.\nExpected: %d\nActual  : %d`, expectedId, actualId)
	}
	expectedAmount := float32(636.42)
	actualAmount := balances[0].Amount
	if actualAmount != expectedAmount {
		t.Errorf("first balance, expected balance amount of %f but got %f", expectedAmount, actualAmount)
	}
	expectedDate, err := parseDateString("2016-06-17")
	if err != nil {
		t.Fatalf("Error parsing date string for use in tests. Error: %s", err.Error())
	}
	actualDate := balances[0].Date
	if !expectedDate.Equal(actualDate) {
		t.Errorf("first balance, expected date of %s but got %s", urlFormatDateString(expectedDate), urlFormatDateString(actualDate))
	}
}

func Test_BalanceCreate(t *testing.T) {
	router := NewRouter()
	endpoint := `/balance/create`

	invalidData := []byte(`INVALID BODY`)
	req := httptest.NewRequest("POST", endpoint, bytes.NewReader(invalidData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()
	expectedCode := http.StatusBadRequest
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, resp.StatusCode)
	}

	type accountBalance struct {
		AccountId int `json:"account_id"`
		GOHMoney.Balance `json:"balance"`
	}

	type testSet struct {
		newBalance     accountBalance
		expectedStatus int
		expectJsonDecodeError bool
		//createdAccount *GOHMoney.Account
	}

	testSets := []testSet{
		{
			newBalance:     accountBalance{},
			expectedStatus: http.StatusBadRequest,
			expectJsonDecodeError:true,
		},
		{
			newBalance: accountBalance{
				AccountId:-1,
			},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
		},
		{
			newBalance: accountBalance{
				AccountId:1,
			},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
		},
		{
			newBalance: accountBalance{
				AccountId:1,
				Balance:GOHMoney.Balance{},
			},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
		},
		{
			newBalance: accountBalance{
				AccountId:1,
				Balance:GOHMoney.Balance{
					Date: time.Now().AddDate(1000, 1, 1),
				},
			},
			expectedStatus:   http.StatusCreated,
			expectJsonDecodeError:false,
		},
		{
			newBalance: accountBalance{
				AccountId:1,
				Balance:GOHMoney.Balance{
					Date:time.Now().AddDate(1000,1,1),
					Amount:-2000,
				},
			},
			expectedStatus:   http.StatusCreated,
			expectJsonDecodeError:false,
		},
	}

	for _, testSet := range testSets {
		newBalance := testSet.newBalance
		data, err := json.Marshal(newBalance)
		if err != nil {
			t.Fatalf("Failed to form json from balance: %s\nError: %s", newBalance, err.Error())
		}
		req = httptest.NewRequest("POST", endpoint, bytes.NewReader(data))
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		resp = w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Error reading response body. Error: %s", err)
		}
		expectedCode = testSet.expectedStatus
		if resp.StatusCode != expectedCode {
			t.Errorf("Expected response code %d. Got %d\nnewBalance: %s\nRequest body: %s\nReponse body: %s", expectedCode, resp.StatusCode, newBalance, data, body)
		}
		createdBalance := GOHMoneyDB.Balance{}
		err = json.Unmarshal(body, &createdBalance)
		if (err == nil) != (testSet.expectJsonDecodeError == false) {
			t.Errorf("Unexpected error when json decoding response body to balance\nExpect error: %t\nActual  : %s\nnewBalance: %s\nBody: %s", testSet.expectJsonDecodeError, err, testSet.newBalance, body)
		}
		if err == nil && createdBalance.Id == 0 {
			t.Errorf("Unexpected Id. Expected non-zero, got %d", createdBalance.Id)
		}
		if !createdBalance.Date.Equal(newBalance.Date.Truncate(time.Hour * 24)) {
			t.Errorf("Unexpected date.\nExpected: %s\nActual  : %s", newBalance.Date, createdBalance.Date)
		}
		if newBalance.Amount != createdBalance.Amount {
			t.Errorf("Unexpected amount.\nExpected: %f\nActual  : %f", newBalance.Amount, createdBalance.Amount)
		}
	}
}

func Test_BalanceUpdate_ValidBalanceId_InvalidAccount(t *testing.T) {
	router := NewRouter()
	endpoint := func(id uint) string { return fmt.Sprintf(`/balance/%d/update`, id) }
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		t.Fatalf("Unable to prepare DB for testing. Error: %s", err.Error())
		return
	}
	account, accountErr := GOHMoney.NewAccount("TEST_ACCOUNT", time.Now(), pq.NullTime{})
	if accountErr != nil {
		t.Fatalf("Unable to create account object for testing. Error: %s", err.Error())
	}
	createdAccount, err := GOHMoneyDB.CreateAccount(db, account)
	if err != nil {
		t.Fatalf("Unable to create account DB entry for testing. Error: %s", err.Error())
	}
	invalidBalanceId := uint(1)
	type accountBalance struct {
		AccountId uint
		GOHMoney.Balance
	}
	update := accountBalance{ AccountId:createdAccount.Id }
	jsonBytes, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Unable to generate json object for testind. Error: %s", err.Error())
	}
	req := httptest.NewRequest("POST", endpoint(invalidBalanceId), bytes.NewReader(jsonBytes))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()
	expectedCode := http.StatusBadRequest
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d (%s). Got %d (%s)", expectedCode, http.StatusText(expectedCode), resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	expectedBody := GOHMoneyDB.InvalidAccountBalanceError{AccountId:createdAccount.Id,BalanceId:invalidBalanceId}.Error()
	if string(body) != expectedBody {
		t.Errorf("Unexpected response body.\nExpected: %s\nActual  : %s", expectedBody, body)
	}
}

func Test_BalanceUpdate_InvalidUpdateData(t *testing.T) {
	router := NewRouter()
	endpoint := func(id uint) string { return fmt.Sprintf(`/balance/%d/update`, id) }
	validBalanceId := uint(1)
	invalidUpdateData := []byte("INVALID ACCOUNT BALANCE DATA BODY")
	req := httptest.NewRequest("POST", endpoint(validBalanceId), bytes.NewReader(invalidUpdateData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()
	expectedCode := http.StatusBadRequest
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d (%s). Got %d (%s)", expectedCode, http.StatusText(expectedCode), resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	var updatedBalance GOHMoneyDB.Balance
	err = json.Unmarshal(body, &updatedBalance)
	if err == nil {
		t.Error("Expected a json unmarshalling error but nil was returned.")
	}
}

func Test_BalanceUpdate_InvalidUpdateBalance(t *testing.T) {
	router := NewRouter()
	endpoint := func(id uint) string { return fmt.Sprintf(`/balance/%d/update`, id) }
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		t.Fatalf("Unable to prepare DB for testing. Error: %s", err.Error())
		return
	}
	account, accountErr := GOHMoney.NewAccount("TEST_ACCOUNT", time.Now(), pq.NullTime{})
	if accountErr != nil {
		t.Fatalf("Unable to create account object for testing. Error: %s", err.Error())
	}
	createdAccount, err := GOHMoneyDB.CreateAccount(db, account)
	if err != nil {
		t.Fatalf("Unable to create account DB entry for testing. Error: %s", err.Error())
	}
	originalBalance, err := createdAccount.InsertBalance(db, GOHMoney.Balance{Date:time.Now(), Amount:100})
	if err != nil {
		t.Fatalf("Unable to insert balance into DB for testing. Error: %s", err.Error())
	}
	invalidUpdateBalance := GOHMoney.Balance{Date:account.Start.Time.AddDate(-1,0,0), Amount:200}
	type accountBalance struct {
		AccountId uint
		GOHMoney.Balance
	}
	update := accountBalance{
		AccountId:createdAccount.Id,
		Balance:invalidUpdateBalance,
	}
	updateData, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Unable to marshal json for testing. Error: %s", err.Error())
	}
	req := httptest.NewRequest("POST", endpoint(originalBalance.Id), bytes.NewReader(updateData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()
	expectedCode := http.StatusBadRequest
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d (%s). Got %d (%s)", expectedCode, http.StatusText(expectedCode), resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	var updatedBalance GOHMoneyDB.Balance
	err = json.Unmarshal(body, &updatedBalance)
	if err == nil {
		t.Error("Expected a json unmarshalling error but nil was returned.")
	}
}

func Test_BalanceUpdate_Valid(t *testing.T){
	router := NewRouter()
	endpoint := func(id uint) string { return fmt.Sprintf(`/balance/%d/update`, id) }
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		t.Fatalf("Unable to prepare DB for testing. Error: %s", err.Error())
		return
	}
	account, accountErr := GOHMoney.NewAccount("TEST_ACCOUNT", time.Now(), pq.NullTime{})
	if accountErr != nil {
		t.Fatalf("Unable to create account object for testing. Error: %s", err.Error())
	}
	createdAccount, err := GOHMoneyDB.CreateAccount(db, account)
	if err != nil {
		t.Fatalf("Unable to create account DB entry for testing. Error: %s", err.Error())
	}
	originalBalance, err := createdAccount.InsertBalance(db, GOHMoney.Balance{Date:time.Now(), Amount:100})
	if err != nil {
		t.Fatalf("Unable to insert balance into DB for testing. Error: %s", err.Error())
	}
	validUpdateBalance := GOHMoney.Balance{Date: account.Start.Time.AddDate(0, 0, 1), Amount: 200}
	type accountBalance struct {
		AccountId uint
		GOHMoney.Balance
	}
	update := accountBalance{
		AccountId:createdAccount.Id,
		Balance:validUpdateBalance,
	}
	updateData, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Unable to marshal json for testing. Error: %s", err.Error())
	}
	req := httptest.NewRequest("POST", endpoint(originalBalance.Id), bytes.NewReader(updateData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()
	expectedCode := http.StatusNoContent
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d (%s). Got %d (%s)", expectedCode, http.StatusText(expectedCode), resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	var updatedBalance GOHMoneyDB.Balance
	err = json.Unmarshal(body, &updatedBalance)
	if err != nil {
		t.Errorf("Expected no json unmarshalling error received: %s", err.Error())
	}
	if updatedBalance.Id != originalBalance.Id {
		t.Errorf("Balance Id changed during update.\n\tOriginal: %d\n\tFinal   : %d", originalBalance.Id, updatedBalance.Id)
	}
	expectedDate := validUpdateBalance.Date.Truncate(24 * time.Hour)
	if !updatedBalance.Balance.Date.Equal(expectedDate) {
		t.Errorf("Unexpected updated balance Date.\n\tExpected: %s\n\tActual  : %s", validUpdateBalance, updatedBalance.Balance)
	}
	if updatedBalance.Amount != validUpdateBalance.Amount {
		t.Errorf("Unexpected updated balance Amount.\n\tExpected: %s\n\tActual  : %s", validUpdateBalance.Amount, updatedBalance.Amount)
	}
}
