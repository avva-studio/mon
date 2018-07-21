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
		as, err := accounts(c)
		if err != nil {
			return errors.Wrap(err, "getting accounts")
		}

		if viper.GetBool(keyQuiet) {
			for _, a := range as {
				fmt.Println(a.ID)
			}
			return nil
		}

		table.Accounts(as, os.Stdout)
		return nil
	},
}

var accountsBalancesCmd = &cobra.Command{
	Use:  "balances",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if atDate.Time == nil {
			now := time.Now()
			atDate.Time = &now
		}

		c := client.Client(viper.GetString(keyServerHost))
		as, err := accounts(c)
		if err != nil {
			return errors.Wrap(err, "getting accounts")
		}

		abs, err := accountsBalances(c, as, *atDate.Time)
		if err != nil {
			return errors.Wrap(err, "getting balances for all accounts")
		}

		table.AccountsWithBalance(abs, os.Stdout)

		cbs := currencyBalances(abs)
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

func accounts(store storage.Storage) (storage.Accounts, error) {
	as, err := store.SelectAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "selecting accounts")
	}

	ac, err := prepareAccountCondition()
	if err != nil {
		return nil, errors.Wrap(err, "preparing conditions")
	}

	*as = ac.Filter(*as)

	if s, ok := sort.AccountSorts()[sortBy.String()]; ok {
		s(*as)
	}

	return *as, nil
}

func accountsBalances(store storage.Storage, as storage.Accounts, at time.Time) ([]accountbalance.AccountBalance, error) {
	var abs []accountbalance.AccountBalance
	for _, a := range as {
		bs, err := store.SelectAccountBalances(a)
		if err != nil {
			return nil, errors.Wrapf(err, "selecting balances for account: %+v", a)
		}
		bbs := bs.InnerBalances()
		if len(*bs) == 0 {
			log.Printf("no balances for account:%+v", a)
			continue
		}
		b, err := bbs.AtTime(at)
		if err != nil {
			log.Println(errors.Wrapf(err, "getting balances at time:%+v for account:%+v", at, a))
			continue
		}
		abs = append(abs, accountbalance.AccountBalance{
			Account: a,
			Balance: b,
		})
	}

	if s, ok := sort.AccountbalanceSorts()[sortBy.String()]; ok {
		s(abs)
	}

	return abs, nil
}

func currencyBalances(abs []accountbalance.AccountBalance) map[currency.Code]balance.Balances {
	cbs := make(map[currency.Code]balance.Balances)
	for _, ab := range abs {
		crncy := ab.Account.Account.CurrencyCode()
		if _, ok := cbs[crncy]; !ok {
			cbs[crncy] = balance.Balances{}
		}
		cbs[crncy] = append(cbs[crncy], ab.Balance)
	}
	return cbs
}

func prepareAccountCondition() (filter.AccountCondition, error) {
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
	accountsCmd.PersistentFlags().Bool(keyOpen, false, "show only open accounts")
	accountsCmd.PersistentFlags().UintSliceVar(&ids, keyIDs, []uint{}, "filter by ids")
	accountsCmd.PersistentFlags().StringSliceVar(&currencies, keyCurrencies, []string{}, "filter by currencies")
	accountsCmd.Flags().BoolP(keyQuiet, "q", false, "show only account ids")
	accountsCmd.PersistentFlags().Var(atDate, keyAtDate, "show balances at a certain date")
	accountsCmd.PersistentFlags().Var(sortBy, keySortBy, fmt.Sprintf("sort by one of %s", strings.Join(sort.AllKeys(), ",")))

	accountsCmd.AddCommand(accountsBalancesCmd)

	for _, cc := range []*cobra.Command{
		accountsCmd, accountsBalancesCmd,
	} {
		err := bindAllFlags(cc)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "binding command:[%s] flags", cc.Use))
		}
	}
}

func bindAllFlags(c *cobra.Command) error {
	if err := viper.BindPFlags(c.Flags()); err != nil {
		return errors.Wrapf(err, "binding local pflags")
	}
	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		return errors.Wrapf(err, "binding persistent pflags")
	}
	return nil
}
