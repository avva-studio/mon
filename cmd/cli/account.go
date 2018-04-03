package main

import (
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

var cmdAccount = &cobra.Command{
	Use: "account",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal(errors.New("no account id given"))
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal(errors.Wrap(err, "parsing account id"))
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			log.Fatal(errors.Wrap(err, "selecting account"))
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)
	},
}

var cmdAccountBalances = &cobra.Command{
	Use: "balances",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal(errors.New("no account id given"))
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal(errors.Wrap(err, "parsing account id"))
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			log.Fatal(errors.Wrap(err, "selecting account"))
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)

		bs, err := c.SelectAccountBalances(*a)
		if err != nil {
			log.Fatal(errors.Wrap(err, "selecting account balances"))
		}

		table.Balances(*bs, os.Stdout)
	},
}

var cmdAccountBalanceInsert = &cobra.Command{
	Use: "balance-insert",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal(errors.New("no account id and amount given"))
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal(errors.Wrap(err, "parsing account id"))
		}
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(errors.Wrap(err, "parsing balance amount"))
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			log.Fatal(errors.Wrap(err, "selecting account"))
		}

		b, err := c.InsertBalance(*a, balance.Balance{
			Date:   time.Now(),
			Amount: amount,
		})
		if err != nil {
			log.Fatal(errors.Wrap(err, "inserting balance"))
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)
		table.Balances(storage.Balances{*b}, os.Stdout)
	},
}

func init() {
	cmdRoot.AddCommand(cmdAccount)
	cmdAccount.AddCommand(cmdAccountBalances)
	cmdAccount.AddCommand(cmdAccountBalanceInsert)
}
