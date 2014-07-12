// Simple Moving Average (SMA)
package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

// A Simple Moving Average Indicator
type SMAWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	periodTotal          float64
	periodHistory        *list.List
	periodCounter        int
	valueAvailableAction ValueAvailableAction
}

// NewSMAWithoutStorage returns a new Simple Moving Average (SMA) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The SMA results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewSMAWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *SMAWithoutStorage, err error) {
	newSMA := SMAWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               timePeriod * -1,
		periodHistory:               list.New()}
	newSMA.selectData = selectData
	newSMA.valueAvailableAction = valueAvailableAction

	return &newSMA, nil
}

// A Simple Moving Average Indicator
type SMA struct {
	*SMAWithoutStorage

	// public variables
	Data []float64
}

// NewSMA returns a new Simple Moving Average (SMA) configured with the
// specified timePeriod. The SMA results are stored in the Data field.
func NewSMA(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *SMA, err error) {
	newSMA := SMA{}
	newSMA.SMAWithoutStorage, err = NewSMAWithoutStorage(timePeriod, selectData, func(dataItem float64, streamBarIndex int) {
		newSMA.Data = append(newSMA.Data, dataItem)
	})

	return &newSMA, err
}

func NewSMAForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *SMA, err error) {
	newSma, err := NewSMA(timePeriod, selectData)
	priceStream.AddTickSubscription(newSma)
	return newSma, err
}

func (sma *SMAWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = sma.selectData(tickData)
	sma.ReceiveTick(selectedData, streamBarIndex)
}

func (sma *SMAWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	sma.periodCounter += 1
	sma.periodHistory.PushBack(tickData)

	if sma.periodCounter > 0 {
		var valueToRemove = sma.periodHistory.Front()
		sma.periodTotal -= valueToRemove.Value.(float64)
	}
	if sma.periodHistory.Len() > sma.GetTimePeriod() {
		var first = sma.periodHistory.Front()
		sma.periodHistory.Remove(first)
	}
	sma.periodTotal += tickData
	var result float64 = sma.periodTotal / float64(sma.GetTimePeriod())
	if sma.periodCounter >= 0 {
		sma.dataLength += 1

		if sma.validFromBar == -1 {
			sma.validFromBar = streamBarIndex
		}

		if result > sma.maxValue {
			sma.maxValue = result
		}

		if result < sma.minValue {
			sma.minValue = result
		}

		sma.valueAvailableAction(result, streamBarIndex)
	}
}
