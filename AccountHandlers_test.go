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
	connectionString, err = GOHMoneyDB.LoadDBConnectionString(usr.HomeDir + `/.gohmoneydbtestconnectionstring`)
	if err != nil {
		fmt.Printf("Unable to open DB connection string for testing. Error: %s", err.Error())
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
		t.Fatalf("Expected accounts min length %d, got length %d", minAccountsLength, len(accounts))
	}

	account := accounts[0]
	innerAccount, err := GOHMoney.NewAccount("Current", time.Date(2013,10,01,0,0,0,0,time.UTC), pq.NullTime{})
	if err != nil {
		t.Fatalf("Error creating account for testing. Error: %s", err.Error())
	}
	expectedAccount := GOHMoneyDB.Account{ Id:1, Account: innerAccount }
	checkAccount(expectedAccount, account, t)
	if t.Failed() {
		t.Logf("Body: %s", body)
	}

	account = accounts[6]
	innerAccount, err = GOHMoney.NewAccount(
		"Patrick",
		time.Date(2015,9,14,0,0,0,0,time.UTC),
		pq.NullTime{
			Valid:true,
			Time:time.Date(2016,6,19,0,0,0,0,time.UTC),
		},
	)
	expectedAccount = GOHMoneyDB.Account{ Id:7,Account:innerAccount }
	checkAccount(expectedAccount, account, t)
}

func checkAccount(expectedAccount GOHMoneyDB.Account, actualAccount GOHMoneyDB.Account, t *testing.T) {
	if actualAccount.Name != expectedAccount.Name {
		t.Errorf("Unexpected Account  name.\nExpected: %s\nActual  : %s", expectedAccount.Name, actualAccount.Name)
	}
	if actualAccount.Id != expectedAccount.Id {
		t.Errorf("Unexpected Account id.\nExpected: %d\nActual  : %d", expectedAccount.Id, actualAccount.Id)
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
		expectedStatus int
		name string
		start time.Time
		end pq.NullTime
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
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			name:"TEST_ACCOUNT",
			start:time.Now(),
			end:pq.NullTime{
				Valid:true,
				Time:time.Now().AddDate(0,0,-1),
			},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			name: "TEST_ACCOUNT",
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			name:"TEST_ACCOUNT",
			start:time.Now(),
			end:pq.NullTime{Valid:true},
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			name:"   ",
			start:time.Now(),
			expectedStatus:http.StatusBadRequest,
			expectJsonDecodeError:true,
			newAccountsCount:0,
		},
		{
			name:"TEST_ACCOUNT",
			start:time.Now(),
			expectedStatus:http.StatusCreated,
			expectJsonDecodeError:false,
			newAccountsCount:1,
		},
		{
			name:"TEST_ACCOUNT",
			start:time.Now().AddDate(0, 0, -1),
			end:pq.NullTime{Valid:true,Time:time.Now()},
			expectedStatus:http.StatusCreated,
			expectJsonDecodeError:false,
			newAccountsCount:1,
		},
	}

	for _, testSet := range testSets {
		originalAccounts := getAccounts(router, t)
		newAccount, _ := GOHMoney.NewAccount(testSet.name, testSet.start, testSet.end)
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
		if (err == nil) != (testSet.expectJsonDecodeError == false) {
			t.Errorf("Unexpected error when json decoding response body to account\nExpect error: %t\nActual  : %s\nBody: %s", testSet.expectJsonDecodeError, err, body)
		}
		if testSet.newAccountsCount > 0 {
			if createdAccount.Id < 0 {
				t.Errorf("Expected positive createdAccount id but got: %d", createdAccount.Id)
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
				0,0,0,0,time.UTC,
			)
			if !expectedDateOpened.Equal(createdAccount.Start()) {
				t.Errorf("Unexpected created account date opened\nExpected: %s\nActual  : %s", expectedDateOpened, createdAccount.Start())
				t.Logf("New account: %s", newAccount)
			}
			expectedDateClosed := pq.NullTime{
				Valid:testSet.end.Valid,
				Time:time.Date(
					testSet.end.Time.Year(),
					testSet.end.Time.Month(),
					testSet.end.Time.Day(),
					0,0,0,0,time.UTC,
				),
			}
			if expectedDateClosed.Valid != createdAccount.End().Valid {
				t.Errorf("Unexpected created account date closed validity\nExpected: %t\nActual  : %t", expectedDateClosed.Valid, createdAccount.End().Valid)
			} else if expectedDateClosed.Valid && !expectedDateClosed.Time.Equal(createdAccount.End().Time){
				t.Errorf("Unexpected created account date closed.\nExpected: %s\nActual  : %s", expectedDateClosed, createdAccount.End())
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

func Test_AccountBalance_AccountWithBalances(t *testing.T) {
	present := time.Now()
	past := present.AddDate(-1, 0, 0)
	future := present.AddDate(1, 0, 0)
	newAccount, err := GOHMoney.NewAccount("TEST ACCOUNT", past, pq.NullTime{})
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
	pastBalance, err := createdAccount.InsertBalance(db,GOHMoney.Balance{Date:past,Amount:0})
	if err != nil {
		t.Fatalf("Error adding balance to account for test. Error: %s", err.Error())
	}
	presentBalance, err := createdAccount.InsertBalance(db,GOHMoney.Balance{Date:present,Amount:1})
	if err != nil {
		t.Fatalf("Error adding balance to account for test. Error: %s", err.Error())
	}
	futureBalance, err := createdAccount.InsertBalance(db,GOHMoney.Balance{Date:future,Amount:2})
	if err != nil {
		t.Fatalf("Error adding balance to account for test. Error: %s", err.Error())
	}

	testSets := []struct{
		paramsString string
		expectedBalance GOHMoneyDB.Balance
		expectedStatusCode int
}{
		{
			paramsString: `?date=`+ urlFormatDateString(past.AddDate(0,0,-1)),
			expectedStatusCode:http.StatusNotFound,
		},
		{
			paramsString:``,
			expectedBalance:presentBalance,
			expectedStatusCode:http.StatusOK,
		},
		{
			paramsString: `?date=`+ urlFormatDateString(past),
			expectedBalance:pastBalance,
			expectedStatusCode:http.StatusOK,
		},
		{
			paramsString: `?date=`+ urlFormatDateString(past.AddDate(0,0,1)),
			expectedBalance:pastBalance,
			expectedStatusCode:http.StatusOK,
		},
		{
			paramsString: `?date=`+ urlFormatDateString(present),
			expectedBalance:presentBalance,
			expectedStatusCode:http.StatusOK,
		},
		{
			paramsString: `?date=`+ urlFormatDateString(present.AddDate(0,0,1)),
			expectedBalance:presentBalance,
			expectedStatusCode:http.StatusOK,
		},
		{
			paramsString: `?date=`+ urlFormatDateString(future),
			expectedBalance:futureBalance,
			expectedStatusCode:http.StatusOK,
		},
		{
			paramsString: `?date=`+ urlFormatDateString(future.AddDate(0,0,1)),
			expectedBalance:futureBalance,
			expectedStatusCode:http.StatusOK,
		},
		{
			paramsString: `?date=`+ urlFormatDateString(future.AddDate(20,0,0)),
			expectedBalance:futureBalance,
			expectedStatusCode:http.StatusOK,
		},
	}
	for _, testSet := range testSets {
		endpoint := fmt.Sprintf(`/account/%d/balance%s`, createdAccount.Id, testSet.paramsString)
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
		if balance.Amount != testSet.expectedBalance.Amount {
			t.Errorf("Unexpected balance amount.\nExpected: %f\nActual  : %f\nParams: %s", testSet.expectedBalance.Amount, balance.Amount, testSet.paramsString)
		}
		if !balance.Date.Equal(testSet.expectedBalance.Date ) {
			t.Errorf("Unexpected Balance Date.\nExpected: %s\nActual  : %s\nParams: %s", testSet.expectedBalance.Date, balance.Date, testSet.paramsString)
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
	newAccount, err := GOHMoney.NewAccount(
		"TEST_ACCOUNT",
		time.Now().AddDate(0, 0, -1),
		pq.NullTime{
			Valid:true,
			Time:time.Now(),
		},
	)
	if err != nil {
		t.Fatalf("Error creating account for testing. Error: %s", err.Error())
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