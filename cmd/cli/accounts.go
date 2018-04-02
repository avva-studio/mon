package main

import (
	"log"
	"os"
	"strconv"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-time"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const dateFormat = `02-01-2006`

var cmdAccounts = &cobra.Command{
	Use: "accounts",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.Client(viper.GetString(keyServerHost))
		as, err := c.SelectAccounts()
		if err != nil {
			log.Fatal(errors.Wrap(err, "selecting accounts"))
		}
		table(*as)
	},
}

func init() {
	cmdRoot.AddCommand(cmdAccounts)
}

func table(as storage.Accounts) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Opened", "Closed", "Currency"})

	for _, a := range as {
		table.Append([]string{
			strconv.FormatUint(uint64(a.ID), 10),
			a.Name(),
			a.Opened().Format(dateFormat),
			closedString(a.Closed()),
			a.CurrencyCode().String(),
		})
	}
	table.Render() // Send output
}

func closedString(t time.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(dateFormat)
}
