package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rubiojr/go-enviroplus/ltr559"
	"log"
	"net/http"
)

const namespace = "enviro"

type sensors struct {
	ltr559 *ltr559.LTR559
}

type environmentMetricCollector struct {
	proximity *prometheus.Desc
	lux       *prometheus.Desc
	sensors   sensors
}

func newEnvironmentMetricCollector() *environmentMetricCollector {
	ltr559Sensor, err := ltr559.New()
	if err != nil {
		panic(err)
	}

	sensors := sensors{
		ltr559: ltr559Sensor,
	}

	return &environmentMetricCollector{
		proximity: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "proximity"),
			"proximity, with larger numbers being closer proximity and vice versa",
			[]string{},
			nil,
		),
		lux: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "lux"),
			"current ambient light level (lux)",
			[]string{}, // labels added here if needed
			nil,
		),
		sensors: sensors,
	}
}

func (c *environmentMetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.proximity
	ch <- c.lux
}

func (c *environmentMetricCollector) Collect(ch chan<- prometheus.Metric) {
	proximity, err := c.sensors.ltr559.Proximity()
	if err != nil {
		panic(err)
	}

	lux, err := c.sensors.ltr559.Lux()
	if err != nil {
		panic(err)
	}

	// labels added here if needed
	ch <- prometheus.MustNewConstMetric(c.proximity, prometheus.GaugeValue, proximity)
	ch <- prometheus.MustNewConstMetric(c.lux, prometheus.GaugeValue, lux)
}

var (
	listenAddress = flag.String("web.listen-address", ":7100", "Address to listen on for web interface.")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

func main() {
	flag.Parse()

	collector := newEnvironmentMetricCollector()
	prometheus.MustRegister(collector)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Enviro Exporter Metrics</title></head>
             <body>
             <h1>Enviro Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	fmt.Printf("listening at http://localhost%s%s", *listenAddress, *metricsPath)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
