package gde

// ValueType is an enumeration of metric types that represent a simple value.
type ValueType string

// Possible values for the ValueType enum.
const (
	Datasource ValueType = "Datasource"
	Dashboard  ValueType = "Dashboard"
)

type Metric interface {
	// Getting data structure functions
	Dir() string
	Type() ValueType
	Title() string
	Content() string
}
