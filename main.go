package main

import (
	"flag"
	"fmt"
	"github.com/thoas/go-funk"
	"log"
	"net/http"
	"time"

	"github.com/latchmihay/edge-pinger/pkg/config"
	"github.com/latchmihay/edge-pinger/pkg/engine"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port       = flag.Int("port", 8080, "Port to listen for Prometheus requests")
	configFile = flag.String("config", "", "Path to edge pinger configuration file")
	debug      = flag.Bool("debug", false, "Debug (default: false)")
)

const (
	defaultCount    = 5
	defaultTimeout  = "15s"
	defaultInterval = "1m"
)

func main() {
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())

	if funk.IsEmpty(*configFile) {
		log.Fatal("Please pass in a config file via --config")
	}

	hclConfig, err := config.LoadConfigFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	edgePingerConfig, err := config.ParseConfig(hclConfig)
	if err != nil {
		log.Fatal(err)
	}

	if edgePingerConfig.Count == 0 {
		edgePingerConfig.Count = defaultCount
	}

	if edgePingerConfig.Timeout == "" {
		edgePingerConfig.Timeout = defaultTimeout
	}
	timeout, err := time.ParseDuration(edgePingerConfig.Timeout)
	if err != nil {
		log.Fatal(err)
	}

	if edgePingerConfig.Interval == "" {
		edgePingerConfig.Interval = defaultInterval
	}
	interval, err := time.ParseDuration(edgePingerConfig.Interval)
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range edgePingerConfig.Addresses {
		client, err := engine.NewPing(addr, edgePingerConfig.Count, timeout, *debug)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Initiating a ping loop for %v Count=%v Timeout=%v Interval=%v", addr, edgePingerConfig.Count, timeout, interval)
		go func() {
			for true {
				client.Run()
				time.Sleep(interval)
			}
		}()
	}
	log.Printf("Listening on :%d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

}
