package services

import (
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	ctrace "github.com/pip-services3-go/pip-services3-components-go/trace"
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

func (c *InstrumentTiming) EndTiming(err error) {
	if err == nil {
		c.EndSuccess()
	} else {
		c.EndFailure(err)
	}
}

func (c *InstrumentTiming) EndSuccess() {
	if c.counterTiming != nil {
		c.counterTiming.EndTiming()
	}
	if c.traceTiming != nil {
		c.traceTiming.EndTrace()
	}

	c.clear()
}

func (c *InstrumentTiming) EndFailure(err error) {
	if c.counterTiming != nil {
		c.counterTiming.EndTiming()
	}

	if err != nil {
		if c.logger != nil {
			c.logger.Error(c.correlationId, err, "Failed to call %s method", c.name)
		}
		if c.counters != nil {
			c.counters.IncrementOne(c.name + "." + c.verb + "_errors")
		}
		if c.traceTiming != nil {
			c.traceTiming.EndFailure(err)
		}
	} else {
		if c.traceTiming != nil {
			c.traceTiming.EndTrace()
		}
	}

	c.clear()
}
