package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Lowest Low Value Bars Indicator (LlvBars), no storage, for use in other indicators
type LlvBarsWithoutStorage struct {
	*baseIndicator
	*baseIntBounds

	// private variables
	periodHistory        *list.List
	valueAvailableAction ValueAvailableActionInt
	currentLow           float64
	currentLowIndex      int64
	timePeriod           int
}

// NewLlvBarsWithoutStorage creates a Lowest Low Value Bars Indicator Indicator (LlvBars) without storage
func NewLlvBarsWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionInt) (indicator *LlvBarsWithoutStorage, err error) {

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

	ind := LlvBarsWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseIntBounds:        newBaseIntBounds(),
		currentLow:           math.MaxFloat64,
		currentLowIndex:      0,
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, nil
}

// A Lowest Low Value Bars Indicator (LlvBars)
type LlvBars struct {
	*LlvBarsWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []int64
}

// NewLlvBars creates a Lowest Low Value Bars Indicator (LlvBars) for online usage
func NewLlvBars(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LlvBars, err error) {
	ind := LlvBars{selectData: selectData}
	ind.LlvBarsWithoutStorage, err = NewLlvBarsWithoutStorage(timePeriod, func(dataItem int64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultLlvBars creates a Lowest Low Value Indicator (LlvBars) for online usage with default parameters
//	- timePeriod: 25
func NewDefaultLlvBars() (indicator *LlvBars, err error) {
	timePeriod := 25
	return NewLlvBars(timePeriod, gotrade.UseClosePrice)
}

// NewLlvBarsWithSrcLen creates a Lowest Low Value Indicator (LlvBars)for offline usage
func NewLlvBarsWithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LlvBars, err error) {
	ind, err := NewLlvBars(timePeriod, selectData)
	ind.Data = make([]int64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultLlvBarsWithSrcLen creates a Lowest Low Value Indicator (LlvBars)for offline usage with default parameters
func NewDefaultLlvBarsWithSrcLen(sourceLength int) (indicator *LlvBars, err error) {
	ind, err := NewDefaultLlvBars()
	ind.Data = make([]int64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewLlvBarsForStream creates a Lowest Low Value Indicator (LlvBars)for online usage with a source data stream
func NewLlvBarsForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LlvBars, err error) {
	ind, err := NewLlvBars(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLlvBarsForStream creates a Lowest Low Value Indicator (LlvBars)for online usage with a source data stream
func NewDefaultLlvBarsForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *LlvBars, err error) {
	ind, err := NewDefaultLlvBars()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewLlvBarsForStreamWithSrcLen creates a Lowest Low Value Indicator (LlvBars)for offline usage with a source data stream
func NewLlvBarsForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LlvBars, err error) {
	ind, err := NewLlvBarsWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLlvBarsForStreamWithSrcLen creates a Lowest Low Value Indicator (LlvBars)for offline usage with a source data stream
func NewDefaultLlvBarsForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *LlvBars, err error) {
	ind, err := NewDefaultLlvBarsWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *LlvBars) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *LlvBarsWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)

	// resize the history
	if ind.periodHistory.Len() > ind.timePeriod {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)

		// make sure we haven't just removed the current low
		if ind.currentLowIndex == int64(ind.timePeriod-1) {
			ind.currentLow = math.MaxFloat64
			// we have we need to find the new low in the history
			var i int = ind.timePeriod - 1
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
		if tickData < ind.currentLow {
			ind.currentLow = tickData
			ind.currentLowIndex = 0
		} else {
			ind.currentLowIndex += 1
		}

		if ind.periodHistory.Len() == ind.timePeriod {
			var result = ind.currentLowIndex

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
