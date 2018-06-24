package file

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
)

type File struct{}

var sampleConfig = `
  output_dir = "<dir>" # default is /tmp/gde
  output_format = "zip" # zip, dir # default is zip
`

func (f *File) SampleConfig() string {
	return sampleConfig
}

func (f *File) Description() string {
	return "Send grafana json to specified directory"
}

func (f *File) Write() error {
	return nil
}

func init() {
	outputs.Add("file", func() gde.Output {
		return &File{}
	})
}
