package indicators

// highest high value in period indicator

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Highest High Value In Period Indicator
type HHVBarsWithoutStorage struct {
	*baseIndicatorWithIntBounds
	*baseIndicatorWithTimePeriod

	// private variables
	periodHistory        *list.List
	valueAvailableAction ValueAvailableActionInt
	currentHigh          float64
	currentHighIndex     int64
}

// NewHHVBarsWithoutStorage returns a new Highest High Value (HHVBars) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The HHVBars results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewHHVBarsWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableActionInt) (indicator *HHVBarsWithoutStorage, err error) {
	newHHVBars := HHVBarsWithoutStorage{baseIndicatorWithIntBounds: newBaseIndicatorWithIntBounds(0),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		currentHigh:                 math.SmallestNonzeroFloat64,
		currentHighIndex:            0,
		periodHistory:               list.New()}
	newHHVBars.selectData = selectData
	newHHVBars.valueAvailableAction = valueAvailableAction

	return &newHHVBars, nil
}

// A Highest High Value Indicator
type HHVBars struct {
	*HHVBarsWithoutStorage

	// public variables
	Data []int64
}

// NewHHVBars returns a new Highest High Value (HHVBars) configured with the
// specified timePeriod. The HHVBars results are stored in the Data field.
func NewHHVBars(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *HHVBars, err error) {
	newHHVBars := HHVBars{}
	newHHVBars.HHVBarsWithoutStorage, err = NewHHVBarsWithoutStorage(timePeriod, selectData, func(dataItem int64, streamBarIndex int) {
		newHHVBars.Data = append(newHHVBars.Data, dataItem)
	})

	return &newHHVBars, err
}

func NewHHVBarsForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *HHVBars, err error) {
	newSma, err := NewHHVBars(timePeriod, selectData)
	priceStream.AddTickSubscription(newSma)
	return newSma, err
}

func (ind *HHVBarsWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *HHVBarsWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)

	if ind.periodHistory.Len() > ind.GetTimePeriod() {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)

		// make sure we haven't just removed the current high
		if ind.currentHighIndex == int64(ind.GetTimePeriod()-1) {
			ind.currentHigh = math.SmallestNonzeroFloat64
			// we have we need to find the new high in the history
			var i int = ind.GetTimePeriod() - 1
			for e := ind.periodHistory.Front(); e != nil; e = e.Next() {
				value := e.Value.(float64)
				if value > ind.currentHigh {
					ind.currentHigh = value
					ind.currentHighIndex = int64(i)
				}
				i -= 1
			}
		} else {
			if tickData > ind.currentHigh {
				ind.currentHigh = tickData
				ind.currentHighIndex = 0
			} else {
				ind.currentHighIndex += 1
			}
		}

	} else {
		if tickData > ind.currentHigh {
			ind.currentHigh = tickData
			ind.currentHighIndex = 0
		} else {
			ind.currentHighIndex += 1
		}
	}

	var result = ind.currentHighIndex

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
