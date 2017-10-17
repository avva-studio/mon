package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/user"
	"strings"
	"testing"
	"time"

	"github.com/GlynOwenHanmer/GOHMoney/account"
	"github.com/GlynOwenHanmer/GOHMoney/balance"
	"github.com/GlynOwenHanmer/GOHMoneyDB"
	gohtime "github.com/GlynOwenHanmer/go-time"
	"github.com/gorilla/mux"
	"github.com/GlynOwenHanmer/GOHMoney/money"
)

func TestMain(m *testing.M) {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Unable to get current user for testing. Error: %s", err.Error())
		return
	}
	if len(usr.HomeDir) < 1 {
		fmt.Printf("Current user has no home directory to load connection string from.")
		return
	}
	connectionString, err = GOHMoneyDB.LoadDBConnectionString(usr.HomeDir + `/.gohmoney/.gohmoneydbtestconnectionstring`)
	if err != nil {
		fmt.Printf("Unable to open DB connection string for testing. Error: %s", err.Error())
		os.Exit(-1)
		return
	}
	os.Exit(m.Run())
}

func Test_Accounts(t *testing.T) {
	req := httptest.NewRequest("GET", "/accounts", nil)
	w := httptest.NewRecorder()
	NewRouter().ServeHTTP(w, req)
	expectedCode := http.StatusOK
	resp := w.Result()
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	accounts := GOHMoneyDB.Accounts{}
	if err := json.Unmarshal(body, &accounts); err != nil {
		t.Errorf("Error unmarshalling response body to Account\nError: %s\nBody: %s", err.Error(), body)
	}
	minAccountsLength := 25
	if len(accounts) < minAccountsLength {
		t.Fatalf("Expected accounts min length %d, got length %d", minAccountsLength, len(accounts))
	}

	a := accounts[0]
	innerAccount, err := account.New("Current", time.Date(2013, 10, 01, 0, 0, 0, 0, time.UTC), gohtime.NullTime{})
	if err != nil {
		t.Fatalf("Error creating a for testing. Error: %s", err.Error())
	}
	expectedAccount := GOHMoneyDB.Account{ID: 1, Account: innerAccount}
	checkAccount(expectedAccount, a, t)
	if t.Failed() {
		t.Logf("Body: %s", body)
	}

	a = accounts[6]
	innerAccount, err = account.New(
		"Patrick",
		time.Date(2015, 9, 14, 0, 0, 0, 0, time.UTC),
		gohtime.NullTime{
			Valid: true,
			Time:  time.Date(2016, 6, 19, 0, 0, 0, 0, time.UTC),
		},
	)
	expectedAccount = GOHMoneyDB.Account{ID: 7, Account: innerAccount}
	checkAccount(expectedAccount, a, t)
}

func checkAccount(expectedAccount GOHMoneyDB.Account, actualAccount GOHMoneyDB.Account, t *testing.T) {
	if actualAccount.Name != expectedAccount.Name {
		t.Errorf("Unexpected Account  name.\nExpected: %s\nActual  : %s", expectedAccount.Name, actualAccount.Name)
	}
	if actualAccount.ID != expectedAccount.ID {
		t.Errorf("Unexpected Account id.\nExpected: %d\nActual  : %d", expectedAccount.ID, actualAccount.ID)
	}
	if !actualAccount.Start().Equal(expectedAccount.Start()) {
		t.Errorf("Unexpected Account Start.\nExpected: %v\nActual  : %v", expectedAccount.Start(), actualAccount.Start())
	}
	if actualAccount.End().Valid != expectedAccount.End().Valid || !actualAccount.End().Time.Equal(expectedAccount.End().Time) {
		t.Errorf("Unexpected Account End.\nExpected: %v\nActual  : %v", expectedAccount.End(), actualAccount.End())
	}
}

func Test_AccountId(t *testing.T) {
	req := httptest.NewRequest("GET", "/account/1", nil)
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
	a := GOHMoneyDB.Account{}
	err = json.Unmarshal(body, &a)
	if err != nil {
		t.Errorf("Error unmarshalling response body to Account\nError: %s\nBody: %s", err.Error(), body)
	}
	expectedAccount := GOHMoneyDB.Account{
		ID: 1,
		Account: account.Account{
			Name: "Current",
		},
	}
	if a.Name != expectedAccount.Name {
		t.Errorf("Unexpected a name.\nExpected: %s\nActual  : %s", expectedAccount.Name, a.Name)
	}
	if a.ID != expectedAccount.ID {
		t.Errorf("Unexpected a id.\nExpected: %d\nActual  : %d", expectedAccount.ID, a.ID)
	}
}

func Test_AccountCreate(t *testing.T) {
	router := NewRouter()
	endpoint := "/account/create"

	invalidData := []byte(`INVALID BODY`)
	req := httptest.NewRequest("POST", endpoint, bytes.NewReader(invalidData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()
	expectedCode := http.StatusBadRequest
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, resp.StatusCode)
	}

	testSets := []struct {
		expectedStatus        int
		name                  string
		start                 time.Time
		end                   gohtime.NullTime
		expectJSONDecodeError bool
		newAccountsCount      int
		createdAccount        *account.Account
	}{
		{
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
			newAccountsCount:      0,
		},
		{
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
			newAccountsCount:      0,
		},
		{
			name:  "TEST_ACCOUNT",
			start: time.Now(),
			end: gohtime.NullTime{
				Valid: true,
				Time:  time.Now().AddDate(0, 0, -1),
			},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
			newAccountsCount:      0,
		},
		{
			name:                  "TEST_ACCOUNT",
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
			newAccountsCount:      0,
		},
		{
			name:                  "TEST_ACCOUNT",
			start:                 time.Now(),
			end:                   gohtime.NullTime{Valid: true},
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
			newAccountsCount:      0,
		},
		{
			name:                  "   ",
			start:                 time.Now(),
			expectedStatus:        http.StatusBadRequest,
			expectJSONDecodeError: true,
			newAccountsCount:      0,
		},
		{
			name:                  "TEST_ACCOUNT",
			start:                 time.Now(),
			expectedStatus:        http.StatusCreated,
			expectJSONDecodeError: false,
			newAccountsCount:      1,
		},
		{
			name:                  "TEST_ACCOUNT",
			start:                 time.Now().AddDate(0, 0, -1),
			end:                   gohtime.NullTime{Valid: true, Time: time.Now()},
			expectedStatus:        http.StatusCreated,
			expectJSONDecodeError: false,
			newAccountsCount:      1,
		},
	}

	for _, testSet := range testSets {
		originalAccounts := getAccounts(router, t)
		newAccount, _ := account.New(testSet.name, testSet.start, testSet.end)
		//if err != nil {
		//	t.Fatalf("Error creating account for testing. Error: %s", err.Error())
		//}
		data, err := json.Marshal(newAccount)
		if err != nil {
			t.Errorf("Failed to form json from account: %s\nError: %s", newAccount, err.Error())
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
			t.Errorf("Expected response code %d. Got %d\nReponse body: %s", expectedCode, resp.StatusCode, body)
		}
		resultantAccounts := getAccounts(router, t)
		newAccountsCount := len(resultantAccounts) - len(originalAccounts)
		if newAccountsCount != testSet.newAccountsCount {
			t.Errorf("Unexpected new accounts count\nExpected: %d\nActual  : %d", testSet.newAccountsCount, newAccountsCount)
		}
		createdAccount := GOHMoneyDB.Account{}
		err = json.Unmarshal(body, &createdAccount)
		if (err == nil) != (testSet.expectJSONDecodeError == false) {
			t.Errorf("Unexpected error when json decoding response body to account\nExpect error: %t\nActual  : %s\nBody: %s", testSet.expectJSONDecodeError, err, body)
		}
		if testSet.newAccountsCount > 0 {
			if createdAccount.ID < 0 {
				t.Errorf("Expected positive createdAccount id but got: %d", createdAccount.ID)
				t.Logf("New account: %s", newAccount)
			}
			if createdAccount.Name != testSet.name {
				t.Errorf("Unexpected created account name\nExpected: %s\nActual  : %s", testSet.name, createdAccount.Name)
				t.Logf("New account: %s", newAccount)
			}
			expectedDateOpened := time.Date(
				testSet.start.Year(),
				testSet.start.Month(),
				testSet.start.Day(),
				0, 0, 0, 0, time.UTC,
			)
			if !expectedDateOpened.Equal(createdAccount.Start()) {
				t.Errorf("Unexpected created account date opened\nExpected: %s\nActual  : %s", expectedDateOpened, createdAccount.Start())
				t.Logf("New account: %s", newAccount)
			}
			expectedDateClosed := gohtime.NullTime{
				Valid: testSet.end.Valid,
				Time: time.Date(
					testSet.end.Time.Year(),
					testSet.end.Time.Month(),
					testSet.end.Time.Day(),
					0, 0, 0, 0, time.UTC,
				),
			}
			if expectedDateClosed.Valid != createdAccount.End().Valid {
				t.Errorf("Unexpected created account date closed validity\nExpected: %t\nActual  : %t", expectedDateClosed.Valid, createdAccount.End().Valid)
			} else if expectedDateClosed.Valid && !expectedDateClosed.Time.Equal(createdAccount.End().Time) {
				t.Errorf("Unexpected created account date closed.\nExpected: %s\nActual  : %s", expectedDateClosed, createdAccount.End())
				t.Logf("New account: %s", newAccount)
			}
		}
	}
}

func getAccounts(router *mux.Router, t *testing.T) account.Accounts {
	req := httptest.NewRequest("GET", "/accounts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	expectedCode := http.StatusOK
	resp := w.Result()
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	accounts := account.Accounts{}
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		t.Errorf("Error unmarshalling response body to Account\nError: %s\nBody: %s", err.Error(), body)
	}
	return accounts
}

func newBalanceIgnoreError(d time.Time, a int64) balance.Balance {
	b, _ := balance.New(d, money.GBP(a))
	return b
}

func Test_AccountBalance_AccountWithBalances(t *testing.T) {
	present := time.Now()
	past := present.AddDate(-1, 0, 0)
	future := present.AddDate(1, 0, 0)
	newAccount, err := account.New("TEST ACCOUNT", past, gohtime.NullTime{})
	if err != nil {
		t.Fatalf("Unable to create new account object for testing. Error: %s", err.Error())
	}
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		t.Fatalf("Unable to prepare DB for testing. Error: %s", err.Error())
		return
	}
	createdAccount, err := GOHMoneyDB.CreateAccount(db, newAccount)
	if err != nil {
		t.Fatalf("Error creating account for test. Error: %s", err.Error())
	}
	pastBalance, err := createdAccount.InsertBalance(db, newBalanceIgnoreError(past, 0))
	if err != nil {
		t.Fatalf("Error adding balance to account for test. Error: %s", err.Error())
	}
	presentBalance, err := createdAccount.InsertBalance(db, newBalanceIgnoreError(present,1))
	if err != nil {
		t.Fatalf("Error adding balance to account for test. Error: %s", err.Error())
	}
	futureBalance, err := createdAccount.InsertBalance(db, newBalanceIgnoreError(future, 2))
	if err != nil {
		t.Fatalf("Error adding balance to account for test. Error: %s", err.Error())
	}

	testSets := []struct {
		paramsString       string
		expectedBalance    GOHMoneyDB.Balance
		expectedStatusCode int
	}{
		{
			paramsString:       `?date=` + urlFormatDateString(past.AddDate(0, 0, -1)),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			paramsString:       ``,
			expectedBalance:    presentBalance,
			expectedStatusCode: http.StatusOK,
		},
		{
			paramsString:       `?date=` + urlFormatDateString(past),
			expectedBalance:    pastBalance,
			expectedStatusCode: http.StatusOK,
		},
		{
			paramsString:       `?date=` + urlFormatDateString(past.AddDate(0, 0, 1)),
			expectedBalance:    pastBalance,
			expectedStatusCode: http.StatusOK,
		},
		{
			paramsString:       `?date=` + urlFormatDateString(present),
			expectedBalance:    presentBalance,
			expectedStatusCode: http.StatusOK,
		},
		{
			paramsString:       `?date=` + urlFormatDateString(present.AddDate(0, 0, 1)),
			expectedBalance:    presentBalance,
			expectedStatusCode: http.StatusOK,
		},
		{
			paramsString:       `?date=` + urlFormatDateString(future),
			expectedBalance:    futureBalance,
			expectedStatusCode: http.StatusOK,
		},
		{
			paramsString:       `?date=` + urlFormatDateString(future.AddDate(0, 0, 1)),
			expectedBalance:    futureBalance,
			expectedStatusCode: http.StatusOK,
		},
		{
			paramsString:       `?date=` + urlFormatDateString(future.AddDate(20, 0, 0)),
			expectedBalance:    futureBalance,
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, testSet := range testSets {
		endpoint := fmt.Sprintf(`/account/%d/balance%s`, createdAccount.ID, testSet.paramsString)
		req := httptest.NewRequest("GET", endpoint, nil)
		w := httptest.NewRecorder()
		NewRouter().ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode == http.StatusNotFound && resp.StatusCode == testSet.expectedStatusCode {
			continue
		}
		if resp.StatusCode != testSet.expectedStatusCode {
			t.Errorf("Unexpected response code. Expected %d, got %d", testSet.expectedStatusCode, resp.StatusCode)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf(`Error reading body from response. Error: %s`, err.Error())
		}
		var balance GOHMoneyDB.Balance
		err = json.Unmarshal(body, &balance)
		if err != nil {
			t.Errorf("Unable to unmarshal json to balance. Error: %s\nBody: %s", err.Error(), body)
			continue
		}
		if equal, _ := balance.Money().Equal(testSet.expectedBalance.Money()); !equal {
			t.Errorf("Unexpected balance amount.\nExpected: %f\nActual  : %f\nParams: %s", testSet.expectedBalance.Money(), balance.Money(), testSet.paramsString)
		}
		if !balance.Date().Equal(testSet.expectedBalance.Date()) {
			t.Errorf("Unexpected Balance Date.\nExpected: %s\nActual  : %s\nParams: %s", testSet.expectedBalance.Date(), balance.Date(), testSet.paramsString)
		}
	}
}

func Test_AccountBalance_InvalidParameter(t *testing.T) {
	testSets := []struct {
		accountID          uint8
		paramsString       string
		expectedAmount     float32
		expectedStatusCode int
	}{
		{
			accountID:          1,
			paramsString:       `?pidgeons=nowthen`,
			expectedAmount: 21.80,
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, testSet := range testSets {
		endpoint := fmt.Sprintf(`/account/%d/balance%s`, testSet.accountID, testSet.paramsString)
		req := httptest.NewRequest("GET", endpoint, nil)
		w := httptest.NewRecorder()
		NewRouter().ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != testSet.expectedStatusCode {
			t.Errorf("Unexpected response code. Expected %d, got %d", testSet.expectedStatusCode, resp.StatusCode)
		}
	}
}

func Test_AccountBalance_AccountWithBalances_SetDate(t *testing.T) {
	accountID := 1
	expectedAmount := int64(2179)
	dateString := `2017-01-18`
	endpoint := fmt.Sprintf(`/account/%d/balance?date=%s`, accountID, dateString)
	req := httptest.NewRequest("GET", endpoint, nil)
	w := httptest.NewRecorder()
	NewRouter().ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected response code. Expected %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var balance GOHMoneyDB.Balance
	err := json.NewDecoder(resp.Body).Decode(&balance)
	if err != nil {
		t.Fatalf("Unable to unmarshal json to balance. Error: %s", err.Error())
	}
	expected := money.GBP(expectedAmount)
	if equal, _ := balance.Money().Equal(expected); !equal {
		t.Errorf("Unexpected balance amount.\nExpected: %+v\nActual  : %+v", expected, balance.Money())
	}
}

func createTestDBAccount(t *testing.T, start time.Time, end gohtime.NullTime) *GOHMoneyDB.Account {
	newAccount, err := account.New("TEST_ACCOUNT", start, end)
	if err != nil {
		t.Fatalf("Error creating account for testing. Error: %s", err.Error())
	}
	db, err := GOHMoneyDB.OpenDBConnection(connectionString)
	if err != nil {
		t.Fatalf("Unable to prepare DB for testing. Error: %s", err.Error())
	}
	defer db.Close()
	createdAccount, err := GOHMoneyDB.CreateAccount(db, newAccount)
	if err != nil {
		t.Fatalf("Error creating account for test. Error: %s", err.Error())
	}
	return createdAccount
}

func Test_AccountBalance_AccountWithoutBalances(t *testing.T) {
	createdAccount := createTestDBAccount(t, time.Now(), gohtime.NullTime{})
	req := httptest.NewRequest("GET", Account(*createdAccount).balanceEndpoint(), nil)
	w := httptest.NewRecorder()
	NewRouter().ServeHTTP(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	expectedStatusCode := http.StatusNotFound
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Unexpected response code. Expected %d, got %d\nBody: %s", expectedStatusCode, resp.StatusCode, body)
	}
}

func TestAccountUpdate_InvalidData(t *testing.T) {
	createdAccount := createTestDBAccount(t, time.Now(), gohtime.NullTime{})
	router := NewRouter()
	invalidUpdateData := []byte("INVALID ACCOUNT BALANCE DATA BODY")
	req := httptest.NewRequest("PUT", Account(*createdAccount).updateEndpoint(), bytes.NewReader(invalidUpdateData))
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

func ifErrorFatal(t *testing.T, err error, message string) {
	if err != nil {
		message = strings.TrimSpace(message)
		if len(message) > 0 {
			message = fmt.Sprintf("%s: ", message)
		}
		t.Fatalf("%s%s", message, err)
	}
}

func TestAccountUpdate_ValidData(t *testing.T) {
	original := createTestDBAccount(t, time.Now(), gohtime.NullTime{})
	router := NewRouter()
	updates, err := account.New(
		"UPDATED ACCOUNT NAME",
		time.Now().Add(24*time.Hour).Truncate(24*time.Hour),
		gohtime.NullTime{
			Valid: true,
			Time:  time.Now().Add(72 * time.Hour).Truncate(24 * time.Hour),
		},
	)
	ifErrorFatal(t, err, "Error creating new account object for testing")
	updateBytes, err := json.Marshal(updates)
	ifErrorFatal(t, err, "Error marshaling json for testing")
	req := httptest.NewRequest("PUT", Account(*original).updateEndpoint(), bytes.NewReader(updateBytes))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()
	expectedCode := http.StatusOK
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected response code %d (%s). Got %d (%s)", expectedCode, http.StatusText(expectedCode), resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	var updated GOHMoneyDB.Account
	err = json.Unmarshal(body, &updated)
	if err != nil {
		t.Errorf("Error unmarshalling updated account body.\nError: %s\nBody: %s", err, body)
	}
	if !updated.Account.Equal(updates) {
		t.Errorf("Returned account does not represent updates applied.\n\tReturned: %s\n\tApplied: %s", updated.Account, updates)
	}
}

func TestAccountDelete(t *testing.T) {
	original := createTestDBAccount(t, time.Now(), gohtime.NullTime{})
	router := NewRouter()
	req := httptest.NewRequest("DELETE", Account(*original).deleteEndpoint(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body. Error: %s", err)
	}
	if expected := http.StatusNoContent; resp.StatusCode != expected {
		t.Errorf("Expected status code %d (%s) but got %d (%s).\nBody: %s", expected, http.StatusText(expected), resp.StatusCode, http.StatusText(resp.StatusCode), body)
	}
}
