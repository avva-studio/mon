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
	"github.com/GlynOwenHanmer/GOHMoney/money"
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
	actualAmount := balances[0].Money()
	expectedMoney := int64(63641)
	if equal, _ := actualAmount.Equal(money.GBP(expectedMoney)); !equal {
		t.Errorf("first balance, expected balance amount of %v but got %v", expectedMoney, actualAmount)
	}
	expectedDate, err := parseDateString("2016-06-17")
	if err != nil {
		t.Fatalf("Error parsing date string for use in tests. Error: %s", err.Error())
	}
	actualDate := balances[0].Date()
	if !expectedDate.Equal(actualDate) {
		t.Errorf("first balance, expected date of %s but got %s", urlFormatDateString(expectedDate), urlFormatDateString(actualDate))
	}
}

// accountBalanceTestJSONHelper is an internal type used to marshal and unmarshal json for methods.
type accountBalanceTestJSONHelper struct {
	AccountID int
	Balance balance.Balance
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

	testSets := []struct {
		newBalance            accountBalanceTestJSONHelper
		expectedStatus        int
		expectJSONDecodeError bool
	}{
		{
			newBalance:            accountBalanceTestJSONHelper{},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalanceTestJSONHelper{
				AccountID: -1,
			},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalanceTestJSONHelper{
				AccountID: int(account.ID),
			},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalanceTestJSONHelper{
				AccountID: int(account.ID),
				Balance:   balance.Balance{},
			},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalanceTestJSONHelper{
				AccountID: int(account.ID),
				Balance: newBalanceIgnoreError(now.AddDate(-1, 0, 0), 0),
			},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
		},
		{
			newBalance: accountBalanceTestJSONHelper{
				AccountID: int(account.ID),
				Balance: newBalanceIgnoreError(now, 0),
			},
			expectedStatus:        http.StatusCreated,
			expectJSONDecodeError: false,
		},
		{
			newBalance: accountBalanceTestJSONHelper{
				AccountID: int(account.ID),
				Balance: newBalanceIgnoreError(now.AddDate(1000, 1, 1), 0),
			},
			expectedStatus:        http.StatusCreated,
			expectJSONDecodeError: false,
		},
		{
			newBalance: accountBalanceTestJSONHelper{
				AccountID: int(account.ID),
				Balance: newBalanceIgnoreError(now.AddDate(1000, 1, 1), -2000),
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
		if !createdBalance.Date().Equal(newBalance.Balance.Date().Truncate(time.Hour * 24)) {
			t.Errorf("Unexpected date.\nExpected: %s\nActual  : %s", newBalance.Balance.Date(), createdBalance.Date())
		}
		if equal, _ := newBalance.Balance.Money().Equal(createdBalance.Money()); !equal {
			t.Errorf("Unexpected amount.\nExpected: %f\nActual  : %f", newBalance.AccountID, createdBalance.Money())
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

	update := accountBalanceTestJSONHelper{AccountID: int(createdAccount.ID), Balance:newBalanceIgnoreError(time.Now(), 123)}
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
		t.Errorf("Unexpected response body.\nExpected: %s\nActual  : %s.\nRequst body: %s", expectedBody, body, jsonBytes)
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
	originalBalance, err := account.InsertBalance(db, newBalanceIgnoreError(time.Now(), 100))
	if err != nil {
		t.Fatalf("Unable to insert balance into DB for testing. Error: %s", err.Error())
	}
	invalidUpdateBalance := newBalanceIgnoreError(account.Start().AddDate(-1, 0, 0), 200)
	update := accountBalanceTestJSONHelper{
		AccountID: int(account.ID),
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
	originalBalance, err := a.InsertBalance(db, newBalanceIgnoreError(time.Now(), 100))
	if err != nil {
		t.Fatalf("Unable to insert balance into DB for testing. Error: %s", err.Error())
	}
	validUpdateBalance := newBalanceIgnoreError(a.Start().AddDate(0, 0, 1), 200)
	update := accountBalanceTestJSONHelper{
		AccountID: int(a.ID),
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
	expectedDate := validUpdateBalance.Date().Truncate(24 * time.Hour)
	if !updatedBalance.Balance.Date().Equal(expectedDate) {
		t.Errorf("Unexpected updated balance Date.\n\tExpected: %s\n\tActual  : %s", validUpdateBalance, updatedBalance.Balance)
	}
	if equal, _ := updatedBalance.Money().Equal(validUpdateBalance.Money()); !equal {
		t.Errorf("Unexpected updated balance Amount.\n\tExpected: %s\n\tActual  : %s", validUpdateBalance.Money(), updatedBalance.Money())
	}
}
