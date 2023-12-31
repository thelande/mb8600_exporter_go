/*
Copyright 2023 Thomas Helander

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"net/http"
	"os"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/thelande/mb8600/pkg/mb8600"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const (
	exporterName  = "mb8600_exporter"
	exporterTitle = "Go Exporter Template"
)

var (
	address = kingpin.Flag(
		"address",
		"Address of the cable modem.",
	).Default("192.168.100.1").IP()
	username = kingpin.Flag(
		"username",
		"Username to use to login to the modem.",
	).Default("admin").Envar("MODEM_USERNAME").String()
	password = kingpin.Flag(
		"password",
		"Password to use to login to the modem.",
	).Default("motorola").Envar("MODEM_PASSWORD").String()
	metricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	webConfig = webflag.AddFlags(kingpin.CommandLine, ":9813")
	logger    log.Logger
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version.Print(exporterName))
	kingpin.Parse()

	logger = promlog.New(promlogConfig)
	level.Info(logger).Log("msg", fmt.Sprintf("Starting %s", exporterName), "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	client := mb8600.NewMotoClient(address.String(), *username, *password, logger)
	if _, err := client.Login(); err != nil {
		level.Error(logger).Log("msg", "failed to login", "err", err)
		os.Exit(1)
	}

	collector := NewCollector(address.String(), *username, *password, logger)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	landingConfig := web.LandingConfig{
		Name:        exporterTitle,
		Description: "Prometheus go-based Exporter",
		Version:     version.Info(),
		Links: []web.LandingLinks{
			{
				Address: *metricsPath,
				Text:    "Metrics",
			},
		},
	}
	landingPage, err := web.NewLandingPage(landingConfig)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.Handle("/", landingPage)

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "HTTP listener stopped", "error", err)
		os.Exit(1)
	}
}
