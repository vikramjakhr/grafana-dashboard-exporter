# Grafana Input Plugin

This plugin calls the Grafana API's to fetch JSON's of Dashboards and DataSources.

### Configuration:

```
# Fetches Grafana json from specified grafana host
[[inputs.grafana]]
  host = "http://<host>:<port>" # required
  authorization = "Bearer <token>" # required
  dashboard = true # true if dashboard needs to be fetched; default true
  datasource = true # true if datasource needs to be fetched; default true
```