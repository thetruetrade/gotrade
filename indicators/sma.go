// Simple Moving Average (Sma)
package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

// A Simple Moving Average Indicator
type SmaWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	periodTotal          float64
	periodHistory        *list.List
	periodCounter        int
	valueAvailableAction ValueAvailableActionFloat
}

// NewSmaWithoutStorage returns a new Simple Moving Average (Sma) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The Sma results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewSmaWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *SmaWithoutStorage, err error) {
	newSma := SmaWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               timePeriod * -1,
		periodHistory:               list.New()}
	newSma.valueAvailableAction = valueAvailableAction

	return &newSma, nil
}

// A Simple Moving Average Indicator
type Sma struct {
	*SmaWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewSma returns a new Simple Moving Average (Sma) configured with the
// specified timePeriod. The Sma results are stored in the Data field.
func NewSma(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Sma, err error) {
	newSma := Sma{}
	newSma.SmaWithoutStorage, err = NewSmaWithoutStorage(
		timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newSma.Data = append(newSma.Data, dataItem)
		})
	newSma.selectData = selectData

	return &newSma, err
}

func NewSmaForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Sma, err error) {
	newSma, err := NewSma(timePeriod, selectData)
	priceStream.AddTickSubscription(newSma)
	return newSma, err
}

func (sma *Sma) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = sma.selectData(tickData)
	sma.ReceiveTick(selectedData, streamBarIndex)
}

func (sma *SmaWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
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
