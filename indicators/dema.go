// Double Exponential Moving Average (DEMA)
package indicators

// DEMA(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))

import (
	"github.com/thetruetrade/gotrade"
)

type DEMAWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableAction
	ema1                 *EMA
	ema2                 *EMA
	currentEMA           float64
}

func NewDEMAWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *DEMAWithoutStorage, err error) {
	newDEMA := DEMAWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(2 * (timePeriod - 1)),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod)}
	newDEMA.selectData = selectData
	newDEMA.valueAvailableAction = valueAvailableAction

	newDEMA.ema1, _ = NewEMA(timePeriod, selectData)

	newDEMA.ema1.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newDEMA.currentEMA = dataItem
		newDEMA.ema2.ReceiveTick(dataItem, streamBarIndex)
	}

	newDEMA.ema2, _ = NewEMA(timePeriod, selectData)

	newDEMA.ema2.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newDEMA.dataLength += 1
		if newDEMA.validFromBar == -1 {
			newDEMA.validFromBar = streamBarIndex
		}

		// DEMA(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))
		dema := (2 * newDEMA.currentEMA) - dataItem

		if dema > newDEMA.maxValue {
			newDEMA.maxValue = dema
		}

		if dema < newDEMA.minValue {
			newDEMA.minValue = dema
		}

		newDEMA.valueAvailableAction(dema, streamBarIndex)
	}

	return &newDEMA, nil
}

// A Double Exponential Moving Average Indicator
type DEMA struct {
	*DEMAWithoutStorage

	// public variables
	Data []float64
}

// NewDEMA returns a new Double Exponential Moving Average (DEMA) configured with the
// specified timePeriod. The DEMA results are stored in the DATA field.
func NewDEMA(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *DEMA, err error) {

	newDEMA := DEMA{}
	newDEMA.DEMAWithoutStorage, err = NewDEMAWithoutStorage(timePeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newDEMA.Data = append(newDEMA.Data, dataItem)
		})

	return &newDEMA, err
}

func NewDEMAForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *DEMA, err error) {
	newDEMA, err := NewDEMA(timePeriod, selectData)
	priceStream.AddTickSubscription(newDEMA)
	return newDEMA, err
}

func (dema *DEMAWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = dema.selectData(tickData)
	dema.ReceiveTick(selectedData, streamBarIndex)
}

func (dema *DEMAWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	dema.ema1.ReceiveTick(tickData, streamBarIndex)
}
