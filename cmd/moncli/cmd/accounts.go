package cmd

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

var accountsCmd = &cobra.Command{
	Use: "accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		as, err := client.Client(viper.GetString(keyServerHost)).SelectAccounts()
		if err != nil {
			return errors.Wrap(err, "selecting accounts")
		}
		if viper.GetBool(keyOpen) {
			*as = filter.Filter(*as, filter.Open())
		}
		if viper.GetBool(keyQuiet) {
			for _, a := range *as {
				fmt.Println(a.ID)
			}
			return nil
		}
		table.Accounts(*as, os.Stdout)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.Flags().BoolP(keyOpen, "", false, "show only open accounts")
	accountsCmd.Flags().BoolP(keyQuiet, "q", false, "show only account ids")
	if err := viper.BindPFlags(accountsCmd.Flags()); err != nil {
		log.Fatal(errors.Wrap(err, "binding pflags"))
	}
}
