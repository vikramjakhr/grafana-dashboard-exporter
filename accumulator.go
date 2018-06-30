package gde

// Accumulator is an interface for "accumulating" metrics from plugin(s).
// The metrics are sent down a channel shared between all plugins.
type Accumulator interface {
	AddOutput(org string, valueType ValueType, title string, content []byte)

	AddError(err error)
}
