package config

import (
	"time"
)

type Config struct {
	Agent   *AgentConfig
	Inputs  []*RunningInput
	Outputs []*RunningOutput
}

func NewConfig() *Config {
	c := &Config{
		// Agent defaults:
		Agent: &AgentConfig{
			Interval:      10 * time.Second,
			RoundInterval: true,
		},

		Inputs:  make([]*RunningInput, 0),
		Outputs: make([]*RunningOutput, 0),
	}
	return c
}

type RunningInput struct {
	Input  Input
}

type Input interface {
	// SampleConfig returns the default configuration of the Input
	SampleConfig() string

	// Description returns a one-sentence description on the Input
	Description() string

	// Process processes the input every "interval"
	Process() error
}


// RunningOutput contains the output configuration
type RunningOutput struct {
	Name   string
	Output Output
}

type Output interface {
	// Description returns a one-sentence description on the Output
	Description() string
	// SampleConfig returns the default configuration of the Output
	SampleConfig() string
	// Write takes in group of points to be written to the Output
	Write() error
}


type AgentConfig struct {
	// Interval at which to gather information
	Interval time.Duration

	// RoundInterval rounds collection interval to 'interval'.
	//     ie, if Interval=10s then always collect on :00, :10, :20, etc.
	RoundInterval bool

	// Debug is the option for running in debug mode
	Debug bool

	// Logfile specifies the file to send logs to
	Logfile string

	// Quiet is the option for running in quiet mode
	Quiet bool
}
