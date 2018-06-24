package outputs

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter"
)

type Creator func() gde.Output

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}
