// config
package execreceiver

import (
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`
	Queries                        []Query `mapstructure:"queries"`
}

func createDefaultConfig() component.Config {
	cfg := scraperhelper.NewDefaultControllerConfig()
	cfg.CollectionInterval = 10 * time.Second
	return &Config{
		ControllerConfig: cfg,
	}
}

func (c Config) Validate() error {
	if len(c.Queries) == 0 {
		return errors.New("'queries' cannot be empty")
	}
	for _, query := range c.Queries {
		if err := query.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type Query struct {
	COMMAND string    `mapstructure:"command"`
	Metric  MetricCfg `mapstructure:"metric"`
}

func (q Query) Validate() error {
	var errs []error
	if q.COMMAND == "" {
		errs = append(errs, errors.New("'command' cannot be empty"))
	}
	if err := q.Metric.Validate(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

type MetricCfg struct {
	MetricName       string            `mapstructure:"metric_name"`
	Monotonic        bool              `mapstructure:"monotonic"`
	ValueType        MetricValueType   `mapstructure:"value_type"`
	DataType         MetricType        `mapstructure:"data_type"`
	Aggregation      MetricAggregation `mapstructure:"aggregation"`
	Unit             string            `mapstructure:"unit"`
	Description      string            `mapstructure:"description"`
	StaticAttributes map[string]string `mapstructure:"static_attributes"`
	StartTsColumn    string            `mapstructure:"start_ts_column"`
	TsColumn         string            `mapstructure:"ts_column"`
}

func (c MetricCfg) Validate() error {
	var errs []error
	if c.MetricName == "" {
		errs = append(errs, errors.New("'metric_name' cannot be empty"))
	}

	if err := c.ValueType.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := c.DataType.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := c.Aggregation.Validate(); err != nil {
		errs = append(errs, err)
	}
	if c.DataType == MetricTypeGauge && c.Aggregation != "" {
		errs = append(errs, fmt.Errorf("aggregation=%s but data_type=%s does not support aggregation", c.Aggregation, c.DataType))
	}
	if errs != nil && c.MetricName != "" {
		errs = append(errs, fmt.Errorf("invalid metric config with metric_name '%s'", c.MetricName))
	}
	return errors.Join(errs...)
}

type MetricType string

const (
	MetricTypeUnspecified MetricType = ""
	MetricTypeGauge       MetricType = "gauge"
	MetricTypeSum         MetricType = "sum"
)

func (t MetricType) Validate() error {
	switch t {
	case MetricTypeUnspecified, MetricTypeGauge, MetricTypeSum:
		return nil
	}
	return fmt.Errorf("metric config has unsupported data_type: '%s'", t)
}

type MetricValueType string

const (
	MetricValueTypeUnspecified MetricValueType = ""
	MetricValueTypeInt         MetricValueType = "int"
	MetricValueTypeDouble      MetricValueType = "double"
)

func (t MetricValueType) Validate() error {
	switch t {
	case MetricValueTypeUnspecified, MetricValueTypeInt, MetricValueTypeDouble:
		return nil
	}
	return fmt.Errorf("metric config has unsupported value_type: '%s'", t)
}

type MetricAggregation string

const (
	MetricAggregationUnspecified MetricAggregation = ""
	MetricAggregationCumulative  MetricAggregation = "cumulative"
	MetricAggregationDelta       MetricAggregation = "delta"
)

func (a MetricAggregation) Validate() error {
	switch a {
	case MetricAggregationUnspecified, MetricAggregationCumulative, MetricAggregationDelta:
		return nil
	}
	return fmt.Errorf("metric config has unsupported aggregation: '%s'", a)
}
