package engine

import (
	"log"
	"time"

	"github.com/latchmihay/edge-pinger/pkg/prom"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sparrc/go-ping"
)

type PingableAddress struct {
	Hostname string
	Count    int
	Timeout  time.Duration
	debug    bool
	gauges   map[string]*prometheus.GaugeVec
	counters map[string]*prometheus.CounterVec
}

var (
	addrProm = "edge_pinger"
	// label for the collectors
	labels = []string{"ip", "hostname"}
	// create collectors
	avgRtt      = prom.AddGauge(addrProm, "rtt_avg", "average round trip time", labels)
	minRtt      = prom.AddGauge(addrProm, "rtt_min", "min round trip time", labels)
	maxRtt      = prom.AddGauge(addrProm, "rtt_max", "max round trip time", labels)
	stdDevRtt   = prom.AddGauge(addrProm, "rtt_stddev", "max round trip time", labels)
	packetLoss  = prom.AddGauge(addrProm, "packet_loss", "percentage of packets lost", labels)
	packetSent  = prom.AddCounter(addrProm, "packets_transmitted", "total number of packets sent", labels)
	packetRecv  = prom.AddCounter(addrProm, "packet_received", "total number of packets received", labels)
	numberPings = prom.AddCounter(addrProm, "total_packets", "total number of packets sent", labels)
	healthy     = prom.AddGauge(addrProm, "healthy", "Failed to resolve and to ping", labels)
)

func Init() {
	prometheus.MustRegister(avgRtt, minRtt, maxRtt, stdDevRtt, packetLoss, packetSent, packetRecv, numberPings, healthy)
}

func NewPing(addr string, count int, timeout time.Duration, debug bool) *PingableAddress {
	return &PingableAddress{
		Hostname: addr,
		Count:    count,
		Timeout:  timeout,
		debug:    debug,
		gauges:   map[string]*prometheus.GaugeVec{"avgRtt": avgRtt, "minRtt": minRtt, "maxRtt": maxRtt, "stdDevRtt": stdDevRtt, "packetLoss": packetLoss, "health": healthy},
		counters: map[string]*prometheus.CounterVec{"packetSent": packetSent, "packetRecv": packetRecv, "numberPings": numberPings},
	}
}

func (pa *PingableAddress) Run() {
	if pa.debug {
		log.Printf("Attempting to ping %v", pa.Hostname)
	}
	pinger, err := ping.NewPinger(pa.Hostname)
	if err != nil {
		pa.gauges["health"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Set(1)
		return
	}
	pa.gauges["health"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Set(0)

	pinger.Count = pa.Count
	pinger.Timeout = pa.Timeout
	pinger.OnRecv = func(pkt *ping.Packet) {
		pa.counters["numberPings"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Add(1)
		if pa.debug {
			log.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}

	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		pa.gauges["packetLoss"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Set(stats.PacketLoss)

		pa.gauges["avgRtt"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Set(stats.AvgRtt.Seconds() * 1e3)
		pa.gauges["minRtt"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Set(stats.MinRtt.Seconds() * 1e3)
		pa.gauges["maxRtt"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Set(stats.MaxRtt.Seconds() * 1e3)
		pa.gauges["stdDevRtt"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Set(stats.StdDevRtt.Seconds() * 1e3)

		pa.counters["packetSent"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Add(float64(stats.PacketsSent))
		pa.counters["packetRecv"].WithLabelValues(pinger.IPAddr().String(), pa.Hostname).Add(float64(stats.PacketsRecv))

		if pa.debug {
			log.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			log.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			log.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}
	}

	pinger.Run()
}
