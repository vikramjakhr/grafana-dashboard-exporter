# File Output Plugin

This plugin stores the Grafana Dashboard and DataSource JSON's to specified directory.

### Configuration:

```
# Send grafana json to specified directory
[[outputs.file]]
  output_dir = "<dir>" # default is /tmp/gde
  output_format = "zip" # zip, dir # default is zip
```