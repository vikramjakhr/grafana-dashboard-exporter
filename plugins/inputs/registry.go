package inputs

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter"
)

type Creator func() gde.Input

var Inputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Inputs[name] = creator
}
