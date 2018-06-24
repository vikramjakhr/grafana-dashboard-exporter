package gde

type Input interface {
	// SampleConfig returns the default configuration of the Input
	SampleConfig() string

	// Description returns a one-sentence description on the Input
	Description() string

	// Process processes the input every "interval"
	Process(Accumulator) error
}
