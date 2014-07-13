package indicators

// highest high value in period indicator

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Lowest Low Value In Period Indicator
type LLVWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	periodHistory        *list.List
	valueAvailableAction ValueAvailableAction
	currentLow           float64
	currentLowIndex      int
}

// NewLLVWithoutStorage returns a new Lowest Low Value (LLV) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The LLV results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewLLVWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *LLVWithoutStorage, err error) {
	newLLV := LLVWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		currentLow:                  math.MaxFloat64,
		currentLowIndex:             0,
		periodHistory:               list.New()}
	newLLV.selectData = selectData
	newLLV.valueAvailableAction = valueAvailableAction

	return &newLLV, nil
}

// A Lowest Low Value Indicator
type LLV struct {
	*LLVWithoutStorage

	// public variables
	Data []float64
}

// NewLLV returns a new Lowest Low Value (LLV) configured with the
// specified timePeriod. The LLV results are stored in the Data field.
func NewLLV(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LLV, err error) {
	newLLV := LLV{}
	newLLV.LLVWithoutStorage, err = NewLLVWithoutStorage(timePeriod, selectData, func(dataItem float64, streamBarIndex int) {
		newLLV.Data = append(newLLV.Data, dataItem)
	})

	return &newLLV, err
}

func NewLLVForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LLV, err error) {
	newSma, err := NewLLV(timePeriod, selectData)
	priceStream.AddTickSubscription(newSma)
	return newSma, err
}

func (ind *LLVWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *LLVWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)

	if ind.periodHistory.Len() > ind.GetTimePeriod() {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)

		// make sure we haven't just removed the current high
		if ind.currentLowIndex == ind.GetTimePeriod()-1 {
			ind.currentLow = math.MaxFloat64
			// we have we need to find the new high in the history
			var i int = ind.GetTimePeriod() - 1
			for e := ind.periodHistory.Front(); e != nil; e = e.Next() {
				value := e.Value.(float64)
				if value < ind.currentLow {
					ind.currentLow = value
					ind.currentLowIndex = i
				}
				i -= 1
			}
		} else {
			if tickData < ind.currentLow {
				ind.currentLow = tickData
				ind.currentLowIndex = 0
			} else {
				ind.currentLowIndex += 1
			}
		}

		var result = ind.currentLow

		ind.dataLength += 1

		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}

		if result > ind.maxValue {
			ind.maxValue = result
		}

		if result < ind.minValue {
			ind.minValue = result
		}

		ind.valueAvailableAction(result, streamBarIndex)

	} else {
		if tickData < ind.currentLow {
			ind.currentLow = tickData
			ind.currentLowIndex = 0
		} else {
			ind.currentLowIndex += 1
		}

		if ind.periodHistory.Len() == ind.GetTimePeriod() {
			var result = ind.currentLow

			ind.dataLength += 1

			if ind.validFromBar == -1 {
				ind.validFromBar = streamBarIndex
			}

			if result > ind.maxValue {
				ind.maxValue = result
			}

			if result < ind.minValue {
				ind.minValue = result
			}

			ind.valueAvailableAction(result, streamBarIndex)
		}
	}

}
