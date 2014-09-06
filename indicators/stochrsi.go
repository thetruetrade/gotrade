package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Stochastic Relative Strength Indicator (StochRsi), no storage, for use in other indicators
type StochRsiWithoutStorage struct {
	*baseIndicatorWithFloatBoundsStoch

	// private variables
	periodCounter     int
	fastDMA           *SmaWithoutStorage
	rsi               *RsiWithoutStorage
	hhv               *HhvWithoutStorage
	llv               *LlvWithoutStorage
	currentRSI        float64
	currentPeriodHigh float64
	currentPeriodLow  float64
	currentFastK      float64
	currentFastDMA    float64
}

// NewStochRsiWithoutStorage creates a Stochastic Relative Strength Indicator (StochRsi) without storage
func NewStochRsiWithoutStorage(timePeriod int, fastKTimePeriod int, fastDTimePeriod int, valueAvailableAction ValueAvailableActionStoch) (indicator *StochRsiWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for this indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod is less than the minimum (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	// the minimum fastKTimePeriod for this indicator is 1
	if fastKTimePeriod < 1 {
		return nil, errors.New("fastKTimePeriod is less than the minimum (1)")
	}

	// check the maximum fastKTimePeriod
	if fastKTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("fastKTimePeriod is greater than the maximum (100000)")
	}

	// the minimum fastDTimePeriod for this indicator is 1
	if fastDTimePeriod < 1 {
		return nil, errors.New("fastDTimePeriod is less than the minimum (1)")
	}

	// check the maximum fastDTimePeriod
	if fastDTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("fastDTimePeriod is greater than the maximum (100000)")
	}

	ind := StochRsiWithoutStorage{
		currentFastDMA: 0.0,
		periodCounter:  (fastKTimePeriod * -1),
	}

	tmpRSI, err := NewRsiWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.periodCounter += 1

		ind.currentRSI = dataItem
		ind.hhv.ReceiveTick(dataItem, streamBarIndex)
		ind.llv.ReceiveTick(dataItem, streamBarIndex)

		if ind.periodCounter >= 0 {
			diff := ind.currentPeriodHigh - ind.currentPeriodLow
			if diff != 0 {
				ind.currentFastK = 100.0 * ((ind.currentRSI - ind.currentPeriodLow) / diff)
			} else {
				ind.currentFastK = 0
			}
			ind.fastDMA.ReceiveTick(ind.currentFastK, streamBarIndex)
		}
	})

	tmpFastDMA, err := NewSmaWithoutStorage(fastDTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentFastDMA = dataItem

		ind.UpdateIndicatorWithNewValue(ind.currentFastK, ind.currentFastDMA, streamBarIndex)
	})

	totalTimePeriod := tmpRSI.GetLookbackPeriod() + tmpFastDMA.GetLookbackPeriod() + fastKTimePeriod - 1

	ind.baseIndicatorWithFloatBoundsStoch = newBaseIndicatorWithFloatBoundsStoch(totalTimePeriod, valueAvailableAction)
	ind.fastDMA = tmpFastDMA
	ind.rsi = tmpRSI
	ind.hhv, err = NewHhvWithoutStorage(fastKTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentPeriodHigh = dataItem
	})
	ind.llv, err = NewLlvWithoutStorage(fastKTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentPeriodLow = dataItem
	})

	return &ind, err
}

// A Stochastic Relative Strength Indicator (StochRsi)
type StochRsi struct {
	*StochRsiWithoutStorage

	// public variables
	SlowK []float64
	SlowD []float64
}

// NewStochRsi creates a Stochastic Relative Strength Indicator (StochRsi) for online usage
func NewStochRsi(timePeriod int, fastKTimePeriod int, fastDTimePeriod int) (indicator *StochRsi, err error) {
	newStochRsi := StochRsi{}
	newStochRsi.StochRsiWithoutStorage, err = NewStochRsiWithoutStorage(timePeriod, fastKTimePeriod, fastDTimePeriod,
		func(dataItemK float64, dataItemD float64, streamBarIndex int) {
			newStochRsi.SlowK = append(newStochRsi.SlowK, dataItemK)
			newStochRsi.SlowD = append(newStochRsi.SlowD, dataItemD)
		})

	return &newStochRsi, err
}

// NewDefaultStochRsi creates a Stochastic Relative Strength Indicator (StochRsi) for online usage with default parameters
//	- timePeriod : 14
//  - fastKTimePeriod : 5
//  - fastDTimePeriod : 3
func NewDefaultStochRsi() (indicator *StochRsi, err error) {
	timePeriod := 14
	fastKTimePeriod := 5
	fastDTimePeriod := 3
	return NewStochRsi(timePeriod, fastKTimePeriod, fastDTimePeriod)
}

// NewStochRsiWithSrcLen creates a Stochastic Relative Strength Indicator (StochRsi) for offline usage
func NewStochRsiWithSrcLen(sourceLength uint, timePeriod int, fastKTimePeriod int, fastDTimePeriod int) (indicator *StochRsi, err error) {
	ind, err := NewStochRsi(timePeriod, fastKTimePeriod, fastDTimePeriod)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.SlowK = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
		ind.SlowD = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultStochRsiWithSrcLen creates a Stochastic Relative Strength Indicator (StochRsi) for offline usage with default parameters
func NewDefaultStochRsiWithSrcLen(sourceLength uint) (indicator *StochRsi, err error) {
	ind, err := NewDefaultStochRsi()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.SlowK = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
		ind.SlowD = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewStochRsiForStream creates a Stochastic Relative Strength Indicator (StochRsi) for online usage with a source data stream
func NewStochRsiForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, fastKTimePeriod int, fastDTimePeriod int) (indicator *StochRsi, err error) {
	ind, err := NewStochRsi(timePeriod, fastKTimePeriod, fastDTimePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultStochRsiForStream creates a Stochastic Relative Strength Indicator (StochRsi) for online usage with a source data stream
func NewDefaultStochRsiForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *StochRsi, err error) {
	ind, err := NewDefaultStochRsi()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewStochRsiForStreamWithSrcLen creates a Stochastic Relative Strength Indicator (StochRsi) for offline usage with a source data stream
func NewStochRsiForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, fastKTimePeriod int, fastDTimePeriod int) (indicator *StochRsi, err error) {
	ind, err := NewStochRsiWithSrcLen(sourceLength, timePeriod, fastKTimePeriod, fastDTimePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultStochRsiForStreamWithSrcLen creates a Stochastic Relative Strength Indicator (StochRsi) for offline usage with a source data stream
func NewDefaultStochRsiForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *StochRsi, err error) {
	ind, err := NewDefaultStochRsiWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *StochRsiWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	ind.rsi.ReceiveTick(tickData.C(), streamBarIndex)
}
