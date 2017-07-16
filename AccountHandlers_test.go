package main

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"encoding/json"
	"github.com/GlynOwenHanmer/GOHMoney"
	"github.com/GlynOwenHanmer/GOHMoneyDB"
	"net/http"
	"time"
	"github.com/lib/pq"
	"bytes"
	"github.com/gorilla/mux"
	"fmt"
	"os"
	"os/user"
)

func TestMain(m *testing.M){
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Unable to get current user for testing. Error: %s", err.Error())
		return
	}
	if len(usr.HomeDir) < 1 {
		fmt.Printf("Current user has no home directory to load connection string from.")
		return
	}
	connectionString, err = loadDBConnectionString(usr.HomeDir + `/.gohmoneydbtestconnectionstring`)
	if err != nil {
		fmt.Printf("Unable to open DB connection for testing. Error: %s", err.Error())
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
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		t.Errorf("Error unmarshalling response body to Account\nError: %s\nBody: %s", err.Error(), body)
	}
	minAccountsLength := 25
	if len(accounts) < minAccountsLength {
		t.Errorf("Expected accounts min length %d, got length %d", minAccountsLength, len(accounts))
	}
	account := accounts[0]
	expectedAccount := GOHMoneyDB.Account{
		Id:1,
		Account: GOHMoney.Account{
			Name:"Current",
			DateOpened:time.Date(2013,10,01,0,0,0,0,time.UTC),
		},
	}

	checkAccount(expectedAccount, account, t)
	account = accounts[6]
	expectedAccount = GOHMoneyDB.Account{
		Id:7,
		Account:GOHMoney.Account{
			Name:"Patrick",
			DateOpened:time.Date(2015,9,14,0,0,0,0,time.UTC),
			DateClosed:pq.NullTime{Valid:true, Time:time.Date(2016,6,19,0,0,0,0,time.UTC)},
		},
	}
	checkAccount(expectedAccount, account, t)
}

func checkAccount(expectedAccount GOHMoneyDB.Account, actualAccount GOHMoneyDB.Account, t *testing.T) {
	if actualAccount.Name != expectedAccount.Name {
		t.Errorf("Unexpected account name.\nExpected: %s\nActual  : %s", expectedAccount.Name, actualAccount.Name)
	}
	if actualAccount.Id != expectedAccount.Id {
		t.Errorf("Unexpected account id.\nExpected: %d\nActual  : %d", expectedAccount.Id, actualAccount.Id)
	}
	if !actualAccount.DateOpened.Equal(expectedAccount.DateOpened) {
		t.Errorf("Unexpected DateOpened.\nExpected: %v\nActual  : %v", expectedAccount.DateOpened.Format(GOHMoneyDB.DbDateFormat), actualAccount.DateOpened.Format(GOHMoneyDB.DbDateFormat))
	}
	if actualAccount.DateClosed.Valid != expectedAccount.DateClosed.Valid {
		t.Errorf("Unexpected DateClosed validity.\nExpected: %t\nActual  : %t", expectedAccount.DateClosed.Valid, actualAccount.DateClosed.Valid)
	}
	if actualAccount.DateClosed.Valid && expectedAccount.DateClosed.Valid && !actualAccount.DateClosed.Time.Equal(expectedAccount.DateClosed.Time) {
		t.Errorf("Unexpected DateClosed.\nExpected: %v\nActual  : %v", expectedAccount.DateClosed, actualAccount.DateClosed)
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
	account := GOHMoneyDB.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		t.Errorf("Error unmarshalling response body to Account\nError: %s\nBody: %s", err.Error(), body)
	}
	expectedAccount := GOHMoneyDB.Account{
		Id: 1,
		Account: GOHMoney.Account{
			Name: "Current",
		},
	}
	if account.Name != expectedAccount.Name {
		t.Errorf("Unexpected account name.\nExpected: %s\nActual  : %s", expectedAccount.Name, account.Name)
	}
	if account.Id != expectedAccount.Id {
		t.Errorf("Unexpected account id.\nExpected: %d\nActual  : %d", expectedAccount.Id, account.Id)
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

	testSets := []struct{
		newAccount GOHMoney.Account
		expectedStatus int
		expectJsonDecodeError bool
		newAccountsCount int
		createdAccount *GOHMoney.Account
	}{
		{
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			newAccount:GOHMoney.Account{},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			newAccount:GOHMoney.Account{
				Name:"TEST_ACCOUNT",
				DateOpened:time.Now(),
				DateClosed:pq.NullTime{Valid:true,Time:time.Now().AddDate(0,0,-1)},
			},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			newAccount:GOHMoney.Account{
				Name: "TEST_ACCOUNT",
			},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			newAccount:GOHMoney.Account{
				Name:"TEST_ACCOUNT",
				DateOpened:time.Now(),
				DateClosed:pq.NullTime{Valid:true},
			},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			newAccount:GOHMoney.Account{
				Name:"   ",
				DateOpened:time.Now(),
			},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			newAccount:GOHMoney.Account{
				Name:"TEST_ACCOUNT",
				DateOpened:time.Now(),
			},
			expectedStatus:http.StatusCreated,
			expectJsonDecodeError:false,
			newAccountsCount:1,
		},
		{
			newAccount:GOHMoney.Account{
				Name:"TEST_ACCOUNT",
				DateOpened:time.Now().AddDate(0, 0, -1),
				DateClosed:pq.NullTime{Valid:true,Time:time.Now()},
			},
			expectedStatus:http.StatusCreated,
			expectJsonDecodeError:false,
			newAccountsCount:1,
		},
	}

	for _, testSet := range testSets {
		originalAccounts := getAccounts(router, t)
		newAccount := testSet.newAccount
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
		if (err == nil) != (testSet.expectJsonDecodeError == false) {
			t.Errorf("Unexpected error when json decoding response body to account\nExpect error: %t\nActual  : %s\nBody: %s", testSet.expectJsonDecodeError, err, body)
		}
		if testSet.newAccountsCount > 0 {
			if createdAccount.Id < 0 {
				t.Errorf("Expected positive createdAccount id but got: %d", createdAccount.Id)
				t.Logf("New account: %s", newAccount)
			}
			if createdAccount.Name != testSet.newAccount.Name {
				t.Errorf("Unexpected created account name\nExpected: %s\nActual  : %s", testSet.newAccount.Name, createdAccount.Name)
				t.Logf("New account: %s", newAccount)
			}
			expectedDateOpened := time.Date(
				testSet.newAccount.DateOpened.Year(),
				testSet.newAccount.DateOpened.Month(),
				testSet.newAccount.DateOpened.Day(),
				0,0,0,0,time.UTC,
			)
			if !expectedDateOpened.Equal(createdAccount.DateOpened) {
				t.Errorf("Unexpected created account date opened\nExpected: %s\nActual  : %s", expectedDateOpened, createdAccount.DateOpened)
				t.Logf("New account: %s", newAccount)
			}
			expectedDateClosed := pq.NullTime{
				Valid:testSet.newAccount.DateClosed.Valid,
				Time:time.Date(
					testSet.newAccount.DateClosed.Time.Year(),
					testSet.newAccount.DateClosed.Time.Month(),
					testSet.newAccount.DateClosed.Time.Day(),
					0,0,0,0,time.UTC,
				),
			}
			if expectedDateClosed.Valid != createdAccount.DateClosed.Valid {
				t.Errorf("Unexpected created account date closed validity\nExpected: %t\nActual  : %t", expectedDateClosed.Valid, createdAccount.DateClosed.Valid)
			} else if expectedDateClosed.Valid && !expectedDateClosed.Time.Equal(createdAccount.DateClosed.Time){
				t.Errorf("Unexpected created account date closed.\nExpected: %s\nActual  : %s", expectedDateClosed, createdAccount.DateClosed)
				t.Logf("New account: %s", newAccount)
			}
		}
	}
}

func getAccounts(router *mux.Router, t *testing.T) GOHMoney.Accounts {
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
	accounts := GOHMoney.Accounts{}
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		t.Errorf("Error unmarshalling response body to Account\nError: %s\nBody: %s", err.Error(), body)
	}
	return accounts
}

func Test_AccountBalance_AccountWithBalances_DefaultDate(t *testing.T) {
	testSets := []struct{
		accountId uint8
		paramsString string
		expectedAmount float32
		expectedStatusCode int
}{
		{
			accountId:1,
			paramsString:``,
			expectedAmount:1476.680054,
			expectedStatusCode:http.StatusOK,
		},
		{
			accountId:1,
			paramsString: `?date=2017-01-18`,
			expectedAmount: 21.80,
			expectedStatusCode:http.StatusOK,
		},
	}
	for _, testSet := range testSets {
		endpoint := fmt.Sprintf(`/account/%d/balance%s`, testSet.accountId, testSet.paramsString)
		req := httptest.NewRequest("GET", endpoint, nil)
		w := httptest.NewRecorder()
		NewRouter().ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != testSet.expectedStatusCode {
			t.Errorf("Unexpected response code. Expected %d, got %d", testSet.expectedStatusCode, resp.StatusCode)
		}
		var balance GOHMoneyDB.Balance
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf(`Error reading body from response. Error: %s`, err.Error())
		}
		err = json.Unmarshal(body, &balance)
		if err != nil {
			t.Errorf("Unable to unmarshal json to balance. Error: %s\nBody: %s", err.Error(), body)
			continue
		}
		if balance.Amount != testSet.expectedAmount {
			t.Errorf("Unexpected balance amount.\nExpected: %f\nActual  : %f\nParams: %s", testSet.expectedAmount, balance.Amount, testSet.paramsString)
		}
	}
}

func Test_AccountBalance_InvalidParameter(t *testing.T) {
	testSets := []struct{
		accountId uint8
		paramsString string
		expectedAmount float32
		expectedStatusCode int
	}{
		{
			accountId:1,
			paramsString: `?pidgeons=nowthen`,
			expectedAmount: 21.80,
			expectedStatusCode:http.StatusBadRequest,
		},
	}
	for _, testSet := range testSets {
		endpoint := fmt.Sprintf(`/account/%d/balance%s`, testSet.accountId, testSet.paramsString)
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
	accountId := 1
	expectedAmount := float32(21.80)
	dateString := `2017-01-18`
	endpoint := fmt.Sprintf(`/account/%d/balance?date=%s`, accountId, dateString)
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
	if balance.Amount != expectedAmount {
		t.Errorf("Unexpected balance amount.\nExpected: %f\nActual  : %f", expectedAmount, balance.Amount)
	}
}

func Test_AccountBalance_AccountWithoutBalances(t *testing.T) {
	newAccount := GOHMoney.Account{
		Name:"TEST_ACCOUNT",
		DateOpened:time.Now().AddDate(0, 0, -1),
		DateClosed:pq.NullTime{Valid:true,Time:time.Now()},
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
	req := httptest.NewRequest("GET", Account(createdAccount).balanceEndpoint(), nil)
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