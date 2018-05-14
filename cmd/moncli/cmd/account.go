package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	gtime "github.com/glynternet/go-time"
	"github.com/glynternet/mon/client"
	"github.com/glynternet/mon/pkg/date"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyDate     = "date"
	keyAmount   = "amount"
	keyName     = "name"
	keyCurrency = "currency"
	keyOpened   = "opened"
	keyClosed   = "closed"
	keyLimit    = "limit"
)

var (
	accountOpened = date.Flag()
	accountClosed = date.Flag()
)

var accountCmd = &cobra.Command{
	Use: "account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument for account ID, received %d", len(args))
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parsing account id")
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			return errors.Wrap(err, "selecting account")
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)
		return nil
	},
}

var accountAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add an account",
	RunE: func(cmd *cobra.Command, args []string) error {
		cc, err := currency.NewCode(viper.GetString(keyCurrency))
		if err != nil {
			return errors.Wrap(err, "creating new currency code")
		}

		opened := gtime.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
		if accountOpened.Time != nil {
			opened.Time = *accountOpened.Time
		}

		var closed gtime.NullTime
		if accountClosed.Time != nil {
			closed = gtime.NullTime{
				Valid: true,
				Time:  *accountClosed.Time,
			}
		}

		var ops []account.Option
		if closed.Valid {
			ops = append(ops, account.CloseTime(closed.Time))
		}
		a, err := account.New(
			viper.GetString(keyName),
			*cc,
			opened.Time,
			ops...,
		)
		if err != nil {
			return errors.Wrap(err, "creating new account for insert")
		}

		i, err := client.Client(viper.GetString(keyServerHost)).InsertAccount(*a)
		if err != nil {
			return errors.Wrap(err, "inserting new account")
		}
		table.Accounts(storage.Accounts{*i}, os.Stdout)
		return nil
	},
}

var accountUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update an account",
	Long:  "update an account with the given details. All of the details of an account must be provided, even if they are exactly the same as the original account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument for account ID, received %d", len(args))
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parsing account id")
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			return errors.Wrap(err, "selecting account to update")
		}

		opened := time.Now()
		if accountOpened.Time != nil {
			opened = *accountOpened.Time
		}

		var ops []account.Option
		if accountClosed.Time != nil {
			ops = append(ops, account.CloseTime(*accountClosed.Time))
		}

		cc, err := currency.NewCode(viper.GetString(keyCurrency))
		if err != nil {
			return errors.Wrap(err, "creating new currency code")
		}

		us, err := account.New(viper.GetString(keyName), *cc, opened, ops...)
		if err != nil {
			return errors.Wrap(err, "creating account for update")
		}

		updated, err := c.UpdateAccount(a, us)
		if err != nil {
			return errors.Wrap(err, "updating account")
		}

		fmt.Println("ORIGINAL")
		table.Accounts(storage.Accounts{*a}, os.Stdout)

		fmt.Println("UPDATED")
		table.Accounts(storage.Accounts{*updated}, os.Stdout)
		return nil
	},
}

var accountBalancesCmd = &cobra.Command{
	Use: "balances",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument for account ID, received %d", len(args))
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parsing account id")
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			return errors.Wrap(err, "selecting account")
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)

		bs, err := c.SelectAccountBalances(*a)
		if err != nil {
			return errors.Wrap(err, "selecting account balances")
		}

		limit := viper.GetInt(keyLimit)
		if limit > len(*bs) {
			limit = len(*bs)
		}
		if limit != 0 {
			*bs = (*bs)[len(*bs)-limit:]
		}

		table.Balances(*bs, os.Stdout)
		return nil
	},
}

var balanceDate = date.Flag()
var accountBalanceInsertCmd = &cobra.Command{
	Use: "balance-insert",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument for account ID, received %d", len(args))
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parsing account id")
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			return errors.Wrap(err, "selecting account")
		}

		t := time.Now()
		if balanceDate.Time != nil {
			t = *balanceDate.Time
		}

		b, err := c.InsertBalance(*a, balance.Balance{
			Date:   t,
			Amount: viper.GetInt(keyAmount),
		})
		if err != nil {
			return errors.Wrap(err, "inserting balance")
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)
		table.Balances(storage.Balances{*b}, os.Stdout)
		return nil
	},
}

func init() {
	accountAddCmd.Flags().StringP(keyName, "n", "", "account name")
	accountAddCmd.Flags().VarP(accountOpened, keyOpened, "o", "account opened date")
	accountAddCmd.Flags().VarP(accountClosed, keyClosed, "c", "account closed date")
	accountAddCmd.Flags().String(keyCurrency, "EUR", "account currency")

	accountUpdateCmd.Flags().StringP(keyName, "n", "", "account name")
	accountUpdateCmd.Flags().VarP(accountOpened, keyOpened, "o", "account opened date")
	accountUpdateCmd.Flags().VarP(accountClosed, keyClosed, "c", "account closed date")
	accountUpdateCmd.Flags().String(keyCurrency, "EUR", "account currency")

	accountBalancesCmd.Flags().UintP(keyLimit, "l", 0, "limit results")

	// TODO: Stop multiple usage of the flag like in this article: http://blog.ralch.com/tutorial/golang-custom-flags/
	accountBalanceInsertCmd.Flags().VarP(balanceDate, keyDate, "d", "date of balance to insert")
	accountBalanceInsertCmd.Flags().IntP(keyAmount, "a", 0, "amount of balance to insert")

	rootCmd.AddCommand(accountCmd)
	for _, c := range []*cobra.Command{
		accountAddCmd,
		accountUpdateCmd,
		accountBalancesCmd,
		accountBalanceInsertCmd,
	} {
		err := viper.BindPFlags(c.Flags())
		if err != nil {
			log.Fatal(errors.Wrap(err, "binding pflags"))
		}
		accountCmd.AddCommand(c)
	}
}
