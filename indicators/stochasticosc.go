package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Stochastic Oscillator Indicator (StochOsc), no storage, for use in other indicators
type StochOscWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionStoch
	periodCounter        int
	slowKMA              *SmaWithoutStorage
	slowDMA              *SmaWithoutStorage
	hhv                  *HhvWithoutStorage
	llv                  *LlvWithoutStorage
	currentPeriodHigh    float64
	currentPeriodLow     float64
	currentFastK         float64
	currentSlowKMA       float64
	currentSlowDMA       float64
}

// NewStochOscWithoutStorage creates a Stochastic Oscillator Indicator (StochOsc) without storage
func NewStochOscWithoutStorage(fastKTimePeriod int, slowKTimePeriod int, slowDTimePeriod int, valueAvailableAction ValueAvailableActionStoch) (indicator *StochOscWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum fastKTimePeriod for this indicator is 1
	if fastKTimePeriod < 1 {
		return nil, errors.New("fastKTimePeriod is less than the minimum (1)")
	}

	// check the maximum fastKTimePeriod
	if fastKTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("fastKTimePeriod is greater than the maximum (100000)")
	}

	// the minimum slowKTimePeriod for this indicator is 1
	if slowKTimePeriod < 1 {
		return nil, errors.New("slowKTimePeriod is less than the minimum (1)")
	}

	// check the maximum slowKTimePeriod
	if slowKTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("slowKTimePeriod is greater than the maximum (100000)")
	}

	// the minimum slowDTimePeriod for this indicator is 1
	if slowDTimePeriod < 1 {
		return nil, errors.New("slowDTimePeriod is less than the minimum (1)")
	}

	// check the maximum slowDTimePeriod
	if slowDTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("slowDTimePeriod is greater than the maximum (100000)")
	}

	ind := StochOscWithoutStorage{
		baseFloatBounds:      newBaseFloatBounds(),
		currentSlowKMA:       0.0,
		currentSlowDMA:       0.0,
		periodCounter:        (fastKTimePeriod * -1),
		valueAvailableAction: valueAvailableAction,
	}

	tmpSlowKMA, err := NewSmaWithoutStorage(slowKTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentSlowKMA = dataItem
		ind.slowDMA.ReceiveTick(ind.currentSlowKMA, streamBarIndex)
	})

	tmpSlowDMA, err := NewSmaWithoutStorage(slowDTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentSlowDMA = dataItem

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

		var max = math.Max(ind.currentSlowKMA, ind.currentSlowDMA)
		var min = math.Min(ind.currentSlowKMA, ind.currentSlowDMA)

		// update the maximum result value
		if max > ind.maxValue {
			ind.maxValue = max
		}

		// update the minimum result value
		if min < ind.minValue {
			ind.minValue = min
		}

		// notify of a new result value though the value available action
		ind.valueAvailableAction(ind.currentSlowKMA, ind.currentSlowDMA, streamBarIndex)
	})

	lookback := fastKTimePeriod - 1 + tmpSlowDMA.GetLookbackPeriod() + tmpSlowKMA.GetLookbackPeriod()

	ind.baseIndicator = newBaseIndicator(lookback)
	ind.slowKMA = tmpSlowKMA
	ind.slowDMA = tmpSlowDMA
	ind.hhv, err = NewHhvWithoutStorage(fastKTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentPeriodHigh = dataItem
	})
	ind.llv, err = NewLlvWithoutStorage(fastKTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentPeriodLow = dataItem
	})

	return &ind, err
}

// A Stochastic Oscillator Indicator (StochOsc)
type StochOsc struct {
	*StochOscWithoutStorage

	// public variables
	SlowK []float64
	SlowD []float64
}

// NewStochOsc creates a Stochastic Oscillator Indicator (StochOsc) for online usage
func NewStochOsc(fastKTimePeriod int, slowKTimePeriod int, slowDTimePeriod int) (indicator *StochOsc, err error) {
	ind := StochOsc{}
	ind.StochOscWithoutStorage, err = NewStochOscWithoutStorage(fastKTimePeriod, slowKTimePeriod, slowDTimePeriod,
		func(dataItemK float64, dataItemD float64, streamBarIndex int) {
			ind.SlowK = append(ind.SlowK, dataItemK)
			ind.SlowD = append(ind.SlowD, dataItemD)
		})

	return &ind, err
}

// NewDefaultStochOsc creates a Stochastic Oscillator Indicator (StochOsc) for online usage with default parameters
//	- fastKTimePeriod : 5
//  - slowKTimePeriod : 3
//  - slowDTimePeriod : 3
func NewDefaultStochOsc() (indicator *StochOsc, err error) {
	fastKTimePeriod := 5
	slowKTimePeriod := 3
	slowDTimePeriod := 3
	return NewStochOsc(fastKTimePeriod, slowKTimePeriod, slowDTimePeriod)
}

// NewStochOscWithSrcLen creates a Stochastic Oscillator Indicator (StochOsc) for offline usage
func NewStochOscWithSrcLen(sourceLength uint, fastKTimePeriod int, slowKTimePeriod int, slowDTimePeriod int) (indicator *StochOsc, err error) {
	ind, err := NewStochOsc(fastKTimePeriod, slowKTimePeriod, slowDTimePeriod)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.SlowK = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
		ind.SlowD = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultStochOscWithSrcLen creates a Stochastic Oscillator Indicator (StochOsc) for offline usage with default parameters
func NewDefaultStochOscWithSrcLen(sourceLength uint) (indicator *StochOsc, err error) {
	ind, err := NewDefaultStochOsc()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.SlowK = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
		ind.SlowD = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewStochOscForStream creates a Stochastic Oscillator Indicator (StochOsc) for online usage with a source data stream
func NewStochOscForStream(priceStream gotrade.DOHLCVStreamSubscriber, fastKTimePeriod int, slowKTimePeriod int, slowDTimePeriod int) (indicator *StochOsc, err error) {
	ind, err := NewStochOsc(fastKTimePeriod, slowKTimePeriod, slowDTimePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultStochOscForStream creates a Stochastic Oscillator Indicator (StochOsc) for online usage with a source data stream
func NewDefaultStochOscForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *StochOsc, err error) {
	ind, err := NewDefaultStochOsc()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewStochOscForStreamWithSrcLen creates a Stochastic Oscillator Indicator (StochOsc) for offline usage with a source data stream
func NewStochOscForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, fastKTimePeriod int, slowKTimePeriod int, slowDTimePeriod int) (indicator *StochOsc, err error) {
	ind, err := NewStochOscWithSrcLen(sourceLength, fastKTimePeriod, slowKTimePeriod, slowDTimePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultStochOscForStreamWithSrcLen creates a Stochastic Oscillator Indicator (StochOsc) for offline usage with a source data stream
func NewDefaultStochOscForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *StochOsc, err error) {
	ind, err := NewDefaultStochOscWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *StochOscWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	ind.hhv.ReceiveTick(tickData.H(), streamBarIndex)
	ind.llv.ReceiveTick(tickData.L(), streamBarIndex)

	if ind.periodCounter >= 0 {
		ind.currentFastK = 100.0 * ((tickData.C() - ind.currentPeriodLow) / (ind.currentPeriodHigh - ind.currentPeriodLow))
		ind.slowKMA.ReceiveTick(ind.currentFastK, streamBarIndex)
	}
}
