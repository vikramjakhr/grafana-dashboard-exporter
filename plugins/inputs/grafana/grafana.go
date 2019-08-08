package grafana

import (
	"encoding/json"
	"fmt"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs/grafana/api"
	"log"
	"strings"
	"time"
)

type Grafana struct {
	Host          string `toml:"host"`
	Authorization string `toml:"authorization"`
	Dashboard     bool   `toml:"dashboard"`
	Meta    	  bool   `toml:"meta"`
	Datasource    bool   `toml:"datasource"`
}

func (_ *Grafana) Description() string {
	return "Fetch grafana json from specified grafana host"
}

var sampleConfig = `
  host = "http://<host>:<port>" # required
  authorization = "Bearer <token>" # required
  dashboard = true # true if dashboard needs to be fetched; default true
  meta = true # true if dashboard metadata is required while fetching; default true
  datasource = true # true if datasource needs to be fetched; default true
`

func (_ *Grafana) SampleConfig() string {
	return sampleConfig
}

func (s *Grafana) Process(acc gde.Accumulator) error {
	if s.Datasource || s.Dashboard {
		gClient, err := api.NewGrafanaClient(s.Authorization, s.Host)
		if err != nil {
			return err
		}

		org, err := gClient.GetCurrentOrg()
		if err != nil {
			return err
		}

		tym := time.Now()

		dir := fmt.Sprintf("%s@%s",
			strings.Replace(org.Name, " ", "", -1),
			tym.Format("2006-January-2T15:04:05"))

		if s.Datasource {
			dSources, err := gClient.GetDataSources()
			if err != nil {
				return err
			}

			for _, ds := range *dSources {
				byts, err := json.Marshal(ds)
				if err != nil {
					return err
				}
				acc.AddOutput(dir, gde.TypeDatasource, gde.ActionCreate, ds.Name, byts)
			}
		}

		if s.Dashboard {
			results, err := gClient.Search(api.SearchTypeDashDB, "")
			if err != nil {
				return err
			}

			for _, db := range *results {
				dashboard, err := gClient.GetDashboard(db.Uri)
				if err != nil {
					return err
				}

				// Remove generic keys
				delete(dashboard.Model, "id")
				delete(dashboard.Model, "uid")
				delete(dashboard.Model, "version")

				var obj interface{} = dashboard
				if s.Meta {
					obj = dashboard.Model
				}
				byts, err := json.Marshal(obj)
				if err != nil {
					return err
				}
				name := dashboard.Model["title"].(string)
				acc.AddOutput(dir, gde.TypeDashboard, gde.ActionCreate, name, byts)
			}
		}

		acc.AddOutput(dir, "", gde.ActionFinish, "", nil)

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
