package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"

	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"gopkg.in/alecthomas/kingpin.v2"
)

var metricState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "patroni_state", Help: "Current Patroni state"}, []string{"state"})

var metricRole = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "patroni_role", Help: "Current database role"}, []string{"role"})

var metricXlogLocation = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "patroni_xlog_location",
	Help: "Current xlog location (only applicable to masters)"}, []string{"role"})

var metricXlogReceivedLocation = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "patroni_xlog_received_location",
	Help: "Current xlog received location (only applicable to replicas)"}, []string{"role"})

var metricXlogReplayedLocation = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "patroni_xlog_replayed_location",
	Help: "Current xlog replayed location (only applicable to replicas)"}, []string{"role"})

type XlogStatus struct {
	Location         float64 `json:"location"`
	ReceivedLocation float64 `json:"received_location"`
	ReplayedLocation float64 `json:"replayed_location"`
}

type PatroniStatus struct {
	State string     `json:"state"`
	Role  string     `json:"role"`
	Xlog  XlogStatus `json:"xlog"`
}

type promHTTPLogger struct {
	logger log.Logger
}

func (l promHTTPLogger) Println(v ...interface{}) {
	level.Error(l.logger).Log("msg", fmt.Sprint(v...))
}

var POSSIBLE_STATES = []string{"running", "rejecting connections", "not responding", "unknown"}

func setState(status PatroniStatus) {
	for _, state := range POSSIBLE_STATES {
		if status.State == state {
			metricState.WithLabelValues(state).Set(1)
		} else {
			metricState.DeleteLabelValues(state)
		}
	}
}

var POSSIBLE_ROLES = []string{"master", "replica"}

func setRole(status PatroniStatus) {
	for _, role := range POSSIBLE_ROLES {
		if status.Role == role {
			metricRole.WithLabelValues(role).Set(1)
		} else {
			metricRole.DeleteLabelValues(role)
		}
	}
}

func setXlogMetrics(status PatroniStatus) {
	if status.Role == "master" {
		metricXlogLocation.WithLabelValues(status.Role).Set(status.Xlog.Location)
		metricXlogReceivedLocation.DeleteLabelValues("replica")
		metricXlogReplayedLocation.DeleteLabelValues("replica")
	} else {
		metricXlogLocation.DeleteLabelValues("master")
		metricXlogReceivedLocation.WithLabelValues(status.Role).Set(status.Xlog.ReceivedLocation)
		metricXlogReplayedLocation.WithLabelValues(status.Role).Set(status.Xlog.ReplayedLocation)
	}
}

func updateMetrics(httpClient http.Client, url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error connecting", err)
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		fmt.Println("Error getting metrics", getErr)
		return
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("Error reading metrics", readErr)
	}

	status := PatroniStatus{}
	jsonErr := json.Unmarshal(body, &status)
	if jsonErr != nil {
		fmt.Println("Error parsing json", jsonErr)
	}

	if exporterVerbose {
		fmt.Println(string(body))
	}

	setState(status)
	setRole(status)
	setXlogMetrics(status)
}

func updateLoop() {
	httpClient := http.Client{Timeout: time.Second * 2}

	for {
		updateMetrics(httpClient, *patroniServer)

		time.Sleep(time.Duration(5) * time.Second)
	}
}

var (
	webConfig = webflag.AddFlags(kingpin.CommandLine)

	listenAddress   = kingpin.Arg("web.listen-address", "Address to listen on for web interface and telemetry.").Default("0.0.0.0:9394").Envar("PATRONI_EXPORTER_LISTEN_ADDRESS").String()
	metricsPath     = kingpin.Arg("web.metrics-path", "Path under which to expose metrics.").Default("/metrics").Envar("PATRONI_EXPORTER_METRICS_PATH").String()
	patroniServer   = kingpin.Arg("patorni.server", "HTTP API address of a Patroni server (prefix with https:// to connect over HTTPS)").Envar("PATRONI_SERVER_URL").Default("http://localhost:8009").String()
	exporterVerbose = false
)

func main() {
	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)

	kingpin.Flag("verbose", "Enable debug/versbose mode").Default("false").Envar("PATRONI_EXPORTER_VERBOSE").BoolVar(&exporterVerbose)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	prometheus.MustRegister(metricState)
	prometheus.MustRegister(metricRole)
	prometheus.MustRegister(metricXlogLocation)
	prometheus.MustRegister(metricXlogReceivedLocation)
	prometheus.MustRegister(metricXlogReplayedLocation)

	go updateLoop()

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)

	http.Handle(*metricsPath,
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer,
			promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					ErrorLog: &promHTTPLogger{
						logger: logger,
					},
				},
			),
		),
	)

	srv := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
