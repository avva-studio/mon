package main

import (
	"log"
	"os"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/accounting-rest/pkg/filter"
	"github.com/glynternet/accounting-rest/pkg/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const keyOpen = "open"

var cmdAccounts = &cobra.Command{
	Use: "accounts",
	Run: func(cmd *cobra.Command, args []string) {
		as, err := client.Client(viper.GetString(keyServerHost)).SelectAccounts()
		if err != nil {
			log.Fatal(errors.Wrap(err, "selecting accounts"))
		}
		if viper.GetBool(keyOpen) {
			*as = filter.Filter(*as, filter.Open())
		}
		table.Accounts(*as, os.Stdout)
	},
}

func init() {
	cmdRoot.AddCommand(cmdAccounts)
	cmdAccounts.Flags().BoolP(keyOpen, "", false, "show only open accounts")
	if err := viper.BindPFlags(cmdAccounts.Flags()); err != nil {
		log.Fatal(errors.Wrap(err, "binding pflags"))
	}
}
