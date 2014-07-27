package indicators

// Dema(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// An Average True Range Indicator (Dema), no storage, for use in other indicators
type DemaWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	ema1                 *EMAWithoutStorage
	ema2                 *EMAWithoutStorage
	currentEMA           float64
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
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		valueAvailableAction: valueAvailableAction,
	}

	ind.ema1, _ = NewEMAWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentEMA = dataItem
		ind.ema2.ReceiveTick(dataItem, streamBarIndex)
	})

	ind.ema2, _ = NewEMAWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

		// Dema(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))
		dema := (2 * ind.currentEMA) - dataItem

		// update the maximum result value
		if dema > ind.maxValue {
			ind.maxValue = dema
		}

		// update the minimum result value
		if dema < ind.minValue {
			ind.minValue = dema
		}

		// notify of a new result value though the value available action
		ind.valueAvailableAction(dema, streamBarIndex)
	})

	return &ind, nil
}

// A Double Exponential Moving Average Indicator (Dema)
type Dema struct {
	*DemaWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewDema creates a Double Exponential Moving Average (Dema) for online usage
func NewDema(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Dema, err error) {

	newDema := Dema{selectData: selectData}
	newDema.DemaWithoutStorage, err = NewDemaWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newDema.Data = append(newDema.Data, dataItem)
		})

	return &newDema, err
}

// NewDefaultDema creates a Double Exponential Moving Average (Dema) for online usage with default parameters
//	- timePeriod: 30
//  - selectData: useClosePrice
func NewDefaultDema() (indicator *Dema, err error) {
	timePeriod := 30
	selectData := gotrade.UseClosePrice
	return NewDema(timePeriod, selectData)
}

// NewDemaWithKnownSourceLength creates a Double Exponential Moving Average (Dema) for offline usage
func NewDemaWithKnownSourceLength(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Dema, err error) {
	ind, err := NewDema(timePeriod, selectData)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultDemaWithKnownSourceLength creates a Double Exponential Moving Average (Dema) for offline usage with default parameters
func NewDefaultDemaWithKnownSourceLength(sourceLength int) (indicator *Dema, err error) {
	ind, err := NewDefaultDema()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewDemaForStream creates a Double Exponential Moving Average (Dema) for online usage with a source data stream
func NewDemaForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Dema, err error) {
	newDema, err := NewDema(timePeriod, selectData)
	priceStream.AddTickSubscription(newDema)
	return newDema, err
}

// NewDefaultDemaForStream creates a Double Exponential Moving Average (Dema) for online usage with a source data stream
func NewDefaultDemaForStream(priceStream *gotrade.DOHLCVStream) (indicator *Dema, err error) {
	ind, err := NewDefaultDema()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDemaForStreamWithKnownSourceLength creates a Double Exponential Moving Average (Dema) for offline usage with a source data stream
func NewDemaForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Dema, err error) {
	ind, err := NewDemaWithKnownSourceLength(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultDemaForStreamWithKnownSourceLength creates a Double Exponential Moving Average (Dema) for offline usage with a source data stream
func NewDefaultDemaForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Dema, err error) {
	ind, err := NewDefaultDemaWithKnownSourceLength(sourceLength)
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
