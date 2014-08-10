package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Minus Directional Indicator (MinusDi), no storage, for use in other indicators
type MinusDiWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	previousHigh         float64
	previousLow          float64
	previousMinusDM      float64
	previousTrueRange    float64
	currentTrueRange     float64
	trueRange            *TrueRange
	timePeriod           int
}

// NewMinusDiWithoutStorage creates a Minus Directional Indicator (MinusDi) without storage
func NewMinusDiWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *MinusDiWithoutStorage, err error) {

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
	ind := MinusDiWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        -1,
		previousMinusDM:      0.0,
		previousTrueRange:    0.0,
		currentTrueRange:     0.0,
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	ind.trueRange, err = NewTrueRange()

	ind.trueRange.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		ind.currentTrueRange = dataItem
	}

	return &ind, nil
}

// A Minus Directional Indicator (MinusDi)
type MinusDi struct {
	*MinusDiWithoutStorage

	// public variables
	Data []float64
}

// NewMinusDi creates a Minus Directional Indicator (MinusDi) for online usage
func NewMinusDi(timePeriod int) (indicator *MinusDi, err error) {
	ind := MinusDi{}
	ind.MinusDiWithoutStorage, err = NewMinusDiWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultMinusDi creates a Minus Directional Indicator (MinusDi) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultMinusDi() (indicator *MinusDi, err error) {
	timePeriod := 14
	return NewMinusDi(timePeriod)
}

// NewMinusDiWithSrcLen creates a Minus Directional Indicator (MinusDi) for offline usage
func NewMinusDiWithSrcLen(sourceLength int, timePeriod int) (indicator *MinusDi, err error) {
	ind, err := NewMinusDi(timePeriod)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultMinusDiWithSrcLen creates a Minus Directional Indicator (MinusDi) for offline usage with default parameters
func NewDefaultMinusDiWithSrcLen(sourceLength int) (indicator *MinusDi, err error) {
	ind, err := NewDefaultMinusDi()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewMinusDiForStream creates a Minus Directional Indicator (MinusDi) for online usage with a source data stream
func NewMinusDiForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *MinusDi, err error) {
	ind, err := NewMinusDi(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMinusDiForStream creates a Minus Directional Indicator (MinusDi) for online usage with a source data stream
func NewDefaultMinusDiForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *MinusDi, err error) {
	ind, err := NewDefaultMinusDi()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewMinusDiForStreamWithSrcLen creates a Minus Directional Indicator (MinusDi) for offline usage with a source data stream
func NewMinusDiForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *MinusDi, err error) {
	ind, err := NewMinusDiWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMinusDiForStreamWithSrcLen creates a Minus Directional Indicator (MinusDi) for offline usage with a source data stream
func NewDefaultMinusDiForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *MinusDi, err error) {
	ind, err := NewDefaultMinusDiWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *MinusDiWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

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
			if (diffM > 0) && (diffP < diffM) && ind.currentTrueRange != 0.0 {
				result = diffM / ind.currentTrueRange
			} else {
				result = 0
			}

			// increment the number of results this indicator can be expected to return
			ind.dataLength += 1

			if ind.validFromBar == -1 {
				// set the streamBarIndex from which this indicator returns valid results
				ind.validFromBar = streamBarIndex
			}

			// update the maximum result value
			if result > ind.maxValue {
				ind.maxValue = result
			}

			// update the minimum result value
			if result < ind.minValue {
				ind.minValue = result
			}

			// notify of a new result value though the value available action
			ind.valueAvailableAction(result, streamBarIndex)
		}
	} else {
		if ind.periodCounter > 0 {
			if ind.periodCounter < ind.timePeriod {
				if (diffM > 0) && (diffP < diffM) {
					ind.previousMinusDM += diffM
				}
				ind.previousTrueRange += ind.currentTrueRange
			} else {
				var result float64
				ind.previousTrueRange = ind.previousTrueRange - (ind.previousTrueRange / float64(ind.timePeriod)) + ind.currentTrueRange
				if (diffM > 0) && (diffP < diffM) {
					ind.previousMinusDM = ind.previousMinusDM - (ind.previousMinusDM / float64(ind.timePeriod)) + diffM
				} else {
					ind.previousMinusDM = ind.previousMinusDM - (ind.previousMinusDM / float64(ind.timePeriod))
				}

				if ind.previousTrueRange != 0.0 {
					result = float64(100.0) * ind.previousMinusDM / ind.previousTrueRange
				} else {
					result = 0.0
				}

				// increment the number of results this indicator can be expected to return
				ind.dataLength += 1

				if ind.validFromBar == -1 {
					// set the streamBarIndex from which this indicator returns valid results
					ind.validFromBar = streamBarIndex
				}

				// update the maximum result value
				if result > ind.maxValue {
					ind.maxValue = result
				}

				// update the minimum result value
				if result < ind.minValue {
					ind.minValue = result
				}

				// notify of a new result value though the value available action
				ind.valueAvailableAction(result, streamBarIndex)
			}
		}
	}

	ind.previousHigh = high
	ind.previousLow = low
}
