package file

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"fmt"
	"os"
	"strings"
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
	if strings.Trim(f.OutputDir, " ") != "" {
		_, err := os.Stat(f.OutputDir)
		return err
	}
	return nil
}

func (f *File) Description() string {
	return "Send grafana json to specified directory"
}

func (f *File) Write(metric gde.Metric) error {
	if metric.Dir() != "" && metric.Type() != "" && metric.Title() != "" && metric.Content() != "" {
		fmt.Println(metric.Dir(), " | ", metric.Type(), " | ", metric.Title())
		dir := "/tmp"

		if strings.Trim(f.OutputDir, " ") != "" {
			dir = strings.TrimRight(f.OutputDir, "/")
		}

		dir = dir + "/" + metric.Dir()

		fmt.Println(dir)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, 0664)
		}

	}
	return nil
}

func init() {
	outputs.Add("file", func() gde.Output {
		return &File{}
	})
}
