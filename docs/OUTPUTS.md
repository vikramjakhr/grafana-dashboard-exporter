### Output Plugins

This section is for developers who want to create a new output sink. Outputs
are created in a similar manner as collection plugins, and their interface has
similar constructs.

### Output Plugin Guidelines

- An output must conform to the [gde.Output][] interface.
- Outputs should call `outputs.Add` in their `init` function to register
  themselves.  See below for a quick example.
- To be available within Telegraf itself, plugins must add themselves to the
  `github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs/all/all.go` file.
- The `SampleConfig` function should return valid toml that describes how the
  plugin can be configured. This is included in `gde config`.  Please
  consult the [SampleConfig][] page for the latest style guidelines.
- The `Description` function should say in one line what this output does.
- Follow the recommended [CodeStyle][].

### Output Plugin Example

```go
package simpleoutput

// simpleoutput.go

import (
    "github.com/vikramjakhr/grafana-dashboard-exporter"
    "github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
)

type Simple struct {
    Ok bool
}

func (s *Simple) Description() string {
    return "a demo output"
}

func (s *Simple) SampleConfig() string {
    return `
  ok = true
`
}

func (s *Simple) Connect() error {
    // Make a connection to the URL here
    return nil
}

func (s *Simple) Write(metrics []gde.Metric) error {
    for _, metric := range metrics {
        // write `metric` to the output sink here
    }
    return nil
}

func init() {
    outputs.Add("simpleoutput", func() gde.Output { return &Simple{} })
}

```

[SampleConfig]: https://github.com/vikramjakhr/grafana-dashboard-exporter/wiki/SampleConfig
[CodeStyle]: https://github.com/vikramjakhr/grafana-dashboard-exporter/wiki/CodeStyle
[gde.Output]: https://godoc.org/github.com/vikramjakhr/grafana-dashboard-exporter#Output