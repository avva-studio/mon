package table

import (
	"io"
	"strconv"

	"github.com/glynternet/go-accounting-storage"
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
			a.Name(),
			a.Opened().Format(dateFormat),
			closedString(a.Closed()),
			a.CurrencyCode().String(),
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
