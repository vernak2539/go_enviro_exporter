//go:generate go run github.com/vernak2539/go_enviro_exporter/cmd/codegen -src ../../.metrics.yml -dest . -pkg metrics

package metrics

//package metrics
//
//import (
//	"github.com/prometheus/client_golang/prometheus"
//)
//
//type metric struct {
//	name string
//	help string
//}
//
//var metrics = []metric{
//	{
//		name: "proximity",
//		help: "proximity, with larger numbers being closer proximity and vice versa",
//	},
//	{
//		name: "lux",
//		help: "current ambient light level (lux)",
//	},
//	{
//		name: "pressure",
//		help: "pressure Pressure measured (hPa)",
//	},
//	{
//		name: "humidity",
//		help: "humidity Relative humidity measured (%)",
//	},
//	{
//		name: "temperature",
//		help: "temperature Temperature measured (*C)",
//	},
//	{
//		name: "PM1",
//		help: "Particulate Matter of diameter less than 1 micron. Measured in micrograms per cubic metre (ug/m3)",
//	},
//}
//
//func BuildPromMetricDescriptions() map[string]*prometheus.Desc {
//	namespace := ""
//	builtMetrics := make(map[string]*prometheus.Desc)
//	for _, m := range metrics {
//		builtMetrics[m.name] = prometheus.NewDesc(
//			prometheus.BuildFQName(namespace, "", m.name),
//			m.help,
//			[]string{},
//			nil,
//		)
//	}
//
//	return builtMetrics
//}
