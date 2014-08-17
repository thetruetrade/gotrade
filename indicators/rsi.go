package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Relative Strength Indicator (Rsi), no storage, for use in other indicators
type RsiWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	previousClose        float64
	previousGain         float64
	previousLoss         float64
	timePeriod           int
}

// NewRsiWithoutStorage creates a Relative Strength Indicator (Rsi) without storage
func NewRsiWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *RsiWithoutStorage, err error) {

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

	lookback := timePeriod
	ind := RsiWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        (timePeriod * -1) - 1,
		previousClose:        0.0,
		previousGain:         0.0,
		previousLoss:         0.0,
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, err
}

// A Relative Strength Indicator (Rsi)
type Rsi struct {
	*RsiWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewRsi creates a Relative Strength Indicator (Rsi) for online usage
func NewRsi(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Rsi, err error) {
	ind := Rsi{selectData: selectData}
	ind.RsiWithoutStorage, err = NewRsiWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	ind.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	}
	return &ind, err
}

// NewDefaultRsi creates a Relative Strength Indicator (Rsi) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultRsi() (indicator *Rsi, err error) {
	timePeriod := 14
	return NewRsi(timePeriod, gotrade.UseClosePrice)
}

// NewRsiWithSrcLen creates a Relative Strength Indicator (Rsi) for offline usage
func NewRsiWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Rsi, err error) {
	ind, err := NewRsi(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultRsiWithSrcLen creates a Relative Strength Indicator (Rsi) for offline usage with default parameters
func NewDefaultRsiWithSrcLen(sourceLength uint) (indicator *Rsi, err error) {
	ind, err := NewDefaultRsi()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewRsiForStream creates a Relative Strength Indicator (Rsi) for online usage with a source data stream
func NewRsiForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Rsi, err error) {
	ind, err := NewRsi(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultRsiForStream creates a Relative Strength Indicator (Rsi) for online usage with a source data stream
func NewDefaultRsiForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Rsi, err error) {
	ind, err := NewDefaultRsi()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewRsiForStreamWithSrcLen creates a Relative Strength Indicator (Rsi) for offline usage with a source data stream
func NewRsiForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Rsi, err error) {
	ind, err := NewRsiWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultRsiForStreamWithSrcLen creates a Relative Strength Indicator (Rsi) for offline usage with a source data stream
func NewDefaultRsiForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Rsi, err error) {
	ind, err := NewDefaultRsiWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Rsi) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *RsiWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1

	if ind.periodCounter > ind.timePeriod*-1 {

		if ind.periodCounter <= 0 {

			if tickData > ind.previousClose {
				ind.previousGain += (tickData - ind.previousClose)
			} else {
				ind.previousLoss -= (tickData - ind.previousClose)
			}
		}

		if ind.periodCounter == 0 {
			ind.previousGain /= float64(ind.timePeriod)
			ind.previousLoss /= float64(ind.timePeriod)

			// increment the number of results this indicator can be expected to return
			ind.dataLength += 1
			if ind.validFromBar == -1 {
				// set the streamBarIndex from which this indicator returns valid results
				ind.validFromBar = streamBarIndex
			}

			var result float64
			//    Rsi = 100 * (prevGain/(prevGain+prevLoss))
			if ind.previousGain+ind.previousLoss == 0.0 {
				result = 0.0
			} else {
				result = 100.0 * (ind.previousGain / (ind.previousGain + ind.previousLoss))
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

		if ind.periodCounter > 0 {
			ind.previousGain *= float64(ind.timePeriod - 1)
			ind.previousLoss *= float64(ind.timePeriod - 1)

			if tickData > ind.previousClose {
				ind.previousGain += (tickData - ind.previousClose)
			} else {
				ind.previousLoss -= (tickData - ind.previousClose)
			}

			ind.previousGain /= float64(ind.timePeriod)
			ind.previousLoss /= float64(ind.timePeriod)

			// increment the number of results this indicator can be expected to return
			ind.dataLength += 1

			var result float64
			//    Rsi = 100 * (prevGain/(prevGain+prevLoss))
			if ind.previousGain+ind.previousLoss == 0.0 {
				result = 0.0
			} else {
				result = 100.0 * (ind.previousGain / (ind.previousGain + ind.previousLoss))
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
	ind.previousClose = tickData
}
