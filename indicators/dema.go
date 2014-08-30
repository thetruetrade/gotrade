package indicators

// Dema(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Double Exponential Moving Average Indicator (Dema), no storage, for use in other indicators
type DemaWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	ema1       *EmaWithoutStorage
	ema2       *EmaWithoutStorage
	currentEMA float64
}

// NewDemaWithoutStorage creates a Double Exponential Moving Average Indicator (Dema) without storage
func NewDemaWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *DemaWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for a Dema indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod is less than the minimum (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := 2 * (timePeriod - 1)
	ind := DemaWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
	}

	ind.ema1, _ = NewEmaWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentEMA = dataItem
		ind.ema2.ReceiveTick(dataItem, streamBarIndex)
	})

	ind.ema2, _ = NewEmaWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {

		// Dema(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))
		result := (2 * ind.currentEMA) - dataItem

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
	})

	return &ind, nil
}

// A Double Exponential Moving Average Indicator (Dema)
type Dema struct {
	*DemaWithoutStorage
	selectData gotrade.DOHLCVDataSelectionFunc

	// public variables
	Data []float64
}

// NewDema creates a Double Exponential Moving Average (Dema) for online usage
func NewDema(timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Dema, err error) {

	if selectData == nil {
		return nil, ErrDOHLCVDataSelectFuncIsNil
	}

	ind := Dema{
		selectData: selectData,
	}

	ind.DemaWithoutStorage, err = NewDemaWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultDema creates a Double Exponential Moving Average (Dema) for online usage with default parameters
//	- timePeriod: 30
//  - selectData: useClosePrice
func NewDefaultDema() (indicator *Dema, err error) {
	timePeriod := 30
	selectData := gotrade.UseClosePrice
	return NewDema(timePeriod, selectData)
}

// NewDemaWithSrcLen creates a Double Exponential Moving Average (Dema) for offline usage
func NewDemaWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Dema, err error) {
	ind, err := NewDema(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultDemaWithSrcLen creates a Double Exponential Moving Average (Dema) for offline usage with default parameters
func NewDefaultDemaWithSrcLen(sourceLength uint) (indicator *Dema, err error) {
	ind, err := NewDefaultDema()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDemaForStream creates a Double Exponential Moving Average (Dema) for online usage with a source data stream
func NewDemaForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Dema, err error) {
	newDema, err := NewDema(timePeriod, selectData)
	priceStream.AddTickSubscription(newDema)
	return newDema, err
}

// NewDefaultDemaForStream creates a Double Exponential Moving Average (Dema) for online usage with a source data stream
func NewDefaultDemaForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Dema, err error) {
	ind, err := NewDefaultDema()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDemaForStreamWithSrcLen creates a Double Exponential Moving Average (Dema) for offline usage with a source data stream
func NewDemaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Dema, err error) {
	ind, err := NewDemaWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultDemaForStreamWithSrcLen creates a Double Exponential Moving Average (Dema) for offline usage with a source data stream
func NewDefaultDemaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Dema, err error) {
	ind, err := NewDefaultDemaWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (dema *Dema) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = dema.selectData(tickData)
	dema.ReceiveTick(selectedData, streamBarIndex)
}

func (dema *DemaWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	dema.ema1.ReceiveTick(tickData, streamBarIndex)
}
