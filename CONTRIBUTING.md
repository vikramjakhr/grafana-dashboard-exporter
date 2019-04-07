### Contributing

1. Open a [new issue][] to discuss the changes you would like to make.  This is
   not strictly required but it may help reduce the amount of rework you need
   to do later.
1. Make changes or write plugin using the guidelines in the following
   documents:
   - [Input Plugins][inputs]
   - [Output Plugins][outputs]
1. Ensure you have added proper unit tests and documentation.
1. Open a new [pull request][].

### GoDoc

Public interfaces for inputs and outputs can be found in the GoDoc:

[![GoDoc](https://godoc.org/vikramjakhr/grafana-dashboard-exporter/gde?status.svg)](https://godoc.org/github.com/vikramjakhr/grafana-dashboard-exporter)

### Common development tasks

**Adding a dependency:**

Assuming you can already build the project, run these in the gde directory:

1. `dep ensure -vendor-only`
2. `dep ensure -add github.com/[dependency]/[new-package]`

**Unit Tests:**

Before opening a pull request you should run the linter checks and
the short tests.

**Run short tests:**

```
make test
```

[new issue]: https://github.com/vikramjakhr/grafana-dashboard-exporter/issues/new/choose
[pull request]: https://github.com/vikramjakhr/grafana-dashboard-exporter/compare
[inputs]: /docs/INPUTS.md
[outputs]: /docs/OUTPUTS.md