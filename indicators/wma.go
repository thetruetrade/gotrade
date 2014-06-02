// Weighted Moving Average (WMA)
package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type WMAWithoutStorage struct {
	*baseIndicator
	*baseIndicatorWithTimePeriod

	// private variables
	periodTotal          float64
	periodHistory        *list.List
	periodCounter        int
	periodWeightTotal    int
	valueAvailableAction ValueAvailableAction
}

// NewAttachedWMA returns a new Simple Moving Average (WMA) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The WMA results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewWMAWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *WMAWithoutStorage, err error) {
	newWMA := WMAWithoutStorage{baseIndicator: newBaseIndicator(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               timePeriod * -1,
		periodHistory:               list.New()}

	var weightedTotal int = 0
	for i := 1; i <= timePeriod; i++ {
		weightedTotal += i
	}
	newWMA.periodWeightTotal = weightedTotal

	newWMA.selectData = selectData
	newWMA.valueAvailableAction = valueAvailableAction
	return &newWMA, nil
}

// A Simple Moving Average Indicator
type WMA struct {
	*WMAWithoutStorage

	// public variables
	Data []float64
}

// NewWMA returns a new Simple Moving Average (WMA) configured with the
// specified timePeriod. The WMA results are stored in the DATA field.
func NewWMA(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *WMA, err error) {
	newWMA := WMA{}
	newWMA.WMAWithoutStorage, err = NewWMAWithoutStorage(timePeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newWMA.Data = append(newWMA.Data, dataItem)
		})
	return &newWMA, err
}

func NewWMAForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *WMA, err error) {
	newWma, err := NewWMA(timePeriod, selectData)
	priceStream.AddTickSubscription(newWma)
	return newWma, err
}

func (wma *WMAWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = wma.selectData(tickData)
	wma.ReceiveTick(selectedData, streamBarIndex)
}

func (wma *WMAWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	wma.periodCounter += 1

	wma.periodHistory.PushBack(tickData)

	if wma.periodCounter > 0 {

	}
	if wma.periodHistory.Len() > wma.GetTimePeriod() {
		var first = wma.periodHistory.Front()
		wma.periodHistory.Remove(first)
	}

	if wma.periodCounter >= 0 {
		wma.dataLength += 1
		if wma.validFromBar == -1 {
			wma.validFromBar = streamBarIndex
		}

		// calculate the wma
		var iter int = 1
		var sum float64 = 0
		for e := wma.periodHistory.Front(); e != nil; e = e.Next() {
			var localSum float64 = 0
			for i := 1; i <= iter; i++ {
				localSum += e.Value.(float64)
			}
			sum += localSum
			iter++
		}
		var result float64 = sum / float64(wma.periodWeightTotal)

		if result > wma.maxValue {
			wma.maxValue = result
		}

		if result < wma.minValue {
			wma.minValue = result
		}

		wma.valueAvailableAction(result, streamBarIndex)
	}
}
