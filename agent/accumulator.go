package agent

import (
	"log"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"github.com/vikramjakhr/grafana-dashboard-exporter/metric"
)

type MetricMaker interface {
	Name() string
}

func NewAccumulator(
	maker MetricMaker,
	metrics chan gde.Metric,
) *accumulator {
	acc := accumulator{
		maker:   maker,
		metrics: metrics,
	}
	return &acc
}

type accumulator struct {
	metrics chan gde.Metric
	maker   MetricMaker
}

func (ac *accumulator) AddOutput(org string, valueType gde.ValueType, content []byte) {
	if org != "" && valueType != "" && len(content) > 0 {
		ac.metrics <- metric.New(org, valueType, content)
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
