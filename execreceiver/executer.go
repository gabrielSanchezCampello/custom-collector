// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package execreceiver

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

func commandToMetric(command string, cfg MetricCfg, dest pmetric.Metric, startTime pcommon.Timestamp, ts pcommon.Timestamp, scrapeCfg scraperhelper.ControllerConfig) error {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Errorf("error ejecutando el comando: %v", err)
		return err
	}

	value := strings.TrimSpace(string(output))

	dest.SetName(cfg.MetricName)
	dest.SetDescription(cfg.Description)
	dest.SetUnit(cfg.Unit)
	dataPointSlice := setMetricFields(cfg, dest)
	dataPoint := dataPointSlice.AppendEmpty()
	var errs []error

	setTimestamp(cfg, dataPoint, startTime, ts, scrapeCfg)

	err = setDataPointValue(cfg, value, dataPoint)
	if err != nil {
		errs = append(errs, fmt.Errorf("rowToMetric: %w", err))
	}
	attrs := dataPoint.Attributes()
	attrs.PutStr("command", command)
	for k, v := range cfg.StaticAttributes {
		attrs.PutStr(k, v)
	}

	return errors.Join(errs...)
}

func setTimestamp(cfg MetricCfg, dp pmetric.NumberDataPoint, startTime pcommon.Timestamp, ts pcommon.Timestamp, scrapeCfg scraperhelper.ControllerConfig) {
	dp.SetTimestamp(ts)

	// Cumulative sum should have a start time set to the beginning of the data points cumulation
	if cfg.Aggregation == MetricAggregationCumulative && cfg.DataType != MetricTypeGauge {
		dp.SetStartTimestamp(startTime)
	}

	// Non-cumulative sum should have a start time set to the previous endpoint
	if cfg.Aggregation == MetricAggregationDelta && cfg.DataType != MetricTypeGauge {
		dp.SetStartTimestamp(pcommon.NewTimestampFromTime(ts.AsTime().Add(-scrapeCfg.CollectionInterval)))
	}
}

func setMetricFields(cfg MetricCfg, dest pmetric.Metric) pmetric.NumberDataPointSlice {
	var out pmetric.NumberDataPointSlice
	switch cfg.DataType {
	case MetricTypeUnspecified, MetricTypeGauge:
		out = dest.SetEmptyGauge().DataPoints()
	case MetricTypeSum:
		sum := dest.SetEmptySum()
		sum.SetIsMonotonic(cfg.Monotonic)
		sum.SetAggregationTemporality(cfgToAggregationTemporality(cfg.Aggregation))
		out = sum.DataPoints()
	}
	return out
}

func cfgToAggregationTemporality(agg MetricAggregation) pmetric.AggregationTemporality {
	var out pmetric.AggregationTemporality
	switch agg {
	case MetricAggregationUnspecified, MetricAggregationCumulative:
		out = pmetric.AggregationTemporalityCumulative
	case MetricAggregationDelta:
		out = pmetric.AggregationTemporalityDelta
	}
	return out
}

func setDataPointValue(cfg MetricCfg, str string, dest pmetric.NumberDataPoint) error {
	switch cfg.ValueType {
	case MetricValueTypeUnspecified, MetricValueTypeInt:
		val, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("setDataPointValue: error converting to integer: %w", err)
		}
		dest.SetIntValue(int64(val))
	case MetricValueTypeDouble:
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("setDataPointValue: error converting to double: %w", err)
		}
		dest.SetDoubleValue(val)
	}
	return nil
}
