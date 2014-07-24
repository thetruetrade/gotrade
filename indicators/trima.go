package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type TRIMAWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	sma1                 *SmaWithoutStorage
	sma2                 *SmaWithoutStorage
	currentSma           float64
}

func NewTRIMAWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *TRIMAWithoutStorage, err error) {
	newTRIMA := TRIMAWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod)}

	newTRIMA.valueAvailableAction = valueAvailableAction

	var sma1Period int
	var sma2Period int

	if timePeriod%2 == 0 {
		// even
		sma1Period = timePeriod / 2
		sma2Period = (timePeriod / 2) + 1
	} else {
		// odd
		sma1Period = (timePeriod + 1) / 2
		sma2Period = (timePeriod + 1) / 2
	}

	newTRIMA.sma1, err = NewSmaWithoutStorage(sma1Period, func(dataItem float64, streamBarIndex int) {
		newTRIMA.currentSma = dataItem
		newTRIMA.sma2.ReceiveTick(dataItem, streamBarIndex)
	})

	newTRIMA.sma2, _ = NewSmaWithoutStorage(sma2Period, func(dataItem float64, streamBarIndex int) {
		newTRIMA.dataLength += 1
		if newTRIMA.validFromBar == -1 {
			newTRIMA.validFromBar = streamBarIndex
		}

		result := dataItem

		if result > newTRIMA.maxValue {
			newTRIMA.maxValue = result
		}

		if result < newTRIMA.minValue {
			newTRIMA.minValue = result
		}

		newTRIMA.valueAvailableAction(result, streamBarIndex)
	})

	return &newTRIMA, err
}

// A Triangular Moving Average Indicator
type TRIMA struct {
	*TRIMAWithoutStorage
	selectData gotrade.DataSelectionFunc
	// public variables
	Data []float64
}

// NewTRIMA returns a new TriangularMoving Average (TRIMA) configured with the
// specified timePeriod. The TRIMA results are stored in the DATA field.
func NewTRIMA(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *TRIMA, err error) {
	newTRIMA := TRIMA{selectData: selectData}
	newTRIMA.TRIMAWithoutStorage, err = NewTRIMAWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newTRIMA.Data = append(newTRIMA.Data, dataItem)
		})
	return &newTRIMA, err
}

func NewTRIMAForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *TRIMA, err error) {
	newTRIMA, err := NewTRIMA(timePeriod, selectData)
	priceStream.AddTickSubscription(newTRIMA)
	return newTRIMA, err
}

func (tema *TRIMA) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = tema.selectData(tickData)
	tema.ReceiveTick(selectedData, streamBarIndex)
}

func (tema *TRIMAWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	tema.sma1.ReceiveTick(tickData, streamBarIndex)
}
