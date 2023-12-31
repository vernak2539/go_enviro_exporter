package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vernak2539/go_enviro_exporter/internal/metrics"
	"github.com/vernak2539/go_enviro_exporter/internal/sensors"
)

type allSensors struct {
	ltr559   *sensors.LtrSensor
	pms5003  *sensors.PmsSensor
	bmxx80   *sensors.BmeSensor
	mics6814 *sensors.MicsSensor
}

type environmentMetricCollector struct {
	metrics *metrics.ExporterMetrics
	sensors allSensors
}

func newEnvironmentMetricCollector(sensors allSensors) *environmentMetricCollector {
	const namespace = ""

	return &environmentMetricCollector{
		sensors: sensors,
		metrics: metrics.CreateExporterMetricPromDescriptors(),
	}
}

func (c *environmentMetricCollector) Describe(ch chan<- *prometheus.Desc) {
	c.metrics.Describe(ch)
}

func (c *environmentMetricCollector) Collect(ch chan<- prometheus.Metric) {
	proximity := c.sensors.ltr559.GetProximity()
	lux := c.sensors.ltr559.GetLux()
	pm := c.sensors.pms5003.GetPmMeasurement()
	humidity := c.sensors.bmxx80.GetHumidity()
	pressure := c.sensors.bmxx80.GetPressure()
	temp := c.sensors.bmxx80.GetTemperature()
	gas := c.sensors.mics6814.GetGasMeasurements()

	pm1Std64 := float64(pm.Pm10Std)
	pm25Std64 := float64(pm.Pm25Std)
	pm100Std64 := float64(pm.Pm100Std)

	// labels added here if needed
	ch <- prometheus.MustNewConstMetric(c.metrics.Proximity, prometheus.GaugeValue, proximity)
	ch <- prometheus.MustNewConstMetric(c.metrics.Lux, prometheus.GaugeValue, lux)
	ch <- prometheus.MustNewConstMetric(c.metrics.Pressure, prometheus.GaugeValue, pressure)
	ch <- prometheus.MustNewConstMetric(c.metrics.Humidity, prometheus.GaugeValue, humidity)
	ch <- prometheus.MustNewConstMetric(c.metrics.Temperature, prometheus.GaugeValue, temp)
	ch <- prometheus.MustNewConstMetric(c.metrics.Pm1, prometheus.GaugeValue, pm1Std64)
	ch <- prometheus.MustNewConstMetric(c.metrics.Pm25, prometheus.GaugeValue, pm25Std64)
	ch <- prometheus.MustNewConstMetric(c.metrics.Pm10, prometheus.GaugeValue, pm100Std64)
	ch <- prometheus.MustNewConstMetric(c.metrics.Oxidising, prometheus.GaugeValue, gas.Oxidising)
	ch <- prometheus.MustNewConstMetric(c.metrics.Reducing, prometheus.GaugeValue, gas.Reducing)
	ch <- prometheus.MustNewConstMetric(c.metrics.Nh3, prometheus.GaugeValue, gas.NH3)

	c.metrics.Pm1_hist.WithLabelValues().Observe(pm1Std64)
	c.metrics.Pm1_hist.Collect(ch)

	c.metrics.Pm25_hist.WithLabelValues().Observe(pm25Std64)
	c.metrics.Pm25_hist.Collect(ch)

	c.metrics.Pm10_hist.WithLabelValues().Observe(pm100Std64)
	c.metrics.Pm10_hist.Collect(ch)
}

var (
	listenAddress = flag.String("web.listen-address", ":7100", "Address to listen on for web interface.")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

func main() {
	flag.Parse()

	ltrSensor := sensors.NewLtrSensor()
	pmsSensor := sensors.NewPmsSensor()
	bmeSensor := sensors.NewBmeSensor()

	micsSensor := sensors.NewMicsSensor()
	defer micsSensor.Close()

	s := allSensors{
		ltr559:   ltrSensor,
		bmxx80:   bmeSensor,
		pms5003:  pmsSensor,
		mics6814: micsSensor,
	}

	collector := newEnvironmentMetricCollector(s)

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
