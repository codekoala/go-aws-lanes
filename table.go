package lanes

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/apcera/termtables"
)

type Rowish interface {
	AddCell(value interface{})
}

type Row struct {
	Cells []string
}

func (this *Row) AddCell(value interface{}) {
	this.Cells = append(this.Cells, fmt.Sprintf("%s", value))
}

func createServerLine(columns ColumnSet, row Rowish, idx int, svr *Server) {
	for _, col := range columns {
		switch col {
		case ColumnIndex:
			row.AddCell(idx + 1)
		case ColumnLane:
			row.AddCell(svr.Lane)
		case ColumnServer:
			row.AddCell(svr.Name)
		case ColumnIP:
			row.AddCell(svr.IP)
		case ColumnState:
			row.AddCell(svr.State)
		case ColumnID:
			row.AddCell(svr.ID)
		case ColumnSSHIdentity:
			row.AddCell(svr.profile.Identity)
		case ColumnUser:
			row.AddCell(svr.profile.GetUser())
		default:
			continue
		}
	}
}

func DisplayServers(servers []*Server) error {
	return DisplayServersWriter(os.Stdout, servers)
}

func DisplayServersCols(servers []*Server, columns ColumnSet) error {
	return DisplayServersColsWriter(os.Stdout, servers, columns)
}

func DisplayServersWriter(writer io.Writer, servers []*Server) (err error) {
	return DisplayServersColsWriter(writer, servers, DefaultColumnSet)
}

func DisplayServersCSVWriter(writer *csv.Writer, servers []*Server, columns ColumnSet) (err error) {
	for idx, svr := range servers {
		var row Row
		createServerLine(columns, &row, idx, svr)

		if err = writer.Write(row.Cells); err != nil {
			return
		}
	}

	return nil
}

func DisplayServersColsWriter(writer io.Writer, servers []*Server, columns ColumnSet) (err error) {
	if len(servers) == 0 {
		return fmt.Errorf("No servers found.")
	}

	if len(columns) == 0 {
		return nil
	}

	if !config.DisableUTF8 {
		termtables.EnableUTF8()
	}

	table := termtables.CreateTable()
	if config.Table.HideBorders {
		table.Style.SkipBorder = true
		table.Style.BorderY = ""
		table.Style.PaddingLeft = 0
	}

	if !config.Table.HideTitle {
		table.AddTitle("AWS Servers")
	}

	for idx, svr := range servers {
		row := table.AddRow()
		createServerLine(columns, row, idx, svr)
	}

	if !config.Table.HideHeaders {
		// add headers after all rows because cell alignment only applies to cells that exist when SetAlign is called
		for idx, col := range columns {
			table.AddHeaders(col)

			switch col {
			case ColumnIndex:
				table.SetAlign(termtables.AlignRight, idx+1)
			}
		}
	}

	fmt.Fprintf(writer, table.Render())

	return nil
}
