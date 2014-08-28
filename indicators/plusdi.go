package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Plus Directional Indicator (PlusDi), no storage, for use in other indicators
type PlusDiWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodCounter     int
	previousHigh      float64
	previousLow       float64
	previousPlusDM    float64
	previousTrueRange float64
	currentTrueRange  float64
	trueRange         *TrueRange
	timePeriod        int
}

// NewPlusDiWithoutStorage creates a Plus Directional Indicator (PlusDi) without storage
func NewPlusDiWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *PlusDiWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for this indicator is 1
	if timePeriod < 1 {
		return nil, errors.New("timePeriod is less than the minimum (1)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := 1
	if timePeriod > 1 {
		lookback = timePeriod
	}
	ind := PlusDiWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		periodCounter:                -1,
		previousPlusDM:               0.0,
		previousTrueRange:            0.0,
		currentTrueRange:             0.0,
		timePeriod:                   timePeriod,
	}

	ind.trueRange, err = NewTrueRange()

	ind.trueRange.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		ind.currentTrueRange = dataItem
	}

	return &ind, nil
}

// A Plus Directional Indicator (PlusDi)
type PlusDi struct {
	*PlusDiWithoutStorage

	// public variables
	Data []float64
}

// NewPlusDi creates a Plus Directional Indicator (PlusDi) for online usage
func NewPlusDi(timePeriod int) (indicator *PlusDi, err error) {
	ind := PlusDi{}
	ind.PlusDiWithoutStorage, err = NewPlusDiWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultPlusDi creates a Plus Directional Indicator (PlusDi) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultPlusDi() (indicator *PlusDi, err error) {
	timePeriod := 14
	return NewPlusDi(timePeriod)
}

// NewPlusDiWithSrcLen creates a Plus Directional Indicator (PlusDi) for offline usage
func NewPlusDiWithSrcLen(sourceLength uint, timePeriod int) (indicator *PlusDi, err error) {
	ind, err := NewPlusDi(timePeriod)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultPlusDiWithSrcLen creates a Plus Directional Indicator (PlusDi) for offline usage with default parameters
func NewDefaultPlusDiWithSrcLen(sourceLength uint) (indicator *PlusDi, err error) {
	ind, err := NewDefaultPlusDi()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewPlusDiForStream creates a Plus Directional Indicator (PlusDi) for online usage with a source data stream
func NewPlusDiForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *PlusDi, err error) {
	ind, err := NewPlusDi(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultPlusDiForStream creates a Plus Directional Indicator (PlusDi) for online usage with a source data stream
func NewDefaultPlusDiForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *PlusDi, err error) {
	ind, err := NewDefaultPlusDi()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewPlusDiForStreamWithSrcLen creates a Plus Directional Indicator (PlusDi) for offline usage with a source data stream
func NewPlusDiForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *PlusDi, err error) {
	ind, err := NewPlusDiWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultPlusDiForStreamWithSrcLen creates a Plus Directional Indicator (PlusDi) for offline usage with a source data stream
func NewDefaultPlusDiForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *PlusDi, err error) {
	ind, err := NewDefaultPlusDiWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *PlusDiWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	// forward to the true range indicator first using previous data
	ind.trueRange.ReceiveDOHLCVTick(tickData, streamBarIndex)

	ind.periodCounter += 1
	high := tickData.H()
	low := tickData.L()
	diffP := high - ind.previousHigh
	diffM := ind.previousLow - low

	if ind.lookbackPeriod == 1 {
		if ind.periodCounter > 0 {

			// forward to the true range indicator first using previous data
			ind.trueRange.ReceiveDOHLCVTick(tickData, streamBarIndex)

			var result float64
			if (diffP > 0) && (diffP > diffM) && ind.currentTrueRange != 0.0 {
				result = diffP / ind.currentTrueRange
			} else {
				result = 0
			}

			ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
		}
	} else {
		if ind.periodCounter > 0 {
			if ind.periodCounter < ind.timePeriod {
				if (diffP > 0) && (diffP > diffM) {
					ind.previousPlusDM += diffP
				}
				ind.previousTrueRange += ind.currentTrueRange
			} else {
				var result float64
				ind.previousTrueRange = ind.previousTrueRange - (ind.previousTrueRange / float64(ind.timePeriod)) + ind.currentTrueRange
				if (diffP > 0) && (diffP > diffM) {
					ind.previousPlusDM = ind.previousPlusDM - (ind.previousPlusDM / float64(ind.timePeriod)) + diffP
				} else {
					ind.previousPlusDM = ind.previousPlusDM - (ind.previousPlusDM / float64(ind.timePeriod))
				}

				if ind.previousTrueRange != 0.0 {
					result = float64(100.0) * ind.previousPlusDM / ind.previousTrueRange
				} else {
					result = 0.0
				}

				ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
			}
		}
	}

	ind.previousHigh = high
	ind.previousLow = low
}
