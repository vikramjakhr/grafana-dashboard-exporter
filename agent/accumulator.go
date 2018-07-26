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

func (ac *accumulator) AddOutput(dir string, valueType gde.ValueType, action gde.Action, title string, content []byte) {
	if action != "" {
		switch action {
		case gde.ActionCreate:
			if dir != "" && valueType != "" && title != "" && len(content) > 0 {
				ac.metrics <- metric.New(dir, valueType, action, title, content)
			}
			break
		case gde.ActionFinish:
			if dir != ""{
				ac.metrics <- metric.New(dir, valueType, action, title, content)
			}
			break
		}
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
