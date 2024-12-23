// execreiver project execreiver.go
package execreceiver

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

func createMetricsReceiver(
	_ context.Context,
	settings receiver.Settings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	sqlCfg := cfg.(*Config)
	var opts []scraperhelper.ScraperControllerOption
	for i, query := range sqlCfg.Queries {
		id := component.MustNewIDWithName("execreceiver", fmt.Sprintf("query-%d: %s", i, query.COMMAND))
		mp := NewScraper(id, query, sqlCfg.ControllerConfig)

		opt := scraperhelper.AddScraper(component.MustNewType("execreceiver"), mp)
		opts = append(opts, opt)
	}
	return scraperhelper.NewScraperControllerReceiver(
		&sqlCfg.ControllerConfig,
		settings,
		consumer,
		opts...,
	)
}
