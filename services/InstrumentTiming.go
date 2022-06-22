package services

import (
	"context"
	ccount "github.com/pip-services3-gox/pip-services3-components-gox/count"
	clog "github.com/pip-services3-gox/pip-services3-components-gox/log"
	ctrace "github.com/pip-services3-gox/pip-services3-components-gox/trace"
)

type InstrumentTiming struct {
	correlationId string
	name          string
	verb          string
	logger        clog.ILogger
	counters      ccount.ICounters
	counterTiming *ccount.CounterTiming
	traceTiming   *ctrace.TraceTiming
}

func NewInstrumentTiming(correlationId string, name string,
	verb string, logger clog.ILogger, counters ccount.ICounters,
	counterTiming *ccount.CounterTiming, traceTiming *ctrace.TraceTiming) *InstrumentTiming {

	if len(verb) == 0 {
		verb = "call"
	}
	return &InstrumentTiming{
		correlationId: correlationId,
		name:          name,
		verb:          verb,
		logger:        logger,
		counters:      counters,
		counterTiming: counterTiming,
		traceTiming:   traceTiming,
	}
}

func (c *InstrumentTiming) clear() {
	// Clear references to avoid double processing
	c.counters = nil
	c.logger = nil
	c.counterTiming = nil
	c.traceTiming = nil
}

func (c *InstrumentTiming) EndTiming(ctx context.Context, err error) {
	if err == nil {
		c.EndSuccess(ctx)
	} else {
		c.EndFailure(ctx, err)
	}
}

func (c *InstrumentTiming) EndSuccess(ctx context.Context) {
	if c.counterTiming != nil {
		c.counterTiming.EndTiming(ctx)
	}
	if c.traceTiming != nil {
		c.traceTiming.EndTrace(ctx)
	}

	c.clear()
}

func (c *InstrumentTiming) EndFailure(ctx context.Context, err error) {
	if c.counterTiming != nil {
		c.counterTiming.EndTiming(ctx)
	}

	if err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, c.correlationId, err, "Failed to call %s method", c.name)
		}
		if c.counters != nil {
			c.counters.IncrementOne(ctx, c.name+"."+c.verb+"_errors")
		}
		if c.traceTiming != nil {
			c.traceTiming.EndFailure(ctx, err)
		}
	} else {
		if c.traceTiming != nil {
			c.traceTiming.EndTrace(ctx)
		}
	}

	c.clear()
}
