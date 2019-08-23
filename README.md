Small utility that pings a list or domain names and serves a prometheus endpoint with the results. 
Prometheus metrics are at `:8080/metrics`

##### Example Configuration:
```hcl
count = 5
timeout = "10s"
interval = "60s"
addresses = [
  "google.com",
  "facebook.com",
  "github.com",
]
```

##### Options:
```console
  -config string
        Path to edge pinger configuration file
  -debug
        Debug (default: false)
  -port int
        Port to listen for Prometheus requests (default 8080)
```

## How to run:
```console
$ go run main.go --config example/edge-pinger-configuration.hcl 
2019/08/22 17:39:27 Initiating a ping loop for google.com Count=5 Timeout=10s Interval=1m0s
2019/08/22 17:39:27 Initiating a ping loop for facebook.com Count=5 Timeout=10s Interval=1m0s
2019/08/22 17:39:27 Initiating a ping loop for github.com Count=5 Timeout=10s Interval=1m0s
2019/08/22 17:39:27 Listening on :8080
```