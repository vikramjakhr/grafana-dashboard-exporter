package s3

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
)

type S3 struct {
	Bucket       string `toml:"bucket"`
	AccessKey    string `toml:"access_key"`
	SecretKey    string `toml:"secret_key"`
	OutputFormat string `toml:"output_format"`
}

var sampleConfig = `
  bucket = "<bucket-name>" # required
  access_key = "$ACCESS_KEY" # required
  secret_key = "$SECRET_KEY" # required

  output_format = "zip" # zip, dir
`

func (f *S3) SampleConfig() string {
	return sampleConfig
}

func (f *S3) Connect() error {
	return nil
}

func (f *S3) Description() string {
	return "Send grafana json to s3"
}

func (f *S3) Write(file string) error {
	return nil
}

func init() {
	outputs.Add("s3", func() gde.Output {
		return &S3{}
	})
}
