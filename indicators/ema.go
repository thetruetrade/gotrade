package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// An Exponential Moving Average Indicator (Ema), no storage, for use in other indicators
type EmaWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	periodTotal          float64
	periodCounter        int
	multiplier           float64
	previousEma          float64
	valueAvailableAction ValueAvailableActionFloat
	timePeriod           int
}

// NewEmaWithoutStorage creates an Exponential Moving Average Indicator (Ema) without storage
func NewEmaWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *EmaWithoutStorage, err error) {

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

	lookback := timePeriod - 1
	ind := EmaWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        timePeriod * -1,
		multiplier:           float64(2.0 / float64(timePeriod+1.0)),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, err
}

// An Exponential Moving Average Indicator (Ema)
type Ema struct {
	*EmaWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewEma creates an Exponential Moving Average Indicator (Ema) for online usage
func NewEma(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Ema, err error) {
	newEma := Ema{selectData: selectData}
	newEma.EmaWithoutStorage, err = NewEmaWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newEma.Data = append(newEma.Data, dataItem)
		})

	return &newEma, err
}

// NewDefaultEma creates an Exponential Moving Average (Ema) for online usage with default parameters
//	- timePeriod: 25
func NewDefaultEma() (indicator *Ema, err error) {
	timePeriod := 25
	return NewEma(timePeriod, gotrade.UseClosePrice)
}

// NewEmaWithSrcLen creates an Exponential Moving Average (Ema) for offline usage
func NewEmaWithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Ema, err error) {
	ind, err := NewEma(timePeriod, selectData)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultEmaWithSrcLen creates an Exponential Moving Average (Ema) for offline usage with default parameters
func NewDefaultEmaWithSrcLen(sourceLength int) (indicator *Ema, err error) {
	ind, err := NewDefaultEma()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewEmaForStream creates an Exponential Moving Average (Ema) for online usage with a source data stream
func NewEmaForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Ema, err error) {
	ind, err := NewEma(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultEmaForStream creates an Exponential Moving Average (Ema) for online usage with a source data stream
func NewDefaultEmaForStream(priceStream *gotrade.DOHLCVStream) (indicator *Ema, err error) {
	ind, err := NewDefaultEma()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewEmaForStreamWithSrcLen creates an Exponential Moving Average (Ema) for offline usage with a source data stream
func NewEmaForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Ema, err error) {
	ind, err := NewEmaWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultEmaForStreamWithSrcLen creates an Exponential Moving Average (Ema) for offline usage with a source data stream
func NewDefaultEmaForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Ema, err error) {
	ind, err := NewDefaultEmaWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ema *Ema) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ema.selectData(tickData)
	ema.ReceiveTick(selectedData, streamBarIndex)
}

func (ema *EmaWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ema.periodCounter += 1
	if ema.periodCounter < 0 {
		ema.periodTotal += tickData
	} else if ema.periodCounter == 0 {

		// increment the number of results this indicator can be expected to return
		ema.dataLength += 1

		if ema.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ema.validFromBar = streamBarIndex
		}

		ema.periodTotal += tickData
		result := ema.periodTotal / float64(ema.timePeriod)
		ema.previousEma = result

		// update the maximum result value
		if result > ema.maxValue {
			ema.maxValue = result
		}

		// update the minimum result value
		if result < ema.minValue {
			ema.minValue = result
		}

		// notify of a new result value though the value available action
		ema.valueAvailableAction(ema.previousEma, streamBarIndex)

	} else if ema.periodCounter > 0 {
		// increment the number of results this indicator can be expected to return
		ema.dataLength += 1

		result := (tickData-ema.previousEma)*ema.multiplier + ema.previousEma
		ema.previousEma = result

		// update the maximum result value
		if result > ema.maxValue {
			ema.maxValue = result
		}

		// update the minimum result value
		if result < ema.minValue {
			ema.minValue = result
		}

		// notify of a new result value though the value available action
		ema.valueAvailableAction(ema.previousEma, streamBarIndex)
	}
}
