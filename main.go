package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CLI struct {
	Teams       []string      `help:"All of the teams to emit metrics for" env:"TEAMS"`
	MaxSize     uint          `help:"The largest that the metric can grow" default:"100000"`
	MaxIncrease int           `help:"The highest that the metric can be increased by" default:"10000"`
	Interval    time.Duration `help:"The interval at which the metric increases" default:"5s"`
	ServerName  string        `help:"The name of the server" optional:""`
}

var (
	teamVals map[string]uint64 = map[string]uint64{}

	processMemoryConsumptionGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "process_memory_total_bytes",
		Help: "The total amount of memory being used by this process",
	}, []string{"server", "team"})
)

func recordMetrics(cli *CLI) {
	for {
		for _, team := range cli.Teams {
			toIncrease := rand.Intn(cli.MaxIncrease)
			if teamVals[team]+uint64(toIncrease) > uint64(cli.MaxSize) {
				teamVals[team] = 0
			}

			teamVals[team] += uint64(toIncrease)
			processMemoryConsumptionGauge.WithLabelValues(cli.ServerName, team).Set(float64(teamVals[team]))
		}

		time.Sleep(cli.Interval)
	}
}

func main() {
	cli := CLI{}
	_ = kong.Parse(&cli)

	if cli.ServerName == "" {
		name, err := os.Hostname()
		if err != nil {
			log.Fatalln(err)
		}
		cli.ServerName = name
	}

	for _, team := range cli.Teams {
		teamVals[team] = 0
	}

	go recordMetrics(&cli)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":23456", nil)
}
