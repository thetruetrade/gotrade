package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Rate of Change Percentage Indicator (RocP), no storage, for use in other indicators
type RocPWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	periodHistory        *list.List
	timePeriod           int
}

// NewRocPWithoutStorage creates a Rate of Change Percentage Indicator (RocP) without storage
func NewRocPWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *RocPWithoutStorage, err error) {

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

	lookback := timePeriod
	ind := RocPWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        (timePeriod * -1),
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, err
}

// A Rate of Change Percentage Indicator (RocP)
type RocP struct {
	*RocPWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewRocP creates a Rate of Change Percentage Indicator (RocP) for online usage
func NewRocP(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RocP, err error) {
	ind := RocP{selectData: selectData}
	ind.RocPWithoutStorage, err = NewRocPWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	ind.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	}
	return &ind, err
}

// NewDefaultRocP creates a Rate of Change Percentage Indicator (RocP) for online usage with default parameters
//	- timePeriod: 10
func NewDefaultRocP() (indicator *RocP, err error) {
	timePeriod := 10
	return NewRocP(timePeriod, gotrade.UseClosePrice)
}

// NewRocPWithSrcLen creates a Rate of Change Percentage Indicator (RocP) for offline usage
func NewRocPWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RocP, err error) {
	ind, err := NewRocP(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultRocPWithSrcLen creates a Rate of Change Percentage Indicator (RocP) for offline usage with default parameters
func NewDefaultRocPWithSrcLen(sourceLength uint) (indicator *RocP, err error) {
	ind, err := NewDefaultRocP()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewRocPForStream creates a Rate of Change Percentage Indicator (RocP) for online usage with a source data stream
func NewRocPForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RocP, err error) {
	ind, err := NewRocP(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultRocPForStream creates a Rate of Change Percentage Indicator (RocP) for online usage with a source data stream
func NewDefaultRocPForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *RocP, err error) {
	ind, err := NewDefaultRocP()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewRocPForStreamWithSrcLen creates a Rate of Change Percentage Indicator (RocP) for offline usage with a source data stream
func NewRocPForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RocP, err error) {
	ind, err := NewRocPWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultRocPForStreamWithSrcLen creates a Rate of Change Percentage Indicator (RocP) for offline usage with a source data stream
func NewDefaultRocPForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *RocP, err error) {
	ind, err := NewDefaultRocPWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *RocP) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *RocPWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

		//    RocP = (price/previousPrice - 1) * 100
		previousPrice := ind.periodHistory.Front().Value.(float64)

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}
		var result float64
		if previousPrice != 0 {
			result = (tickData - previousPrice) / previousPrice
		} else {
			result = 0.0
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

	if ind.periodHistory.Len() > ind.timePeriod {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}
}
