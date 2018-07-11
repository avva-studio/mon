package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	"github.com/glynternet/mon/client"
	"github.com/glynternet/mon/pkg/date"
	"github.com/glynternet/mon/pkg/filter"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyOpen       = "open"
	keyIDs        = "ids"
	keyCurrencies = "currencies"
	keyQuiet      = "quiet"
	keyBalances   = "balances"
	keyAtDate     = "at-date"
)

var (
	atDate     = date.Flag()
	ids        []uint
	currencies []string
)

var accountsCmd = &cobra.Command{
	Use:  "accounts",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if atDate.Time == nil {
			now := time.Now()
			atDate.Time = &now
		}

		c := client.Client(viper.GetString(keyServerHost))
		as, err := c.SelectAccounts()
		if err != nil {
			return errors.Wrap(err, "selecting accounts")
		}

		fs, err := prepareFilters()
		if err != nil {
			return errors.Wrap(err, "preparing filters")
		}

		for _, f := range fs {
			*as = filter.Filter(*as, f)
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
				current, err := bbs.AtTime(*atDate.Time)
				if err != nil {
					log.Println(errors.Wrapf(err, "getting balances at time:%+v for account:%+v", *atDate.Time, a))
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
			return errors.Wrap(table.Basic(totals, os.Stdout), "printing basic table for totals")
		}
		table.Accounts(*as, os.Stdout)
		return nil
	},
}

func prepareFilters() ([]filter.AccountFilter, error) {
	fs := []filter.AccountFilter{
		filter.Existed(*atDate.Time),
	}

	if viper.GetBool(keyOpen) {
		fs = append(fs, filter.OpenAt(*atDate.Time))
	}

	if len(ids) > 0 {
		fs = append(fs, filter.IDs(ids))
	}

	if len(currencies) > 0 {
		cs, err := currencyStringsToCodes(currencies...)
		if err != nil {
			return nil, errors.Wrap(err, "converting currency string to currency codes")
		}
		fs = append(fs, filter.Currencies(cs...))
	}

	return fs, nil
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.Flags().BoolP(keyOpen, "", false, "show only open accounts")
	accountsCmd.Flags().UintSliceVar(&ids, keyIDs, []uint{}, "filter by ids")
	accountsCmd.Flags().StringSliceVar(&currencies, keyCurrencies, []string{}, "filter by currencies")
	accountsCmd.Flags().BoolP(keyQuiet, "q", false, "show only account ids")
	accountsCmd.Flags().BoolP(keyBalances, "b", false, "show balances for each account")
	accountsCmd.Flags().Var(atDate, keyAtDate, "show balances at a certain date")
	if err := viper.BindPFlags(accountsCmd.Flags()); err != nil {
		log.Fatal(errors.Wrap(err, "binding pflags"))
	}
}

// TODO: this should be handled as a flags type perhaps?
func currencyStringsToCodes(css ...string) ([]currency.Code, error) {
	var codes []currency.Code
	for _, cs := range css {
		c, err := currency.NewCode(cs)
		if err != nil {
			return nil, errors.Wrap(err, "creating new code")
		}
		codes = append(codes, *c)
	}
	return codes, nil
}
