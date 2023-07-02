package main

import (
	"bytes"
	_ "embed"
	"flag"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type rawMetric struct {
	Metric         string   `json:"metric"`
	Help           string   `json:"help"`
	Type           string   `json:"type"`
	Namespace      string   `json:"namespace"`
	Buckets        string   `json:"buckets"`
	ConstLabels    []string `json:"const_labels"`
	VariableLabels []string `json:"variable_labels"`
}

type metricData struct {
	Name           string
	Metric         string
	Type           string
	Namespace      string
	Buckets        string
	Help           string
	ConstLabels    []string
	VariableLabels []string
}

type args struct {
	src  string
	dest string
	pkg  string
}

func main() {
	log.SetFlags(0)
	var args args
	flag.StringVar(&args.src, "src", "./.metrics.yml", "path to .metrics.yml")
	flag.StringVar(&args.dest, "dest", "./internal/metrics", "package destination")
	flag.StringVar(&args.pkg, "pkg", "metrics", "package name")
	flag.Parse()

	data, err := ioutil.ReadFile(args.src)
	if err != nil {
		log.Fatalln("unable to open file:", err)
	}

	var rawMetrics map[string]rawMetric
	if err := yaml.Unmarshal(data, &rawMetrics); err != nil {
		log.Fatalln("unable to unmarshal yml:", err)
	}

	metrics := make([]metricData, 0, len(rawMetrics))
	for entry, m := range rawMetrics {
		metrics = append(metrics, metricData{
			Name:           strings.Title(entry),
			Metric:         m.Metric,
			Type:           m.Type,
			Namespace:      m.Namespace,
			Buckets:        m.Buckets,
			Help:           m.Help,
			ConstLabels:    m.ConstLabels,
			VariableLabels: m.VariableLabels,
		})
	}

	sort.Slice(metrics, func(i, j int) bool {
		return strings.ToLower(metrics[i].Name) < strings.ToLower(metrics[j].Name)
	})

	tm := template.Must(template.New("metrics.go").Parse(metricsTemplate))
	vars := struct {
		Metrics     []metricData
		PackageName string
		Src         string
	}{
		Metrics:     metrics,
		PackageName: args.pkg,
		Src:         args.src,
	}
	var b bytes.Buffer
	must(tm.Execute(&b, vars))

	err = os.MkdirAll(args.dest, os.ModePerm)
	must(err)

	destpath := filepath.Join(args.dest, "metrics.go")
	writeGoFile(destpath, &b, false)
}

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func writeGoFile(path string, b *bytes.Buffer, verify bool) {
	formatted, err := format.Source(b.Bytes())
	must(err)
	if verify {
		destdata, err := os.ReadFile(path)
		must(err)
		if !reflect.DeepEqual(formatted, destdata) {
			log.Fatalf("metrics.go is a generated file by this script.\n")
			log.Fatalf("To add a rawMetric, add it to .metrics.yml and re-run this script.\n")
			log.Fatalf("%s is out of sync.\n", path)
			os.Exit(1)
		}
	}

	must(os.WriteFile(path, formatted, 0644))
}

//go:embed metrics.tmpl
var metricsTemplate string
