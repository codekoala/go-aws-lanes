package lanes

import (
	"strings"
)

type Column string

const (
	ColumnIndex   Column = "IDX"
	ColumnLane           = "LANE"
	ColumnServer         = "SERVER"
	ColumnIP             = "IP ADDRESS"
	ColumnState          = "STATE"
	ColumnID             = "ID"
	ColumnInvalid        = ""
)

var DefaultColumns = []Column{
	ColumnIndex,
	ColumnLane,
	ColumnServer,
	ColumnIP,
	ColumnState,
	ColumnID,
}

func ParseColumnList(columns string) (out []Column) {
	for _, cn := range strings.Split(columns, ",") {
		col := ColumnFromName(cn)
		if col == ColumnInvalid {
			continue
		}

		out = append(out, col)
	}

	return out
}

func ColumnFromName(name string) (c Column) {
	switch name {
	case string(ColumnIndex):
		c = ColumnIndex
	case string(ColumnLane):
		c = ColumnLane
	case string(ColumnServer), "NAME":
		c = ColumnServer
	case string(ColumnIP), "IP":
		c = ColumnIP
	case string(ColumnState):
		c = ColumnState
	case string(ColumnID):
		c = ColumnID
	default:
		c = ColumnInvalid
	}

	return c
}

func GetColumnList(columns ...Column) string {
	var names []string
	for _, col := range columns {
		var name string
		switch col {
		case ColumnIP:
			name = "IP"
		default:
			name = string(col)
		}

		names = append(names, name)
	}

	return strings.Join(names, ",")
}

func GetDefaultColumnList() string {
	return GetColumnList(DefaultColumns...)
}
