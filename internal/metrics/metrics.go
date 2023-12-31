// Code generated by `codegen` DO NOT EDIT.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ExporterMetrics struct {
	Humidity *prometheus.Desc

	Lux *prometheus.Desc

	Nh3 *prometheus.Desc

	Oxidising *prometheus.Desc

	Pm1 *prometheus.Desc

	Pm10 *prometheus.Desc

	Pm10_hist *prometheus.HistogramVec

	Pm1_hist *prometheus.HistogramVec

	Pm25 *prometheus.Desc

	Pm25_hist *prometheus.HistogramVec

	Pressure *prometheus.Desc

	Proximity *prometheus.Desc

	Reducing *prometheus.Desc

	Temperature *prometheus.Desc
}

func CreateExporterMetricPromDescriptors() *ExporterMetrics {
	return &ExporterMetrics{

		Humidity: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "humidity"),
			"humidity Relative humidity measured (%)",
			[]string{},
			nil,
		),

		Lux: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "lux"),
			"current ambient light level (lux)",
			[]string{},
			nil,
		),

		Nh3: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "NH3"),
			"mostly Ammonia but could also include Hydrogen, Ethanol, Propane, Iso-butane (Ohms)",
			[]string{},
			nil,
		),

		Oxidising: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "oxidising"),
			"Mostly nitrogen dioxide but could include NO and Hydrogen (Ohms)",
			[]string{},
			nil,
		),

		Pm1: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "PM1"),
			"Particulate Matter of diameter less than 1 micron. Measured in micrograms per cubic metre (ug/m3)",
			[]string{},
			nil,
		),

		Pm10: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "PM10"),
			"Particulate Matter of diameter less than 10 microns. Measured in micrograms per cubic metre (ug/m3)",
			[]string{},
			nil,
		),

		Pm10_hist: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "pm10_measurements",
				Help:        "Histogram of Particulate Matter of diameter less than 10 microns measurements",
				ConstLabels: nil,
				Buckets:     []float64{0.0, 5.0, 10.0, 15.0, 20.0, 25.0, 30.0, 35.0, 40.0, 45.0, 50.0, 55.0, 60.0, 65.0, 70.0, 75.0, 80.0, 85.0, 90.0, 95.0, 100.0},
			},
			[]string{},
		),

		Pm1_hist: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "pm1_measurements",
				Help:        "Histogram of Particulate Matter of diameter less than 1 micron measurements",
				ConstLabels: nil,
				Buckets:     []float64{0.0, 5.0, 10.0, 15.0, 20.0, 25.0, 30.0, 35.0, 40.0, 45.0, 50.0, 55.0, 60.0, 65.0, 70.0, 75.0, 80.0, 85.0, 90.0, 95.0, 100.0},
			},
			[]string{},
		),

		Pm25: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "PM25"),
			"Particulate Matter of diameter less than 2.5 microns. Measured in micrograms per cubic metre (ug/m3)",
			[]string{},
			nil,
		),

		Pm25_hist: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "pm25_measurements",
				Help:        "Histogram of Particulate Matter of diameter less than 2.5 microns measurements",
				ConstLabels: nil,
				Buckets:     []float64{0.0, 5.0, 10.0, 15.0, 20.0, 25.0, 30.0, 35.0, 40.0, 45.0, 50.0, 55.0, 60.0, 65.0, 70.0, 75.0, 80.0, 85.0, 90.0, 95.0, 100.0},
			},
			[]string{},
		),

		Pressure: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "pressure"),
			"pressure Pressure measured (hPa)",
			[]string{},
			nil,
		),

		Proximity: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "proximity"),
			"proximity, with larger numbers being closer proximity and vice versa",
			[]string{},
			nil,
		),

		Reducing: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "reducing"),
			"Mostly carbon monoxide but could include H2S, Ammonia, Ethanol, Hydrogen, Methane, Propane, Iso-butane (Ohms)",
			[]string{},
			nil,
		),

		Temperature: prometheus.NewDesc(
			prometheus.BuildFQName("", "", "temperature"),
			"temperature Temperature measured (*C)",
			[]string{},
			nil,
		),
	}
}

func (m *ExporterMetrics) Describe(ch chan<- *prometheus.Desc) {

	ch <- m.Humidity

	ch <- m.Lux

	ch <- m.Nh3

	ch <- m.Oxidising

	ch <- m.Pm1

	ch <- m.Pm10

	m.Pm10_hist.Describe(ch)

	m.Pm1_hist.Describe(ch)

	ch <- m.Pm25

	m.Pm25_hist.Describe(ch)

	ch <- m.Pressure

	ch <- m.Proximity

	ch <- m.Reducing

	ch <- m.Temperature

}
