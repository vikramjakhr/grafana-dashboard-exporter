### Input Plugins

This section is for developers who want to create new collection inputs.
GDE is entirely plugin driven. This interface allows for operators to
pick and chose what is gathered and makes it easy for developers
to create new ways of generating metrics.

Plugin authorship is kept as simple as possible to promote people to develop
and submit new inputs.

### Input Plugin Guidelines

- A plugin must conform to the [gde.Input][] interface.
- Input Plugins should call `inputs.Add` in their `init` function to register
  themselves.  See below for a quick example.
- Input Plugins must be added to the
  `github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs/all/all.go` file.
- The `SampleConfig` function should return valid toml that describes how the
  plugin can be configured. This is included in `gde config`.  Please
  consult the [SampleConfig][] page for the latest style
  guidelines.
- The `Description` function should say in one line what this plugin does.
- Follow the recommended [CodeStyle][].

Let's say you've written a plugin that emits metrics about processes on the
current host.

### Input Plugin Example

```go
package simple

// simple.go

import (
    "github.com/vikramjakhr/grafana-dashboard-exporter"
    "github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs"
)

type Simple struct {
    Ok bool
}

func (s *Simple) Description() string {
    return "a demo plugin"
}

func (s *Simple) SampleConfig() string {
    return `
  ## Indicate if everything is fine
  ok = true
`
}

func (s *Simple) Process(acc gde.Accumulator) error {
    // Your logic here

    return nil
}

func init() {
    inputs.Add("simple", func() gde.Input { return &Simple{} })
}
```

[SampleConfig]: https://github.com/vikramjakhr/grafana-dashboard-exporter/wiki/SampleConfig
[CodeStyle]: https://github.com/vikramjakhr/grafana-dashboard-exporter/wiki/CodeStyle
[gde.Input]: https://godoc.org/github.com/vikramjakhr/grafana-dashboard-exporter#Input