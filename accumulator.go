package gde

// Accumulator is an interface for "accumulating" metrics from plugin(s).
// The metrics are sent down a channel shared between all plugins.
type Accumulator interface {
	AddOutput(dir string, valueType ValueType, action Action, title string, content []byte)

	AddError(err error)
}
