package grafana

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"log"
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs/grafana/api"
	"encoding/json"
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
	log.Printf("collecting...")
	if s.Org || s.Datasource || s.Dashboard {
		gClient, err := api.NewGrafanaClient(s.Authorization, s.Host)
		if err != nil {
			return err
		}

		if s.Org {
			org, err := gClient.GetCurrentOrg()
			if err != nil {
				return err
			}
			out, err := json.Marshal(org)
			if err != nil {
				return err
			}
			acc.AddFile(string(out))
		}

		if s.Datasource {
			org, err := gClient.GetDataSources()
			if err != nil {
				return err
			}
			out, err := json.Marshal(org)
			if err != nil {
				return err
			}
			acc.AddFile(string(out))
		}
	} else {
		log.Printf("E! Error in grafana input plugin. Atleast one of Org, Datasource and Dashboard must be true.")
	}
	return nil
}

func init() {
	inputs.Add("grafana", func() gde.Input {
		return &Grafana{}
	})
}
