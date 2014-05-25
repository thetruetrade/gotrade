// Simple Moving Average (SMA)
package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

// A Simple Moving Average Indicator
type SMAWithoutStorage struct {
	*baseIndicatorWithLookback

	// private variables
	periodTotal          float64
	periodHistory        *list.List
	periodCounter        int
	valueAvailableAction ValueAvailableAction
}

// NewSMAWithoutStorage returns a new Simple Moving Average (SMA) configured with the
// specified lookbackPeriod, this version is intended for use by other indicators.
// The SMA results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewSMAWithoutStorage(lookbackPeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *SMAWithoutStorage, err error) {
	newSMA := SMAWithoutStorage{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod),
		periodCounter: lookbackPeriod * -1,
		periodHistory: list.New()}
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
// specified lookbackPeriod. The SMA results are stored in the Data field.
func NewSMA(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *SMA, err error) {
	newSMA := SMA{}
	newSMA.SMAWithoutStorage, err = NewSMAWithoutStorage(lookbackPeriod, selectData, func(dataItem float64, streamBarIndex int) {
		newSMA.Data = append(newSMA.Data, dataItem)
	})

	return &newSMA, err
}

func NewSMAForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *SMA, err error) {
	newSma, err := NewSMA(lookbackPeriod, selectData)
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
	if sma.periodHistory.Len() > sma.lookbackPeriod {
		var first = sma.periodHistory.Front()
		sma.periodHistory.Remove(first)
	}
	sma.periodTotal += tickData
	var result float64 = sma.periodTotal / float64(sma.lookbackPeriod)
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
