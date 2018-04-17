package table

import (
	"io"
	"strconv"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-time"
	"github.com/olekukonko/tablewriter"
)

const dateFormat = `02-01-2006`

func Accounts(as storage.Accounts, w io.Writer) {
	t := newDefaultTable(w)
	t.SetHeader([]string{"ID", "Name", "Opened", "Closed", "Currency"})

	for _, a := range as {
		t.Append([]string{
			strconv.FormatUint(uint64(a.ID), 10),
			a.Account.Name(),
			a.Account.Opened().Format(dateFormat),
			closedString(a.Account.Closed()),
			a.Account.CurrencyCode().String(),
		})
	}
	t.Render() // Send output
}

func AccountsWithBalance(abs map[storage.Account]balance.Balance, w io.Writer) {
	t := newDefaultTable(w)
	t.SetHeader([]string{"ID", "Name", "Opened", "Closed", "Currency", "Balance"})

	for a, b := range abs {
		t.Append([]string{
			strconv.FormatUint(uint64(a.ID), 10),
			a.Account.Name(),
			a.Account.Opened().Format(dateFormat),
			closedString(a.Account.Closed()),
			a.Account.CurrencyCode().String(),
			strconv.Itoa(b.Amount),
		})
	}
	t.Render() // Send output
}

func newDefaultTable(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetAutoWrapText(false)
	return table
}

func closedString(t time.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(dateFormat)
}

func Balances(bs storage.Balances, w io.Writer) {
	t := newDefaultTable(w)
	t.SetHeader([]string{"ID", "Amount", "Date"})

	for _, b := range bs {
		t.Append([]string{
			strconv.FormatUint(uint64(b.ID), 10),
			strconv.Itoa(b.Amount),
			b.Date.Format(dateFormat),
		})
	}

	t.Render()
}
