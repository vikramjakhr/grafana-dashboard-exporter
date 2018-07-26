# GDE (Grafana Dashboard Exporter)

GDE is an extremely powerful solution for backing up grafana dashboards. It's based on the influxdata telegraf theme.

# Installation Procedure of GDE Agent
##### Step 1: Download the [binary](https://github.com/vikramjakhr/grafana-dashboard-exporter/releases/download/beta-v1.0.0/gde). Example command below.
```
wget https://github.com/vikramjakhr/grafana-dashboard-exporter/releases/download/beta-v1.0.0/gde
```
# Sample GDE config (/etc/gde/gde.conf)
```
[agent]
  interval = "5m"
  round_interval = true
  debug = true
  quiet = false
  logfile = "/var/log/gde/gde.log"

[[outputs.file]]
  output_dir = "<dir>" # required
  output_format = "zip" # zip, dir

[[outputs.s3]]
  bucket = "<bucket-name>" # required
  access_key = ""
  secret_key = ""
  region = ""
  bucket_prefix = ""
  output_format = "zip" # zip, dir

[[inputs.grafana]]
  host = "http://<grafana-host>"
  authorization = "Bearer <token>"
  datasource = true
  dashboard = true
  org = true
```