package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rubiojr/go-enviroplus/ltr559"
	"github.com/rubiojr/go-enviroplus/pms5003"
	"log"
	"net/http"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
	"periph.io/x/periph/host"
	"reflect"
)

type sensors struct {
	ltr559  *ltr559.LTR559
	bmxx80  *bmxx80.Dev
	pms5003 *pms5003.Device
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
	sensors sensors
}

func newEnvironmentMetricCollector(sensors sensors) *environmentMetricCollector {
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

	humidity := float64(bmxData.Humidity) / float64(physic.PercentRH)
	pressure := float64(bmxData.Pressure) / float64(physic.KiloPascal/10) // convert from nano pascal to hectopascal
	pm := c.sensors.pms5003.LastValue()

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

func initLightProximitySensor() *ltr559.LTR559 {
	sensor, err := ltr559.New()
	if err != nil {
		panic(err)
	}

	return sensor
}

func initBmeSensor() (*bmxx80.Dev, func() error) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}

	sensor, err := bmxx80.NewI2C(b, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}

	return sensor, b.Close
}

func initPmsSensor() *pms5003.Device {
	sensor, err := pms5003.New()
	if err != nil {
		panic(err)
	}

	go func() {
		sensor.StartReading()
	}()

	return sensor
}

func main() {
	flag.Parse()

	ltrSensor := initLightProximitySensor()
	pmsSensor := initPmsSensor()
	bmeSensor, bmxCleanUp := initBmeSensor()
	defer bmxCleanUp()

	sensors := sensors{
		ltr559:  ltrSensor,
		bmxx80:  bmeSensor,
		pms5003: pmsSensor,
	}

	collector := newEnvironmentMetricCollector(sensors)

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
