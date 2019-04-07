# GDE (Grafana Dashboard Exporter) [![CircleCI](https://circleci.com/gh/vikramjakhr/grafana-dashboard-exporter/tree/master.svg?style=svg)](https://circleci.com/gh/vikramjakhr/grafana-dashboard-exporter/tree/master)

GDE is an extremely powerful open source agent for backing up grafana dashboards. It's based on the influxdata telegraf theme.

Design goals are to have a minimal memory footprint with a plugin system so
that developers in the community can easily add support for collecting
metrics.

GDE is plugin-driven and has the concept of 2 distinct plugin types:

1. [Input Plugins](#input-plugins) collect grafana dashboards json from the grafana server
2. [Output Plugins](#output-plugins) write metrics to various destinations

New plugins are designed to be easy to contribute, we'll eagerly accept pull
requests and will manage the set of plugins that GDE supports.

## Contributing

There are many ways to contribute:
- Fix and [report bugs](https://github.com/vikramjakhr/grafana-dashboard-exporter/issues/new)
- [Improve documentation](https://github.com/vikramjakhr/grafana-dashboard-exporter/issues?q=is%3Aopen+label%3Adocumentation)
- [Review code and feature proposals](https://github.com/vikramjakhr/grafana-dashboard-exporter/pulls)
- [Contribute plugins](CONTRIBUTING.md)

## Installation:

You can download the binaries directly from 
the [releases](https://github.com/vikramjakhr/grafana-dashboard-exporter/releases) section.

### Ansible Role:

Ansible role: In progress :) 

### From Source:

GDE requires golang version 1.9 or newer, the Makefile requires GNU make.

1. [Install Go](https://golang.org/doc/install) >=1.9 (1.11 recommended)
2. [Install dep](https://golang.github.io/dep/docs/installation.html) ==v0.5.0
3. Download Telegraf source:
   ```
   go get -d github.com/vikramjakhr/grafana-dashboard-exporter
   ```
4. Run make from the source directory
   ```
   cd "$HOME/go/src/github.com/vikramjakhr/grafana-dashboard-exporter"
   make
   ```
   
### Changelog

View the [changelog](/CHANGELOG.md) for the latest updates and changes by
version.

## How to use it:

See usage with:

```
gde --help
```

#### Generate a gde config file:

```
gde config > gde.conf
```

#### Generate config with only grafana input & S3 output plugins defined:

```
gde --input-filter grafana --output-filter s3 config
```

#### Run a single gde collection:

```
gde --config gde.conf --test
```

#### Run gde with all plugins defined in config file:

```
gde --config gde.conf
```

## Input Plugins

* [grafana](./plugins/inputs/grafana)

## Output Plugins

* [file](./plugins/outputs/file)
* [s3](./plugins/outputs/s3)