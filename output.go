package gde

type Output interface {
	// Connect to the Output
	Connect() error
	// Description returns a one-sentence description on the Output
	Description() string
	// SampleConfig returns the default configuration of the Output
	SampleConfig() string
	// Write takes in group of points to be written to the Output
	Write(metric Metric) error
}
