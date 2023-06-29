// Code generated by `codegen` DO NOT EDIT.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ExporterMetrics struct {
	Humidity *prometheus.Desc

	Lux *prometheus.Desc

	PM1 *prometheus.Desc

	Pressure *prometheus.Desc

	Proximity *prometheus.Desc

	Temperature *prometheus.Desc
}

func CreateExporterMetricPromDescriptors() *ExporterMetrics {
	return &ExporterMetrics{
		Humidity: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "Humidity"),
			"humidity Relative humidity measured (%)",
			[]string{},
			nil,
		),

		Lux: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "Lux"),
			"current ambient light level (lux)",
			[]string{},
			nil,
		),

		PM1: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "PM1"),
			"Particulate Matter of diameter less than 1 micron. Measured in micrograms per cubic metre (ug/m3)",
			[]string{},
			nil,
		),

		Pressure: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "Pressure"),
			"pressure Pressure measured (hPa)",
			[]string{},
			nil,
		),

		Proximity: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "Proximity"),
			"proximity, with larger numbers being closer proximity and vice versa",
			[]string{},
			nil,
		),

		Temperature: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "Temperature"),
			"temperature Temperature measured (*C)",
			[]string{},
			nil,
		),
	}
}
