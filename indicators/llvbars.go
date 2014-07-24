package indicators

// lowest low value in period indicator

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Lowest Low Value In Period Indicator
type LLVBarsWithoutStorage struct {
	*baseIndicatorWithIntBounds
	*baseIndicatorWithTimePeriod

	// private variables
	periodHistory        *list.List
	valueAvailableAction ValueAvailableActionInt
	currentLow           float64
	currentLowIndex      int64
}

// NewLLVBarsWithoutStorage returns a new Lowest Low Value (LLVBars) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The LLVBars results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewLLVBarsWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionInt) (indicator *LLVBarsWithoutStorage, err error) {
	newLLVBars := LLVBarsWithoutStorage{baseIndicatorWithIntBounds: newBaseIndicatorWithIntBounds(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		currentLow:                  math.MaxFloat64,
		currentLowIndex:             0,
		periodHistory:               list.New()}
	newLLVBars.valueAvailableAction = valueAvailableAction

	return &newLLVBars, nil
}

// A Lowest Low Value Indicator
type LLVBars struct {
	*LLVBarsWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []int64
}

// NewLLVBars returns a new Lowest Low Value (LLVBars) configured with the
// specified timePeriod. The LLVBars results are stored in the Data field.
func NewLLVBars(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LLVBars, err error) {
	newLLVBars := LLVBars{selectData: selectData}
	newLLVBars.LLVBarsWithoutStorage, err = NewLLVBarsWithoutStorage(timePeriod, func(dataItem int64, streamBarIndex int) {
		newLLVBars.Data = append(newLLVBars.Data, dataItem)
	})

	return &newLLVBars, err
}

func NewLLVBarsForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LLVBars, err error) {
	newSma, err := NewLLVBars(timePeriod, selectData)
	priceStream.AddTickSubscription(newSma)
	return newSma, err
}

func (ind *LLVBars) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *LLVBarsWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)

	if ind.periodHistory.Len() > ind.GetTimePeriod() {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)

		// make sure we haven't just removed the current low
		if ind.currentLowIndex == int64(ind.GetTimePeriod()-1) {
			ind.currentLow = math.MaxFloat64
			// we have we need to find the new low in the history
			var i int = ind.GetTimePeriod() - 1
			for e := ind.periodHistory.Front(); e != nil; e = e.Next() {
				value := e.Value.(float64)
				if value < ind.currentLow {
					ind.currentLow = value
					ind.currentLowIndex = int64(i)
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

		var result = ind.currentLowIndex

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
			var result = ind.currentLowIndex

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
