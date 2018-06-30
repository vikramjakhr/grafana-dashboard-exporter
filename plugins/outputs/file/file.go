package file

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"fmt"
)

type File struct {
	OutputDir    string `toml:"output_dir"`
	OutputFormat string `toml:"output_format"`
}

var sampleConfig = `
  output_dir = "<dir>" # default is /tmp/gde
  output_format = "zip" # zip, dir # default is zip
`

func (f *File) SampleConfig() string {
	return sampleConfig
}

func (f *File) Connect() error {
	return nil
}

func (f *File) Description() string {
	return "Send grafana json to specified directory"
}

func (f *File) Write(metric gde.Metric) error {
	if metric.Org() != "" && metric.Type() != "" && metric.Title() != "" && metric.Content() != "" {
		fmt.Println(metric.Org(), " | ", metric.Type(), " | ", metric.Title())
	}
	return nil
}

func init() {
	outputs.Add("file", func() gde.Output {
		return &File{}
	})
}
