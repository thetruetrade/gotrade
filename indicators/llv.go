package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Lowest Low Value Indicator (Llv), no storage, for use in other indicators
type LlvWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	periodHistory        *list.List
	valueAvailableAction ValueAvailableActionFloat
	currentLow           float64
	currentLowIndex      int
	timePeriod           int
}

// NewLlvWithoutStorage creates a Lowest Low Value Indicator Indicator (Llv) without storage
func NewLlvWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *LlvWithoutStorage, err error) {

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
	ind := LlvWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		currentLow:           math.MaxFloat64,
		currentLowIndex:      0,
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, nil
}

// A Lowest Low Value Indicator (hhv)
type Llv struct {
	*LlvWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewLlv creates a Lowest Low Value Indicator (Llv) for online usage
func NewLlv(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Llv, err error) {
	ind := Llv{selectData: selectData}
	ind.LlvWithoutStorage, err = NewLlvWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultLlv creates a Lowest Low Value Indicator (Llv) for online usage with default parameters
//	- timePeriod: 25
func NewDefaultLlv() (indicator *Llv, err error) {
	timePeriod := 25
	return NewLlv(timePeriod, gotrade.UseClosePrice)
}

// NewLlvWithSrcLen creates a Lowest Low Value Indicator (Llv)for offline usage
func NewLlvWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Llv, err error) {
	ind, err := NewLlv(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultLlvWithSrcLen creates a Lowest Low Value Indicator (Llv)for offline usage with default parameters
func NewDefaultLlvWithSrcLen(sourceLength uint) (indicator *Llv, err error) {
	ind, err := NewDefaultLlv()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewLlvForStream creates a Lowest Low Value Indicator (Llv)for online usage with a source data stream
func NewLlvForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Llv, err error) {
	ind, err := NewLlv(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLlvForStream creates a Lowest Low Value Indicator (Llv)for online usage with a source data stream
func NewDefaultLlvForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Llv, err error) {
	ind, err := NewDefaultLlv()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewLlvForStreamWithSrcLen creates a Lowest Low Value Indicator (Llv)for offline usage with a source data stream
func NewLlvForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Llv, err error) {
	ind, err := NewLlvWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLlvForStreamWithSrcLen creates a Lowest Low Value Indicator (Llv)for offline usage with a source data stream
func NewDefaultLlvForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Llv, err error) {
	ind, err := NewDefaultLlvWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Llv) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *LlvWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)

	// resize the history
	if ind.periodHistory.Len() > ind.timePeriod {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)

		// make sure we haven't just removed the current low
		if ind.currentLowIndex == ind.timePeriod-1 {
			ind.currentLow = math.MaxFloat64
			// we have we need to find the new low in the history
			var i int = ind.timePeriod - 1
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
			var result = ind.currentLow

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
