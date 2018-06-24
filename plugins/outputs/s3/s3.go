package s3

import (
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
)

type S3 struct{}

var sampleConfig = `
  bucket = "<bucket-name>" # required
  access_key = "$ACCESS_KEY" # required
  secret_key = "$SECRET_KEY" # required
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
