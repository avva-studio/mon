package main

import (
	"log"
	"os"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/accounting-rest/pkg/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdAccounts = &cobra.Command{
	Use: "accounts",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.Client(viper.GetString(keyServerHost))
		as, err := c.SelectAccounts()
		if err != nil {
			log.Fatal(errors.Wrap(err, "selecting accounts"))
		}
		table.Accounts(*as, os.Stdout)
	},
}

func init() {
	cmdRoot.AddCommand(cmdAccounts)
}
