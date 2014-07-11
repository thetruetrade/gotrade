package indicators

// highest high value in period indicator

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Highest High Value In Period Indicator
type HHVWithoutStorage struct {
	*baseIndicator
	*baseIndicatorWithTimePeriod

	// private variables
	periodHistory        *list.List
	valueAvailableAction ValueAvailableAction
	currentHigh          float64
	currentHighIndex     int
}

// NewHHVWithoutStorage returns a new Highest High Value (HHV) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The HHV results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewHHVWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *HHVWithoutStorage, err error) {
	newHHV := HHVWithoutStorage{baseIndicator: newBaseIndicator(0),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		currentHigh:                 math.SmallestNonzeroFloat64,
		currentHighIndex:            0,
		periodHistory:               list.New()}
	newHHV.selectData = selectData
	newHHV.valueAvailableAction = valueAvailableAction

	return &newHHV, nil
}

// A Highest High Value Indicator
type HHV struct {
	*HHVWithoutStorage

	// public variables
	Data []float64
}

// NewHHV returns a new Highest High Value (HHV) configured with the
// specified timePeriod. The HHV results are stored in the Data field.
func NewHHV(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *HHV, err error) {
	newHHV := HHV{}
	newHHV.HHVWithoutStorage, err = NewHHVWithoutStorage(timePeriod, selectData, func(dataItem float64, streamBarIndex int) {
		newHHV.Data = append(newHHV.Data, dataItem)
	})

	return &newHHV, err
}

func NewHHVForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *HHV, err error) {
	newSma, err := NewHHV(timePeriod, selectData)
	priceStream.AddTickSubscription(newSma)
	return newSma, err
}

func (ind *HHVWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *HHVWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)

	if ind.periodHistory.Len() > ind.GetTimePeriod() {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)

		// make sure we haven't just removed the current high
		if ind.currentHighIndex == ind.GetTimePeriod()-1 {
			ind.currentHigh = math.SmallestNonzeroFloat64
			// we have we need to find the new high in the history
			var i int = ind.GetTimePeriod() - 1
			for e := ind.periodHistory.Front(); e != nil; e = e.Next() {
				value := e.Value.(float64)
				if value > ind.currentHigh {
					ind.currentHigh = value
					ind.currentHighIndex = i
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

	var result = ind.currentHigh

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
