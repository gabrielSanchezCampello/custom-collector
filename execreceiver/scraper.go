// scraper
package execreceiver

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"go.opentelemetry.io/collector/scraper"
	"go.opentelemetry.io/collector/scraper/scrapererror"
)

type Scraper struct {
	id        component.ID
	Query     Query
	ScrapeCfg scraperhelper.ControllerConfig
	StartTime pcommon.Timestamp
}

var ErrNullValueWarning = errors.New("NULL value")

var _ scraper.Metrics = (*Scraper)(nil)

func NewScraper(id component.ID, query Query, scrapeCfg scraperhelper.ControllerConfig) *Scraper {
	return &Scraper{
		id:        id,
		Query:     query,
		ScrapeCfg: scrapeCfg,
	}
}

func (s *Scraper) ID() component.ID {
	return s.id
}

func (s *Scraper) Start(context.Context, component.Host) error {
	s.StartTime = pcommon.NewTimestampFromTime(time.Now())
	return nil
}

func (s *Scraper) ScrapeMetrics(ctx context.Context) (pmetric.Metrics, error) {
	out := pmetric.NewMetrics()
	ts := pcommon.NewTimestampFromTime(time.Now())
	rms := out.ResourceMetrics()
	rm := rms.AppendEmpty()
	sms := rm.ScopeMetrics()
	sm := sms.AppendEmpty()
	ms := sm.Metrics()
	var errs []error
	metricCfg := s.Query.Metric
	if err := commandToMetric(s.Query.COMMAND, metricCfg, ms.AppendEmpty(), s.StartTime, ts, s.ScrapeCfg); err != nil {
		errs = append(errs, err)
	}
	if errs != nil {
		return out, scrapererror.NewPartialScrapeError(errors.Join(errs...), len(errs))
	}
	return out, nil
}

func (s *Scraper) Shutdown(_ context.Context) error {
	return nil
}
