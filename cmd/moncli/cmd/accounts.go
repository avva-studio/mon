package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/accounting-rest/pkg/filter"
	"github.com/glynternet/accounting-rest/pkg/storage"
	"github.com/glynternet/accounting-rest/pkg/table"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	gtime "github.com/glynternet/go-time"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyOpen     = "open"
	keyQuiet    = "quiet"
	keyBalances = "balances"
	keyAtDate   = "at-date"
)

var accountsCmd = &cobra.Command{
	Use: "accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := parseNullTime(viper.GetString(keyAtDate))
		if err != nil {
			return errors.Wrapf(err, "parsing %s", keyAtDate)
		}
		if !t.Valid {
			t = gtime.NullTime{Valid: true, Time: time.Now()}
		}

		c := client.Client(viper.GetString(keyServerHost))
		as, err := c.SelectAccounts()
		if err != nil {
			return errors.Wrap(err, "selecting accounts")
		}

		*as = filter.Filter(*as, filter.Existed(t.Time))

		if viper.GetBool(keyOpen) {
			*as = filter.Filter(*as, filter.OpenAt(t.Time))
		}
		if viper.GetBool(keyQuiet) {
			for _, a := range *as {
				fmt.Println(a.ID)
			}
			return nil
		}
		if viper.GetBool(keyBalances) {
			// TODO: Should these be kept in order of ID here?
			abs := make(map[storage.Account]balance.Balance)
			cbs := make(map[currency.Code]balance.Balances)
			for _, a := range *as {
				bs, err := c.SelectAccountBalances(a)
				if err != nil {
					return errors.Wrapf(err, "selecting balances for account: %+v", a)
				}
				if len(*bs) == 0 {
					log.Printf("no balances for account:%+v", a)
					continue
				}
				var bbs balance.Balances
				for _, b := range *bs {
					bbs = append(bbs, b.Balance)
				}
				current, err := bbs.AtTime(t.Time)
				if err != nil {
					log.Println(errors.Wrapf(err, "getting balances at time:%+v for account:%+v", t, a))
					continue
				}
				abs[a] = current

				crncy := a.Account.CurrencyCode()
				if _, ok := cbs[crncy]; !ok {
					cbs[crncy] = balance.Balances{}
				}
				cbs[crncy] = append(cbs[crncy], current)
			}
			table.AccountsWithBalance(abs, os.Stdout)
			if len(cbs) == 0 {
				return nil
			}
			totals := [][]string{{"Currency", "Amount"}}
			for crncy, bs := range cbs {
				totals = append(totals, []string{crncy.String(), strconv.Itoa(bs.Sum())})
			}
			table.Basic(totals, os.Stdout)
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
	accountsCmd.Flags().BoolP(keyBalances, "b", false, "show balances for each account")
	accountsCmd.Flags().String(keyAtDate, "", "show balances at a certain date")
	if err := viper.BindPFlags(accountsCmd.Flags()); err != nil {
		log.Fatal(errors.Wrap(err, "binding pflags"))
	}
}
