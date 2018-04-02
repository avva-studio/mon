package main

import (
	"log"
	"os"
	"strconv"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/accounting-rest/pkg/table"
	"github.com/glynternet/go-accounting-storage"
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
		idString := args[0]
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			log.Fatal(errors.Wrap(err, "parsing account id"))
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			log.Fatal(errors.Wrap(err, "selecting account"))
		}

		bs, err := c.SelectAccountBalances(*a)
		if err != nil {
			log.Fatal(errors.Wrap(err, "selecting account balances"))
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)
		table.Balances(*bs, os.Stdout)
	},
}

func init() {
	cmdRoot.AddCommand(cmdAccount)
}
