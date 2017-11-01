package lanes

type Column string

const (
	ColumnIndex  Column = "IDX"
	ColumnLane          = "LANE"
	ColumnServer        = "SERVER"
	ColumnIP            = "IP ADDRESS"
	ColumnState         = "STATE"
	ColumnID            = "ID"
)

var DefaultColumns = []Column{
	ColumnIndex,
	ColumnLane,
	ColumnServer,
	ColumnIP,
	ColumnState,
	ColumnID,
}
