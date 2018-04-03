package main

import (
	"fmt"
	"log"
	"os"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/accounting-rest/pkg/filter"
	"github.com/glynternet/accounting-rest/pkg/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyOpen  = "open"
	keyQuiet = "quiet"
)

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
		if viper.GetBool(keyQuiet) {
			for _, a := range *as {
				fmt.Println(a.ID)
			}
			return
		}
		table.Accounts(*as, os.Stdout)
	},
}

func init() {
	cmdRoot.AddCommand(cmdAccounts)
	cmdAccounts.Flags().BoolP(keyOpen, "", false, "show only open accounts")
	cmdAccounts.Flags().BoolP(keyQuiet, "q", false, "show only account ids")
	if err := viper.BindPFlags(cmdAccounts.Flags()); err != nil {
		log.Fatal(errors.Wrap(err, "binding pflags"))
	}
}
