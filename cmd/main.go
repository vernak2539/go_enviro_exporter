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

	pmHistBuckets := map[float64]uint64{0: 0, 5: 5, 10: 10, 15: 15, 20: 20, 25: 25, 30: 30, 35: 35, 40: 40, 45: 45, 50: 50, 55: 55, 60: 60, 65: 65, 70: 70, 75: 75, 80: 80, 85: 85, 90: 90, 95: 95, 100: 100}

	// labels added here if needed
	ch <- prometheus.MustNewConstMetric(c.metrics.Proximity, prometheus.GaugeValue, proximity)
	ch <- prometheus.MustNewConstMetric(c.metrics.Lux, prometheus.GaugeValue, lux)
	ch <- prometheus.MustNewConstMetric(c.metrics.Pressure, prometheus.GaugeValue, pressure)
	ch <- prometheus.MustNewConstMetric(c.metrics.Humidity, prometheus.GaugeValue, humidity)
	ch <- prometheus.MustNewConstMetric(c.metrics.Temperature, prometheus.GaugeValue, temp)
	ch <- prometheus.MustNewConstMetric(c.metrics.Pm1, prometheus.GaugeValue, float64(pm.Pm10Std))
	ch <- prometheus.MustNewConstMetric(c.metrics.Pm25, prometheus.GaugeValue, float64(pm.Pm25Std))
	ch <- prometheus.MustNewConstMetric(c.metrics.Pm10, prometheus.GaugeValue, float64(pm.Pm100Std))
	ch <- prometheus.MustNewConstMetric(c.metrics.Oxidising, prometheus.GaugeValue, gas.Oxidising)
	ch <- prometheus.MustNewConstMetric(c.metrics.Reducing, prometheus.GaugeValue, gas.Reducing)
	ch <- prometheus.MustNewConstMetric(c.metrics.Nh3, prometheus.GaugeValue, gas.NH3)
	ch <- prometheus.MustNewConstHistogram(c.metrics.Pm1_hist, uint64(pm.Pm10Std), float64(pm.Pm10Std), pmHistBuckets)
	ch <- prometheus.MustNewConstHistogram(c.metrics.Pm25_hist, uint64(pm.Pm25Std), float64(pm.Pm25Std) - float64(pm.Pm10Std), pmHistBuckets)
	ch <- prometheus.MustNewConstHistogram(c.metrics.Pm10_hist, uint64(pm.Pm100Std), float64(pm.Pm100Std) - float64(pm.Pm25Std), pmHistBuckets)
}

var (
	listenAddress = flag.String("web.listen-address", ":7100", "Address to listen on for web interface.")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

func main() {
	flag.Parse()

	ltrSensor := sensors.NewLtrSensor()
	pmsSensor := sensors.NewPmsSensor()

	micsSensor := sensors.NewMicsSensor()
	defer micsSensor.Close()

	bmeSensor := sensors.NewBmeSensor()
	defer bmeSensor.Close()

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
