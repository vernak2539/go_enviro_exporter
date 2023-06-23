package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rubiojr/go-enviroplus/ltr559"
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
		proximity: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "proximity"),
			"Proximity metric",
			[]string{"name", "path"}, nil,
		),
		lux: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "lux"),
			"Lux metric",
			[]string{"name", "path"}, nil,
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

	ch <- prometheus.MustNewConstMetric(c.proximity, prometheus.GaugeValue, proximity)
	ch <- prometheus.MustNewConstMetric(c.lux, prometheus.GaugeValue, lux)
}

var metricsPath = "/metrics"

func main() {
	collector := newEnvironmentMetricCollector()
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Mirth Channel Exporter</title></head>
             <body>
             <h1>Mirth Channel Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
}
