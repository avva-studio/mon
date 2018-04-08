package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/accounting-rest/pkg/table"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/balance"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyDate   = "date"
	keyAmount = "amount"
)

var accountCmd = &cobra.Command{
	Use: "account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("no account id given")
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

var accountBalancesCmd = &cobra.Command{
	Use: "balances",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("no account id given")
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

		table.Balances(*bs, os.Stdout)
		return nil
	},
}

var accountBalanceInsertCmd = &cobra.Command{
	Use: "balance-insert",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument, receieved %d", len(args))
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
		if d := viper.GetString(keyDate); d != "" {
			fmt.Println("parsing the d")
			t, err = parseDateString(d)
			if err != nil {
				return errors.Wrap(err, "parsing date string")
			}
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

func parseDateString(dateString string) (time.Time, error) {
	return time.Parse("2006-01-02", dateString)
}

func init() {
	accountBalanceInsertCmd.Flags().StringP(keyDate, "d", "", "date of balance to insert")
	accountBalanceInsertCmd.Flags().StringP(keyAmount, "a", "", "amount of balance to insert")
	err := viper.BindPFlags(accountBalanceInsertCmd.Flags())
	if err != nil {
		log.Fatal(errors.Wrap(err, "binding pflags"))
	}
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(accountBalancesCmd)
	accountCmd.AddCommand(accountBalanceInsertCmd)
}
