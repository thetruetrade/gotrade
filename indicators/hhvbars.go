package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Highest High Value Bars Indicator (HhvBars), no storage, for use in other indicators
type HhvBarsWithoutStorage struct {
	*baseIndicator
	*baseIntBounds

	// private variables
	periodHistory        *list.List
	valueAvailableAction ValueAvailableActionInt
	currentHigh          float64
	currentHighIndex     int64
	timePeriod           int
}

// NewHhvBarsWithoutStorage creates a Highest High Value Bars Indicator Indicator (HhvBars) without storage
func NewHhvBarsWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionInt) (indicator *HhvBarsWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for this indicator is 1
	if timePeriod < 1 {
		return nil, errors.New("timePeriod is less than the minimum (1)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := timePeriod - 1

	ind := HhvBarsWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseIntBounds:        newBaseIntBounds(),
		currentHigh:          math.SmallestNonzeroFloat64,
		currentHighIndex:     0,
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, nil
}

// A Highest High Value Bars Indicator (HhvBars)
type HhvBars struct {
	*HhvBarsWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []int64
}

// NewHhvBars creates a Highest High Value Bars Indicator (HhvBars) for online usage
func NewHhvBars(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *HhvBars, err error) {
	ind := HhvBars{selectData: selectData}
	ind.HhvBarsWithoutStorage, err = NewHhvBarsWithoutStorage(timePeriod, func(dataItem int64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultHhvBars creates a Highest High Value Indicator (HhvBars) for online usage with default parameters
//	- timePeriod: 25
func NewDefaultHhvBars() (indicator *HhvBars, err error) {
	timePeriod := 25
	return NewHhvBars(timePeriod, gotrade.UseClosePrice)
}

// NewHhvBarsWithSrcLen creates a Highest High Value Indicator (HhvBars)for offline usage
func NewHhvBarsWithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *HhvBars, err error) {
	ind, err := NewHhvBars(timePeriod, selectData)
	ind.Data = make([]int64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultHhvBarsWithSrcLen creates a Highest High Value Indicator (HhvBars)for offline usage with default parameters
func NewDefaultHhvBarsWithSrcLen(sourceLength int) (indicator *HhvBars, err error) {
	ind, err := NewDefaultHhvBars()
	ind.Data = make([]int64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewHhvBarsForStream creates a Highest High Value Indicator (HhvBars)for online usage with a source data stream
func NewHhvBarsForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *HhvBars, err error) {
	ind, err := NewHhvBars(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultHhvBarsForStream creates a Highest High Value Indicator (HhvBars)for online usage with a source data stream
func NewDefaultHhvBarsForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *HhvBars, err error) {
	ind, err := NewDefaultHhvBars()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewHhvBarsForStreamWithSrcLen creates a Highest High Value Indicator (HhvBars)for offline usage with a source data stream
func NewHhvBarsForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *HhvBars, err error) {
	ind, err := NewHhvBarsWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultHhvBarsForStreamWithSrcLen creates a Highest High Value Indicator (HhvBars)for offline usage with a source data stream
func NewDefaultHhvBarsForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *HhvBars, err error) {
	ind, err := NewDefaultHhvBarsWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *HhvBars) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *HhvBarsWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)

	// resize the history
	if ind.periodHistory.Len() > ind.timePeriod {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)

		// make sure we haven't just removed the current high
		if ind.currentHighIndex == int64(ind.timePeriod-1) {
			ind.currentHigh = math.SmallestNonzeroFloat64
			// we have we need to find the new high in the history
			var i int = ind.timePeriod - 1
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

		var result = ind.currentHighIndex

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1

		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

		// update the maximum result value
		if result > ind.maxValue {
			ind.maxValue = result
		}

		// update the minimum result value
		if result < ind.minValue {
			ind.minValue = result
		}

		// notify of a new result value though the value available action
		ind.valueAvailableAction(result, streamBarIndex)

	} else {
		if tickData > ind.currentHigh {
			ind.currentHigh = tickData
			ind.currentHighIndex = 0
		} else {
			ind.currentHighIndex += 1
		}

		if ind.periodHistory.Len() == ind.timePeriod {
			var result = ind.currentHighIndex

			// increment the number of results this indicator can be expected to return
			ind.dataLength += 1

			if ind.validFromBar == -1 {
				// set the streamBarIndex from which this indicator returns valid results
				ind.validFromBar = streamBarIndex
			}

			// update the maximum result value
			if result > ind.maxValue {
				ind.maxValue = result
			}

			// update the minimum result value
			if result < ind.minValue {
				ind.minValue = result
			}

			// notify of a new result value though the value available action
			ind.valueAvailableAction(result, streamBarIndex)
		}
	}

}
