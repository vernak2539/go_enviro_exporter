package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vernak2539/go_enviro_exporter/internal/sensors"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

type allSensors struct {
	ltr559  *sensors.LtrSensor
	pms5003 *sensors.PmsSensor
	bmxx80  *bmxx80.Dev
}

type metrics struct {
	proximity   *prometheus.Desc
	lux         *prometheus.Desc
	pressure    *prometheus.Desc
	humidity    *prometheus.Desc
	temperature *prometheus.Desc
	pm1         *prometheus.Desc
}

type environmentMetricCollector struct {
	metrics metrics
	sensors allSensors
}

func newEnvironmentMetricCollector(sensors allSensors) *environmentMetricCollector {
	const namespace = ""

	return &environmentMetricCollector{
		sensors: sensors,
		metrics: metrics{
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
			pressure: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "pressure"),
				"pressure Pressure measured (hPa)",
				[]string{}, // labels added here if needed
				nil,
			),
			humidity: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "humidity"),
				"humidity Relative humidity measured (%)",
				[]string{}, // labels added here if needed
				nil,
			),
			temperature: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "temperature"),
				"temperature Temperature measured (*C)",
				[]string{}, // labels added here if needed
				nil,
			),
			pm1: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "PM1"),
				"Particulate Matter of diameter less than 1 micron. Measured in micrograms per cubic metre (ug/m3)",
				[]string{}, // labels added here if needed
				nil,
			),
		},
	}
}

func (c *environmentMetricCollector) Describe(ch chan<- *prometheus.Desc) {
	v := reflect.ValueOf(c.metrics)
	values := make([]*prometheus.Desc, v.NumField())

	for _, metric := range values {
		ch <- metric
	}
}

func (c *environmentMetricCollector) Collect(ch chan<- prometheus.Metric) {
	proximity := c.sensors.ltr559.GetProximity()
	lux := c.sensors.ltr559.GetLux()
	pm := c.sensors.pms5003.GetPmMeasurement()

	bmxData := physic.Env{}
	if err := c.sensors.bmxx80.Sense(&bmxData); err != nil {
		log.Fatal(err)
	}

	humidity := float64(bmxData.Humidity) / float64(physic.PercentRH)
	pressure := float64(bmxData.Pressure) / float64(physic.KiloPascal/10) // convert from nano pascal to hectopascal

	// labels added here if needed
	ch <- prometheus.MustNewConstMetric(c.metrics.proximity, prometheus.GaugeValue, proximity)
	ch <- prometheus.MustNewConstMetric(c.metrics.lux, prometheus.GaugeValue, lux)
	ch <- prometheus.MustNewConstMetric(c.metrics.pressure, prometheus.GaugeValue, pressure)
	ch <- prometheus.MustNewConstMetric(c.metrics.humidity, prometheus.GaugeValue, humidity)
	ch <- prometheus.MustNewConstMetric(c.metrics.temperature, prometheus.GaugeValue, bmxData.Temperature.Celsius())
	ch <- prometheus.MustNewConstMetric(c.metrics.pm1, prometheus.GaugeValue, float64(pm.Pm10Std))
}

var (
	listenAddress = flag.String("web.listen-address", ":7100", "Address to listen on for web interface.")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

func initBmeSensor() (*bmxx80.Dev, func() error, func() error) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	sensor, err := bmxx80.NewI2C(bus, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}

	return sensor, bus.Close, sensor.Halt
}

func main() {
	flag.Parse()

	ltrSensor := sensors.NewLtrSensor()
	pmsSensor := sensors.NewPmsSensor()
	bmeSensor, bmxBusCleanUp, bmxSensorCleanUp := initBmeSensor()
	defer bmxBusCleanUp()
	defer bmxSensorCleanUp()

	s := allSensors{
		ltr559:  ltrSensor,
		bmxx80:  bmeSensor,
		pms5003: pmsSensor,
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
