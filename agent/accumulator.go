package agent

import (
	"log"
)

type MetricMaker interface {
	Name() string
}

func NewAccumulator(
	maker MetricMaker,
	metrics chan string,
) *accumulator {
	acc := accumulator{
		maker:   maker,
		metrics: metrics,
	}
	return &acc
}

type accumulator struct {
	metrics chan string
	maker   MetricMaker
}

func (ac *accumulator) AddFile(file string) {
	if file != "" {
		ac.metrics <- file
	}
}

// AddError passes a runtime error to the accumulator.
// The error will be tagged with the plugin name and written to the log.
func (ac *accumulator) AddError(err error) {
	if err == nil {
		return
	}
	log.Printf("E! Error in plugin [%s]: %s", ac.maker.Name(), err)
}
