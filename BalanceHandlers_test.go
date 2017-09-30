package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GlynOwenHanmer/GOHMoney/account"
	"github.com/GlynOwenHanmer/GOHMoney/balance"
	"github.com/GlynOwenHanmer/GOHMoneyDB"
	gohtime "github.com/GlynOwenHanmer/go-time"
)

func prepareTestDB(t *testing.T) *sql.DB {
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		t.Fatalf("Unable to prepare DB for testing. Error: %s", err.Error())
		return nil
	}
	return db
}

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
	expectedID := uint(1)
	actualID := balances[0].ID
	if expectedID != actualID {
		t.Errorf(`Unexpected Id.\nExpected: %d\nActual  : %d`, expectedID, actualID)
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

	now := time.Now()
	a, err := account.New("TEST_ACCOUNT", now, gohtime.NullTime{})
	db := prepareTestDB(t)
	account, err := GOHMoneyDB.CreateAccount(db, a)
	if err != nil {
		t.Fatalf("Error creating new Account for testing. Error: %s", err.Error())
	}

	type accountBalance struct {
		AccountID       int `json:"account_id"`
		balance.Balance `json:"balance"`
	}

	testSets := []struct {
		newBalance            accountBalance
		expectedStatus        int
		expectJSONDecodeError bool
	}{
		{
			newBalance:            accountBalance{},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalance{
				AccountID: -1,
			},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalance{
				AccountID: int(account.ID),
			},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalance{
				AccountID: int(account.ID),
				Balance:   balance.Balance{},
			},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalance{
				AccountID: int(account.ID),
				Balance: balance.Balance{
					Date: now.AddDate(-1, 0, 0),
				},
			},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalance{
				AccountID: int(account.ID),
				Balance: balance.Balance{
					Date: now,
				},
			},
			expectedStatus:        http.StatusCreated,
			expectJSONDecodeError: false,
		},
		{
			newBalance: accountBalance{
				AccountID: int(account.ID),
				Balance: balance.Balance{
					Date: now.AddDate(1000, 1, 1),
				},
			},
			expectedStatus:        http.StatusCreated,
			expectJSONDecodeError: false,
		},
		{
			newBalance: accountBalance{
				AccountID: int(account.ID),
				Balance: balance.Balance{
					Date:   time.Now().AddDate(1000, 1, 1),
					Amount: -2000,
				},
			},
			expectedStatus:        http.StatusCreated,
			expectJSONDecodeError: false,
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
		if (err == nil) != (testSet.expectJSONDecodeError == false) {
			t.Errorf("Unexpected error when json decoding response body to balance\nExpect error: %t\nActual  : %s\nnewBalance: %s\nBody: %s", testSet.expectJSONDecodeError, err, testSet.newBalance, body)
		}
		if err != nil {
			continue
		}
		if createdBalance.ID == 0 {
			t.Errorf("Unexpected Id. Expected non-zero, got %d", createdBalance.ID)
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
	db := prepareTestDB(t)
	account, err := account.New("TEST_ACCOUNT", time.Now(), gohtime.NullTime{})
	if err != nil {
		t.Fatalf("Unable to create account object for testing. Error: %s", err.Error())
	}
	createdAccount, err := GOHMoneyDB.CreateAccount(db, account)
	if err != nil {
		t.Fatalf("Unable to create account DB entry for testing. Error: %s", err.Error())
	}
	invalidBalanceID := uint(1)
	type accountBalance struct {
		AccountID uint
		balance.Balance
	}
	update := accountBalance{AccountID: createdAccount.ID}
	jsonBytes, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Unable to generate json object for testind. Error: %s", err.Error())
	}
	req := httptest.NewRequest("POST", endpoint(invalidBalanceID), bytes.NewReader(jsonBytes))
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
	expectedBody := GOHMoneyDB.InvalidAccountBalanceError{AccountID: createdAccount.ID, BalanceID: invalidBalanceID}.Error()
	if string(body) != expectedBody {
		t.Errorf("Unexpected response body.\nExpected: %s\nActual  : %s", expectedBody, body)
	}
}

func Test_BalanceUpdate_InvalidUpdateData(t *testing.T) {
	router := NewRouter()
	endpoint := func(id uint) string { return fmt.Sprintf(`/balance/%d/update`, id) }
	validBalanceID := uint(1)
	invalidUpdateData := []byte("INVALID ACCOUNT BALANCE DATA BODY")
	req := httptest.NewRequest("POST", endpoint(validBalanceID), bytes.NewReader(invalidUpdateData))
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
	account := createTestDBAccount(t, time.Now(), gohtime.NullTime{})
	db := prepareTestDB(t)
	defer db.Close()
	originalBalance, err := account.InsertBalance(db, balance.Balance{Date: time.Now(), Amount: 100})
	if err != nil {
		t.Fatalf("Unable to insert balance into DB for testing. Error: %s", err.Error())
	}
	invalidUpdateBalance := balance.Balance{Date: account.Start().AddDate(-1, 0, 0), Amount: 200}
	type accountBalance struct {
		AccountID uint
		balance.Balance
	}
	update := accountBalance{
		AccountID: account.ID,
		Balance:   invalidUpdateBalance,
	}
	updateData, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Unable to marshal json for testing. Error: %s", err.Error())
	}
	req := httptest.NewRequest("POST", endpoint(originalBalance.ID), bytes.NewReader(updateData))
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

func Test_BalanceUpdate_Valid(t *testing.T) {
	router := NewRouter()
	a := createTestDBAccount(t, time.Now(), gohtime.NullTime{})
	db := prepareTestDB(t)
	originalBalance, err := a.InsertBalance(db, balance.Balance{Date: time.Now(), Amount: 100})
	if err != nil {
		t.Fatalf("Unable to insert balance into DB for testing. Error: %s", err.Error())
	}
	validUpdateBalance := balance.Balance{Date: a.Start().AddDate(0, 0, 1), Amount: 200}
	type accountBalance struct {
		AccountID uint
		balance.Balance
	}
	update := accountBalance{
		AccountID: a.ID,
		Balance:   validUpdateBalance,
	}
	updateData, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Unable to marshal json for testing. Error: %s", err.Error())
	}
	req := httptest.NewRequest("POST", Balance(originalBalance).balanceUpdateEndpoint(), bytes.NewReader(updateData))
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
	if updatedBalance.ID != originalBalance.ID {
		t.Errorf("Balance Id changed during update.\n\tOriginal: %d\n\tFinal   : %d", originalBalance.ID, updatedBalance.ID)
	}
	expectedDate := validUpdateBalance.Date.Truncate(24 * time.Hour)
	if !updatedBalance.Balance.Date.Equal(expectedDate) {
		t.Errorf("Unexpected updated balance Date.\n\tExpected: %s\n\tActual  : %s", validUpdateBalance, updatedBalance.Balance)
	}
	if updatedBalance.Amount != validUpdateBalance.Amount {
		t.Errorf("Unexpected updated balance Amount.\n\tExpected: %s\n\tActual  : %s", validUpdateBalance.Amount, updatedBalance.Amount)
	}
}
