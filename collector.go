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

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thelande/mb8600_exporter/mb8600"
)

const namespace = "mb8600"

var (
	labels     = []string{"channel"}
	infoLabels = []string{
		"channel",
		"channel_id",
		"modulation",
	}

	upDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the modem reachable?",
		nil,
		nil,
	)

	downstreamInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "info"),
		"Additional information about the channel.",
		infoLabels,
		nil,
	)

	downstreamFrequencyDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "frequency_mhz"),
		"Frequency of the channel.",
		labels,
		nil,
	)

	downstreamPowerDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "power_dbmv"),
		"Power of the channel.",
		labels,
		nil,
	)

	downstreamSignalToNoiseDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "snr_db"),
		"Signal-to-noise of the channel.",
		labels,
		nil,
	)

	downstreamCorrectedErrorsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "corrected_errors"),
		"Corrected error count.",
		labels,
		nil,
	)

	downstreamUncorrectedErrorsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "uncorrected_errors"),
		"Uncorrected error count.",
		labels,
		nil,
	)

	downstreamLockStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "lock_state"),
		"The lock state of the channel (UNLOCKED = 0, LOCKED = 1)",
		labels,
		nil,
	)

	upstreamInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "info"),
		"Additional information about the channel.",
		infoLabels,
		nil,
	)

	upstreamFrequencyDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "frequency_mhz"),
		"Frequency of the channel.",
		labels,
		nil,
	)

	upstreamPowerDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "power_dbmv"),
		"Power of the channel.",
		labels,
		nil,
	)

	upstreamSymbolRateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "symbol_rate_ksyms"),
		"Symbol rate of the channel",
		labels,
		nil,
	)

	upstreamLockStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "lock_state"),
		"The lock state of the channel (UNLOCKED = 0, LOCKED = 1)",
		labels,
		nil,
	)

	downstreamMetrics = []*prometheus.Desc{
		downstreamInfo,
		downstreamFrequencyDesc,
		downstreamPowerDesc,
		downstreamSignalToNoiseDesc,
		downstreamCorrectedErrorsDesc,
		downstreamUncorrectedErrorsDesc,
		downstreamLockStatus,
	}

	upstreamMetrics = []*prometheus.Desc{
		upstreamInfo,
		upstreamFrequencyDesc,
		upstreamPowerDesc,
		upstreamSymbolRateDesc,
		upstreamLockStatus,
	}
)

type Collector struct {
	Client *mb8600.MotoClient
}

func NewCollector(address, username, password string, logger log.Logger) *Collector {
	return &Collector{
		Client: mb8600.NewMotoClient(address, username, password, logger),
	}
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc

	for _, desc := range downstreamMetrics {
		ch <- desc
	}

	for _, desc := range upstreamMetrics {
		ch <- desc
	}
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	var up float64 = 1

	_, err := c.Client.Login()
	if err == nil {
		downstream, err := c.Client.GetDownstreamChannels()
		if err != nil {
			level.Warn(c.Client.Logger).Log("msg", "Failed to get downstream channels", "err", err)
			up = 0
		}

		upstream, err := c.Client.GetUpstreamChannels()
		if err != nil {
			level.Warn(c.Client.Logger).Log("msg", "Failed to get upstream channels", "err", err)
			up = 0
		}

		for _, channel := range downstream {
			ch <- prometheus.MustNewConstMetric(
				downstreamInfo,
				prometheus.GaugeValue,
				1,
				fmt.Sprintf("%d", channel.Channel),
				fmt.Sprintf("%d", channel.ChannelID),
				channel.Modulation,
			)

			ch <- prometheus.MustNewConstMetric(
				downstreamFrequencyDesc,
				prometheus.GaugeValue,
				channel.Frequency,
				fmt.Sprintf("%d", channel.Channel),
			)

			ch <- prometheus.MustNewConstMetric(
				downstreamPowerDesc,
				prometheus.GaugeValue,
				channel.Power,
				fmt.Sprintf("%d", channel.Channel),
			)

			ch <- prometheus.MustNewConstMetric(
				downstreamSignalToNoiseDesc,
				prometheus.GaugeValue,
				channel.SignalToNoise,
				fmt.Sprintf("%d", channel.Channel),
			)

			ch <- prometheus.MustNewConstMetric(
				downstreamCorrectedErrorsDesc,
				prometheus.GaugeValue,
				channel.CorrectedErrors,
				fmt.Sprintf("%d", channel.Channel),
			)

			ch <- prometheus.MustNewConstMetric(
				downstreamUncorrectedErrorsDesc,
				prometheus.GaugeValue,
				channel.UncorrectedErrors,
				fmt.Sprintf("%d", channel.Channel),
			)

			status := 0
			if channel.LockStatus == "Locked" {
				status = 1
			}

			ch <- prometheus.MustNewConstMetric(
				downstreamLockStatus,
				prometheus.GaugeValue,
				float64(status),
				fmt.Sprintf("%d", channel.Channel),
			)
		}

		for _, channel := range upstream {
			ch <- prometheus.MustNewConstMetric(
				upstreamInfo,
				prometheus.GaugeValue,
				1,
				fmt.Sprintf("%d", channel.Channel),
				fmt.Sprintf("%d", channel.ChannelID),
				channel.ChannelType,
			)

			ch <- prometheus.MustNewConstMetric(
				upstreamFrequencyDesc,
				prometheus.GaugeValue,
				channel.Frequency,
				fmt.Sprintf("%d", channel.Channel),
			)

			ch <- prometheus.MustNewConstMetric(
				upstreamPowerDesc,
				prometheus.GaugeValue,
				channel.Power,
				fmt.Sprintf("%d", channel.Channel),
			)

			ch <- prometheus.MustNewConstMetric(
				upstreamSymbolRateDesc,
				prometheus.GaugeValue,
				channel.SymbolRate,
				fmt.Sprintf("%d", channel.Channel),
			)

			status := 0
			if channel.LockStatus == "Locked" {
				status = 1
			}

			ch <- prometheus.MustNewConstMetric(
				upstreamLockStatus,
				prometheus.GaugeValue,
				float64(status),
				fmt.Sprintf("%d", channel.Channel),
			)
		}
	}

	ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, up)
}
