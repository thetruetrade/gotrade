// Triple Exponential Moving Average (TEMA)
package indicators

// TEMA(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))

import (
	"github.com/thetruetrade/gotrade"
)

type TEMAWithoutStorage struct {
	*baseIndicatorWithLookback

	// private variables
	valueAvailableAction ValueAvailableAction
	ema1                 *EMA
	ema2                 *EMA
	ema3                 *EMA
	currentEMA           float64
	currentEMA2          float64
}

func NewTEMAWithoutStorage(lookbackPeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *TEMAWithoutStorage, err error) {
	newTEMA := TEMAWithoutStorage{baseIndicatorWithLookback: newBaseIndicatorWithLookback(3*lookbackPeriod - 2)}
	newTEMA.selectData = selectData
	newTEMA.valueAvailableAction = valueAvailableAction

	newTEMA.ema1, err = NewEMA(lookbackPeriod, selectData)

	newTEMA.ema1.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newTEMA.currentEMA = dataItem
		newTEMA.ema2.ReceiveTick(dataItem, streamBarIndex)
	}

	newTEMA.ema2, _ = NewEMA(lookbackPeriod, selectData)
	newTEMA.ema2.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newTEMA.currentEMA2 = dataItem
		newTEMA.ema3.ReceiveTick(dataItem, streamBarIndex)
	}

	newTEMA.ema3, _ = NewEMA(lookbackPeriod, selectData)

	newTEMA.ema3.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newTEMA.dataLength += 1
		if newTEMA.validFromBar == -1 {
			newTEMA.validFromBar = streamBarIndex
		}

		//T-EMA = (3*EMA – 3*EMA(EMA)) + EMA(EMA(EMA))
		tema := (3*newTEMA.currentEMA - 3*newTEMA.currentEMA2) + dataItem

		if tema > newTEMA.maxValue {
			newTEMA.maxValue = tema
		}

		if tema < newTEMA.minValue {
			newTEMA.minValue = tema
		}

		newTEMA.valueAvailableAction(tema, streamBarIndex)
	}

	return &newTEMA, err
}

// A Double Exponential Moving Average Indicator
type TEMA struct {
	*TEMAWithoutStorage

	// public variables
	Data []float64
}

// NewTEMA returns a new Double Exponential Moving Average (TEMA) configured with the
// specified lookbackPeriod. The TEMA results are stored in the DATA field.
func NewTEMA(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *TEMA, err error) {
	newTEMA := TEMA{}
	newTEMA.TEMAWithoutStorage, err = NewTEMAWithoutStorage(lookbackPeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newTEMA.Data = append(newTEMA.Data, dataItem)
		})
	return &newTEMA, err
}

func NewTEMAForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *TEMA, err error) {
	newTEMA, err := NewTEMA(lookbackPeriod, selectData)
	priceStream.AddTickSubscription(newTEMA)
	return newTEMA, err
}

func (tema *TEMAWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = tema.selectData(tickData)
	tema.ReceiveTick(selectedData, streamBarIndex)
}

func (tema *TEMAWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	tema.ema1.ReceiveTick(tickData, streamBarIndex)
}
