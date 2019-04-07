# S3 Output Plugin

This plugin stores the Grafana Dashboard and DataSource JSON's to specified AWS S3 Bucket.

### Configuration:

```
# Send grafana json to s3
[[outputs.s3]]
  bucket = "<bucket-name>" # required
  access_key = "$ACCESS_KEY" # required
  secret_key = "$SECRET_KEY" # required
  region = "s3-bucket-region"
  bucketPrefix = "<prefix>"
  output_format = "zip" # zip, dir
```