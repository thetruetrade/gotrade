package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Minus Directional Movement Indicator (MinusDm), no storage, for use in other indicators
type MinusDmWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	previousHigh         float64
	previousLow          float64
	previousMinusDm      float64
	timePeriod           int
}

// NewMinusDmWithoutStorage creates a Minus Directional Movement Indicator (MinusDm) without storage
func NewMinusDmWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *MinusDmWithoutStorage, err error) {

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
		lookback = timePeriod - 1
	}
	ind := MinusDmWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        -1,
		previousMinusDm:      0.0,
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, nil
}

// A Minus Directional Movement Indicator (MinusDm)
type MinusDm struct {
	*MinusDmWithoutStorage

	// public variables
	Data []float64
}

// NewMinusDm creates a Minus Directional Movement Indicator (MinusDm) for online usage
func NewMinusDm(timePeriod int) (indicator *MinusDm, err error) {
	ind := MinusDm{}
	ind.MinusDmWithoutStorage, err = NewMinusDmWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultMinusDm creates a Minus Directional Movement Indicator (MinusDm) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultMinusDm() (indicator *MinusDm, err error) {
	timePeriod := 14
	return NewMinusDm(timePeriod)
}

// NewMinusDmWithSrcLen creates a Minus Directional Movement Indicator (MinusDm) for offline usage
func NewMinusDmWithSrcLen(sourceLength uint, timePeriod int) (indicator *MinusDm, err error) {
	ind, err := NewMinusDm(timePeriod)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultMinusDmWithSrcLen creates a Minus Directional Movement Indicator (MinusDm) for offline usage with default parameters
func NewDefaultMinusDmWithSrcLen(sourceLength uint) (indicator *MinusDm, err error) {
	ind, err := NewDefaultMinusDm()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewMinusDmForStream creates a Minus Directional Movement Indicator (MinusDm) for online usage with a source data stream
func NewMinusDmForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *MinusDm, err error) {
	ind, err := NewMinusDm(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMinusDmForStream creates a Minus Directional Movement Indicator (MinusDm) for online usage with a source data stream
func NewDefaultMinusDmForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *MinusDm, err error) {
	ind, err := NewDefaultMinusDm()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewMinusDmForStreamWithSrcLen creates a Minus Directional Movement Indicator (MinusDm) for offline usage with a source data stream
func NewMinusDmForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *MinusDm, err error) {
	ind, err := NewMinusDmWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMinusDmForStreamWithSrcLen creates a Minus Directional Movement Indicator (MinusDm) for offline usage with a source data stream
func NewDefaultMinusDmForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *MinusDm, err error) {
	ind, err := NewDefaultMinusDmWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *MinusDmWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	high := tickData.H()
	low := tickData.L()
	diffP := high - ind.previousHigh
	diffM := ind.previousLow - low

	if ind.lookbackPeriod == 1 {
		if ind.periodCounter > 0 {

			var result float64
			if (diffM > 0) && (diffP < diffM) {
				result = diffM
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
					ind.previousMinusDm += diffM
				}

				if ind.periodCounter == ind.timePeriod-1 {

					result := ind.previousMinusDm

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
				var result float64
				if (diffM > 0) && (diffP < diffM) {
					result = ind.previousMinusDm - (ind.previousMinusDm / float64(ind.timePeriod)) + diffM
				} else {
					result = ind.previousMinusDm - (ind.previousMinusDm / float64(ind.timePeriod))
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

				ind.previousMinusDm = result
			}
		}
	}

	ind.previousHigh = high
	ind.previousLow = low
}
