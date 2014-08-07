package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Highest High Value Indicator (Hhv), no storage, for use in other indicators
type HhvWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	periodHistory        *list.List
	valueAvailableAction ValueAvailableActionFloat
	currentHigh          float64
	currentHighIndex     int
	timePeriod           int
}

// NewHhvWithoutStorage creates a Highest High Value Indicator Indicator (Hhv) without storage
func NewHhvWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *HhvWithoutStorage, err error) {

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
	ind := HhvWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		currentHigh:          math.SmallestNonzeroFloat64,
		currentHighIndex:     0,
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, nil
}

// A Highest High Value Indicator (hhv)
type Hhv struct {
	*HhvWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewHhv creates a Highest High Value Indicator (Hhv) for online usage
func NewHhv(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Hhv, err error) {
	ind := Hhv{selectData: selectData}
	ind.HhvWithoutStorage, err = NewHhvWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultHhv creates a Highest High Value Indicator (Hhv) for online usage with default parameters
//	- timePeriod: 25
func NewDefaultHhv() (indicator *Hhv, err error) {
	timePeriod := 25
	return NewHhv(timePeriod, gotrade.UseClosePrice)
}

// NewHhvWithSrcLen creates a Highest High Value Indicator (Hhv)for offline usage
func NewHhvWithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Hhv, err error) {
	ind, err := NewHhv(timePeriod, selectData)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultHhvWithSrcLen creates a Highest High Value Indicator (Hhv)for offline usage with default parameters
func NewDefaultHhvWithSrcLen(sourceLength int) (indicator *Hhv, err error) {
	ind, err := NewDefaultHhv()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewHhvForStream creates a Highest High Value Indicator (Hhv)for online usage with a source data stream
func NewHhvForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Hhv, err error) {
	ind, err := NewHhv(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultHhvForStream creates a Highest High Value Indicator (Hhv)for online usage with a source data stream
func NewDefaultHhvForStream(priceStream *gotrade.DOHLCVStream) (indicator *Hhv, err error) {
	ind, err := NewDefaultHhv()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewHhvForStreamWithSrcLen creates a Highest High Value Indicator (Hhv)for offline usage with a source data stream
func NewHhvForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Hhv, err error) {
	ind, err := NewHhvWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultHhvForStreamWithSrcLen creates a Highest High Value Indicator (Hhv)for offline usage with a source data stream
func NewDefaultHhvForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Hhv, err error) {
	ind, err := NewDefaultHhvWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Hhv) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *HhvWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)

	// resize the history
	if ind.periodHistory.Len() > ind.timePeriod {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)

		// make sure we haven't just removed the current high
		if ind.currentHighIndex == ind.timePeriod-1 {
			ind.currentHigh = math.SmallestNonzeroFloat64
			// we have we need to find the new high in the history
			var i int = ind.timePeriod - 1
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
		var result = ind.currentHigh

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
			var result = ind.currentHigh

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
