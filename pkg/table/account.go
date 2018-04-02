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
	table := tablewriter.NewWriter(w)
	table.SetAutoWrapText(false)
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
