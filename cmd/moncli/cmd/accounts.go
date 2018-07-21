package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	"github.com/glynternet/mon/internal/accountbalance"
	"github.com/glynternet/mon/internal/client"
	"github.com/glynternet/mon/internal/sort"
	"github.com/glynternet/mon/pkg/date"
	"github.com/glynternet/mon/pkg/filter"
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
	keySortBy     = "sort-by"
)

var (
	atDate     = date.Flag()
	sortBy     = sort.NewKey()
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

		ac, err := prepareCondition()
		if err != nil {
			return errors.Wrap(err, "preparing conditions")
		}

		*as = ac.Filter(*as)

		if s, ok := sort.AccountSorts()[sortBy.String()]; ok {
			s(*as)
		}

		if viper.GetBool(keyQuiet) {
			for _, a := range *as {
				fmt.Println(a.ID)
			}
			return nil
		}

		if !viper.GetBool(keyBalances) {
			table.Accounts(*as, os.Stdout)
			return nil
		}

		var abs []accountbalance.AccountBalance
		for _, a := range *as {
			bs, err := c.SelectAccountBalances(a)
			if err != nil {
				return errors.Wrapf(err, "selecting balances for account: %+v", a)
			}
			bbs := bs.InnerBalances()
			if len(*bs) == 0 {
				log.Printf("no balances for account:%+v", a)
				continue
			}
			current, err := bbs.AtTime(*atDate.Time)
			if err != nil {
				log.Println(errors.Wrapf(err, "getting balances at time:%+v for account:%+v", *atDate.Time, a))
				continue
			}
			abs = append(abs, accountbalance.AccountBalance{
				Account: a,
				Balance: current,
			})
		}

		if s, ok := sort.AccountbalanceSorts()[sortBy.String()]; ok {
			s(abs)
		}

		cbs := make(map[currency.Code]balance.Balances)
		for _, ab := range abs {
			crncy := ab.Account.Account.CurrencyCode()
			if _, ok := cbs[crncy]; !ok {
				cbs[crncy] = balance.Balances{}
			}
			cbs[crncy] = append(cbs[crncy], ab.Balance)
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
	},
}

func prepareCondition() (filter.AccountCondition, error) {
	cs := filter.AccountConditions{
		filter.Existed(*atDate.Time),
	}

	if viper.GetBool(keyOpen) {
		cs = append(cs, filter.OpenAt(*atDate.Time))
	}

	if len(ids) > 0 {
		cs = append(cs, idsCondition(ids))
	}

	if len(currencies) > 0 {
		ccs, err := currenciesCondition(currencies)
		if err != nil {
			return nil, errors.Wrap(err, "creating currencies condition")
		}
		cs = append(cs, ccs)
	}

	// Account must meet all AccountConditions
	return cs.And, nil
}

func idsCondition(ids []uint) filter.AccountCondition {
	var cs filter.AccountConditions
	for _, id := range ids {
		cs = append(cs, filter.ID(id))
	}
	return cs.Or
}

func currenciesCondition(cs []string) (filter.AccountCondition, error) {
	ccs, err := currencyStringsToCodes(cs...)
	if err != nil {
		return nil, errors.Wrap(err, "converting currency string to currency codes")
	}
	var acs filter.AccountConditions
	for _, c := range ccs {
		acs = append(acs, filter.Currency(c))
	}
	return acs.Or, nil
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

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.Flags().Bool(keyOpen, false, "show only open accounts")
	accountsCmd.Flags().UintSliceVar(&ids, keyIDs, []uint{}, "filter by ids")
	accountsCmd.Flags().StringSliceVar(&currencies, keyCurrencies, []string{}, "filter by currencies")
	accountsCmd.Flags().BoolP(keyQuiet, "q", false, "show only account ids")
	accountsCmd.Flags().BoolP(keyBalances, "b", false, "show balances for each account")
	accountsCmd.Flags().Var(atDate, keyAtDate, "show balances at a certain date")
	accountsCmd.Flags().Var(sortBy, keySortBy, fmt.Sprintf("sort by one of %s", strings.Join(sort.AllKeys(), ",")))
	if err := viper.BindPFlags(accountsCmd.Flags()); err != nil {
		log.Fatal(errors.Wrap(err, "binding pflags"))
	}
}
