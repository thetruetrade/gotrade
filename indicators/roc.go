package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Rate of Change Indicator (Roc), no storage, for use in other indicators
type RocWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	periodHistory        *list.List
	timePeriod           int
}

// NewRocWithoutStorage creates a Rate of Change Indicator (Roc) without storage
func NewRocWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *RocWithoutStorage, err error) {

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
	ind := RocWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        (timePeriod * -1),
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, nil
}

// A Rate of Change Indicator (Roc)
type Roc struct {
	*RocWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewRoc creates a Rate of Change Indicator (Roc) for online usage
func NewRoc(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Roc, err error) {
	ind := Roc{selectData: selectData}
	ind.RocWithoutStorage, err = NewRocWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	ind.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	}
	return &ind, err
}

// NewDefaultRoc creates a Rate of Change Indicator (Roc) for online usage with default parameters
//	- timePeriod: 10
func NewDefaultRoc() (indicator *Roc, err error) {
	timePeriod := 10
	return NewRoc(timePeriod, gotrade.UseClosePrice)
}

// NewRocWithSrcLen creates a Rate of Change Indicator (Roc) for offline usage
func NewRocWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Roc, err error) {
	ind, err := NewRoc(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultRocWithSrcLen creates a Rate of Change Indicator (Roc) for offline usage with default parameters
func NewDefaultRocWithSrcLen(sourceLength uint) (indicator *Roc, err error) {
	ind, err := NewDefaultRoc()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewRocForStream creates a Rate of Change Indicator (Roc) for online usage with a source data stream
func NewRocForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Roc, err error) {
	ind, err := NewRoc(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultRocForStream creates a Rate of Change Indicator (Roc) for online usage with a source data stream
func NewDefaultRocForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Roc, err error) {
	ind, err := NewDefaultRoc()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewRocForStreamWithSrcLen creates a Rate of Change Indicator (Roc) for offline usage with a source data stream
func NewRocForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Roc, err error) {
	ind, err := NewRocWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultRocForStreamWithSrcLen creates a Rate of Change Indicator (Roc) for offline usage with a source data stream
func NewDefaultRocForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Roc, err error) {
	ind, err := NewDefaultRocWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Roc) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *RocWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

		//    Roc = (price/previousPrice - 1) * 100
		previousPrice := ind.periodHistory.Front().Value.(float64)

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}
		var result float64
		if previousPrice != 0 {
			result = 100.0 * ((tickData / previousPrice) - 1)
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
