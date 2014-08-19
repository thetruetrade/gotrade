package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Momentum Indicator (Mom), no storage, for use in other indicators
type MomWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	periodHistory        *list.List
	timePeriod           int
}

// NewMomWithoutStorage creates a Momentum Indicator (Mom) without storage
func NewMomWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *MomWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for this indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod is less than the minimum (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := timePeriod
	ind := MomWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        (timePeriod * -1),
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, err
}

// A Momentum Indicator (Mom)
type Mom struct {
	*MomWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewMom creates a Momentum (Mom) for online usage
func NewMom(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Mom, err error) {
	ind := Mom{selectData: selectData}
	ind.MomWithoutStorage, err = NewMomWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultMom creates a Momentum (Mom) for online usage with default parameters
//	- timePeriod: 10
//  - selectData: useClosePrice
func NewDefaultMom() (indicator *Mom, err error) {
	timePeriod := 10
	selectData := gotrade.UseClosePrice
	return NewMom(timePeriod, selectData)
}

// NewMomWithSrcLen creates a Momentum (Mom) for offline usage
func NewMomWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Mom, err error) {
	ind, err := NewMom(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultMomWithSrcLen creates a Momentum (Mom) for offline usage with default parameters
func NewDefaultMomWithSrcLen(sourceLength uint) (indicator *Mom, err error) {
	ind, err := NewDefaultMom()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewMomForStream creates a Momentum (Mom) for online usage with a source data stream
func NewMomForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Mom, err error) {
	newMom, err := NewMom(timePeriod, selectData)
	priceStream.AddTickSubscription(newMom)
	return newMom, err
}

// NewDefaultMomForStream creates a Momentum (Mom) for online usage with a source data stream
func NewDefaultMomForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Mom, err error) {
	ind, err := NewDefaultMom()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewMomForStreamWithSrcLen creates a Momentum (Mom) for offline usage with a source data stream
func NewMomForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Mom, err error) {
	ind, err := NewMomWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMomForStreamWithSrcLen creates a Momentum (Mom) for offline usage with a source data stream
func NewDefaultMomForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Mom, err error) {
	ind, err := NewDefaultMomWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Mom) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *MomWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

		// Mom = price - previousPrice
		previousPrice := ind.periodHistory.Front().Value.(float64)

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}
		var result float64 = tickData - previousPrice

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

	if ind.periodHistory.Len() > ind.timePeriod {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}
}
