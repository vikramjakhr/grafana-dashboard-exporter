package file

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"log"
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
	if metric.Dir() != "" && metric.Type() != "" && metric.Title() != "" && len(metric.Content()) > 0 {
		dir := "/tmp"

		if strings.Trim(f.OutputDir, " ") != "" {
			dir = strings.TrimRight(f.OutputDir, "/")
		}

		dir = fmt.Sprintf("%s/%s/%ss/", dir, metric.Dir(), string(metric.Type()))

		fmt.Println(dir)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0774)
			if err != nil {
				log.Printf("E! Unable to create direcotry. %v", err)
				return err
			}
		}

		switch metric.Type() {
		case gde.Datasource:
			filename := fmt.Sprintf("%s%s.json", dir, strings.Replace(metric.Title(), " ", "", -1))
			err := ioutil.WriteFile(filename, metric.Content(), 0644)
			if err != nil {
				log.Printf("E! Unable to create file. %v", err)
				return err
			}
			break
		case gde.Dashboard:
			filename := fmt.Sprintf("%s%s.json", dir, strings.Replace(metric.Title(), " ", "", -1))
			err := ioutil.WriteFile(filename, metric.Content(), 0644)
			if err != nil {
				log.Printf("E! Unable to create file. %v", err)
				return err
			}
			break
		}

	}
	return nil
}

func init() {
	outputs.Add("file", func() gde.Output {
		return &File{}
	})
}
