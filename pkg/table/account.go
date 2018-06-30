package table

import (
	"fmt"
	"io"
	"strconv"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-time"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/stringgrid"
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
	t.SetHeader([]string{"ID", "Name", "Opened", "Closed", "Currency", "Balance Date", "Balance Amount"})

	for a, b := range abs {
		t.Append([]string{
			strconv.FormatUint(uint64(a.ID), 10),
			a.Account.Name(),
			a.Account.Opened().Format(dateFormat),
			closedString(a.Account.Closed()),
			a.Account.CurrencyCode().String(),
			b.Date.Format(dateFormat),
			strconv.Itoa(b.Amount),
		})
	}
	t.Render() // Send output
}

func Basic(data [][]string, w io.Writer) error {
	if len(data) < 2 {
		return fmt.Errorf("requires at least 2 rows of data")
	}
	t := newDefaultTable(w)
	t.SetHeader(data[0])
	for _, row := range data[1:] {
		t.Append(row)
	}
	t.Render()
	return nil
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

	g := stringgrid.Columns{
		stringgrid.SimpleColumn(func(i uint) string {
			return strconv.FormatUint(uint64(bs[i].ID), 10)
		}),
		stringgrid.SimpleColumn(func(i uint) string {
			return strconv.Itoa(bs[i].Amount)
		}),
		stringgrid.SimpleColumn(func(i uint) string {
			return bs[i].Date.Format(dateFormat)
		}),
	}

	data, err := g.Generate(uint(len(bs)))
	if err != nil {
		return nil, nil
	}

	for _, row := range  {

	}

	for _, b := range bs {
		t.Append([]string{
			strconv.FormatUint(uint64(b.ID), 10),
			strconv.Itoa(b.Amount),
			b.Date.Format(dateFormat),
		})
	}

	t.Render()
}
