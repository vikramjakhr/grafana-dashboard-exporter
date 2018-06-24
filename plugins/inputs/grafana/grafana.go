package grafana

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"fmt"
)

type Grafana struct {
	Host          string `toml:"host"`
	Authorization string `toml:"authorization"`
	Dashboard     bool `toml:"dashboard"`
	Datasource    bool `toml:"datasource"`
	Org           bool `toml:"org"`
}

func (_ *Grafana) Description() string {
	return "Fetch grafana json from specified grafana host"
}

var sampleConfig = `
  host = "http://<host>:<port>" # required
  authorization = "Bearer <token>" # required
  dashboard = true # true if dashboard needs to be fetched; default true
  datasource = true # true if datasource needs to be fetched; default true
  org = true # true if organization needs to be fetched; default true
`

func (_ *Grafana) SampleConfig() string {
	return sampleConfig
}

func (s *Grafana) Process(acc gde.Accumulator) error {
	fmt.Println("collecting...")
	acc.AddFile("fileeeee")
	return nil
}

func init() {
	inputs.Add("grafana", func() gde.Input {
		return &Grafana{}
	})
}
