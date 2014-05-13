// Double Exponential Moving Average (DEMA)
package indicators

// DEMA(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))

import (
	"github.com/thetruetrade/gotrade"
)

type baseDEMA struct {
	*baseIndicatorWithLookback

	// private variables
	valueAvailableAction ValueAvailableAction
	ema1                 *EMA
	ema2                 *EMA
	currentEMA           float64
}

func newBaseDEMA(lookbackPeriod int) *baseDEMA {
	newDEMA := baseDEMA{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod)}
	return &newDEMA
}

// A Double Exponential Moving Average Indicator
type DEMA struct {
	*baseDEMA

	// public variables
	Data []float64
}

// NewDEMA returns a new Double Exponential Moving Average (DEMA) configured with the
// specified lookbackPeriod. The DEMA results are stored in the DATA field.
func NewDEMA(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *DEMA, err error) {
	newDEMA := DEMA{baseDEMA: newBaseDEMA(2*(lookbackPeriod) - 1)}
	newDEMA.ema1, _ = NewEMA(lookbackPeriod, selectData)

	newDEMA.ema1.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newDEMA.currentEMA = dataItem
		newDEMA.ema2.ReceiveTick(dataItem, streamBarIndex)
	}

	newDEMA.ema2, _ = NewEMA(lookbackPeriod, selectData)

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

	newDEMA.selectData = selectData
	newDEMA.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newDEMA.Data = append(newDEMA.Data, dataItem)
	}
	return &newDEMA, nil
}

func NewDEMAForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *DEMA, err error) {
	newDEMA, err := NewDEMA(lookbackPeriod, selectData)
	priceStream.AddTickSubscription(newDEMA)
	return newDEMA, err
}

func (dema *baseDEMA) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = dema.selectData(tickData)
	dema.ReceiveTick(selectedData, streamBarIndex)
}

func (dema *baseDEMA) ReceiveTick(tickData float64, streamBarIndex int) {
	dema.ema1.ReceiveTick(tickData, streamBarIndex)
}
