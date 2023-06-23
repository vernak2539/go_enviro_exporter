package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rubiojr/go-enviroplus/ltr559"
	"log"
	"net/http"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
	"periph.io/x/periph/host"
)

const namespace = ""

func initBmxSensor() (*bmxx80.Dev, i2c.BusCloser) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}

	d, err := bmxx80.NewI2C(b, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}

	return d, b
}

type sensors struct {
	ltr559 *ltr559.LTR559
	bmxx80 *bmxx80.Dev
	i2cBc  i2c.BusCloser
}

type environmentMetricCollector struct {
	proximity   *prometheus.Desc
	lux         *prometheus.Desc
	pressure    *prometheus.Desc
	humidity    *prometheus.Desc
	temperature *prometheus.Desc
	sensors     sensors
}

func newEnvironmentMetricCollector() *environmentMetricCollector {
	ltr559Sensor, err := ltr559.New()
	if err != nil {
		panic(err)
	}

	bmxx80Sensor, i2cBusCloser := initBmxSensor()

	sensors := sensors{
		ltr559: ltr559Sensor,
		bmxx80: bmxx80Sensor,
		i2cBc:  i2cBusCloser,
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
		sensors: sensors,
	}
}

func (c *environmentMetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.proximity
	ch <- c.lux
	ch <- c.pressure
	ch <- c.humidity
	ch <- c.temperature
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

	bmxData := physic.Env{}
	if err := c.sensors.bmxx80.Sense(&bmxData); err != nil {
		log.Fatal(err)
	}

	// labels added here if needed
	ch <- prometheus.MustNewConstMetric(c.proximity, prometheus.GaugeValue, proximity)
	ch <- prometheus.MustNewConstMetric(c.lux, prometheus.GaugeValue, lux)
	ch <- prometheus.MustNewConstMetric(c.lux, prometheus.GaugeValue, float64(bmxData.Pressure))
	ch <- prometheus.MustNewConstMetric(c.lux, prometheus.GaugeValue, float64(bmxData.Humidity))
	ch <- prometheus.MustNewConstMetric(c.lux, prometheus.GaugeValue, bmxData.Temperature.Celsius())
}

var (
	listenAddress = flag.String("web.listen-address", ":7100", "Address to listen on for web interface.")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

func main() {
	flag.Parse()

	collector := newEnvironmentMetricCollector()
	defer collector.sensors.i2cBc.Close()

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
