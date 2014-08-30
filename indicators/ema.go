package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// An Exponential Moving Average Indicator (Ema), no storage, for use in other indicators
type EmaWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodTotal   float64
	periodCounter int
	multiplier    float64
	previousEma   float64
	timePeriod    int
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
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		periodCounter:                timePeriod * -1,
		multiplier:                   float64(2.0 / float64(timePeriod+1.0)),
		timePeriod:                   timePeriod,
	}

	return &ind, err
}

// An Exponential Moving Average Indicator (Ema)
type Ema struct {
	*EmaWithoutStorage
	selectData gotrade.DOHLCVDataSelectionFunc

	// public variables
	Data []float64
}

// NewEma creates an Exponential Moving Average Indicator (Ema) for online usage
func NewEma(timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Ema, err error) {
	if selectData == nil {
		return nil, ErrDOHLCVDataSelectFuncIsNil
	}

	ind := Ema{
		selectData: selectData,
	}

	ind.EmaWithoutStorage, err = NewEmaWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultEma creates an Exponential Moving Average (Ema) for online usage with default parameters
//	- timePeriod: 25
func NewDefaultEma() (indicator *Ema, err error) {
	timePeriod := 25
	return NewEma(timePeriod, gotrade.UseClosePrice)
}

// NewEmaWithSrcLen creates an Exponential Moving Average (Ema) for offline usage
func NewEmaWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Ema, err error) {
	ind, err := NewEma(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultEmaWithSrcLen creates an Exponential Moving Average (Ema) for offline usage with default parameters
func NewDefaultEmaWithSrcLen(sourceLength uint) (indicator *Ema, err error) {
	ind, err := NewDefaultEma()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewEmaForStream creates an Exponential Moving Average (Ema) for online usage with a source data stream
func NewEmaForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Ema, err error) {
	ind, err := NewEma(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultEmaForStream creates an Exponential Moving Average (Ema) for online usage with a source data stream
func NewDefaultEmaForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Ema, err error) {
	ind, err := NewDefaultEma()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewEmaForStreamWithSrcLen creates an Exponential Moving Average (Ema) for offline usage with a source data stream
func NewEmaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Ema, err error) {
	ind, err := NewEmaWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultEmaForStreamWithSrcLen creates an Exponential Moving Average (Ema) for offline usage with a source data stream
func NewDefaultEmaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Ema, err error) {
	ind, err := NewDefaultEmaWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Ema) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *EmaWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	if ind.periodCounter < 0 {
		ind.periodTotal += tickData
	} else if ind.periodCounter == 0 {

		ind.periodTotal += tickData
		result := ind.periodTotal / float64(ind.timePeriod)
		ind.previousEma = result

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)

	} else if ind.periodCounter > 0 {

		result := (tickData-ind.previousEma)*ind.multiplier + ind.previousEma
		ind.previousEma = result

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
	}
}
