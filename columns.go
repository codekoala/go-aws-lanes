package lanes

import (
	"strings"
)

type Column string
type ColumnSet []Column

const (
	ColumnIndex       Column = "IDX"
	ColumnLane               = "LANE"
	ColumnServer             = "SERVER"
	ColumnIP                 = "IP ADDRESS"
	ColumnState              = "STATE"
	ColumnID                 = "ID"
	ColumnSSHIdentity        = "SSH_IDENTITY"
	ColumnUser               = "USER"
	ColumnInvalid            = ""
)

var DefaultColumnSet = ColumnSet{
	ColumnIndex,
	ColumnLane,
	ColumnServer,
	ColumnIP,
	ColumnState,
	ColumnID,
}

func ParseColumnSet(columns string) (out ColumnSet) {
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
	case string(ColumnSSHIdentity):
		c = ColumnSSHIdentity
	case string(ColumnUser):
		c = ColumnUser
	default:
		c = ColumnInvalid
	}

	return c
}

func GetDefaultColumnList() string {
	return DefaultColumnSet.String()
}

func (this ColumnSet) Add(columns ...Column) {
	this = append(this, columns...)
}

func (this ColumnSet) Remove(columns ...Column) (sanitized ColumnSet) {
	for _, col := range this {
		show := true
		for _, rem := range columns {
			if col == rem {
				show = false
				break
			}
		}

		if show {
			sanitized = append(sanitized, col)
		}
	}

	return sanitized
}

func (this ColumnSet) String() string {
	var names []string

	for _, col := range this {
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
