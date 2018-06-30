package gde

// ValueType is an enumeration of metric types that represent a simple value.
type ValueType string
type Action string

// Possible values for the ValueType enum.
const (
	TypeDatasource ValueType = "Datasource"
	TypeDashboard  ValueType = "Dashboard"
	ActionCreate   Action    = "Create"
	ActionZIP      Action    = "ZIP"
)

type Metric interface {
	// Getting data structure functions
	Dir() string
	Type() ValueType
	Action() Action
	Title() string
	Content() []byte
}
