package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Rate of Change Ratio 100 Scale Indicator (RocR100), no storage, for use in other indicators
type RocR100WithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	periodHistory        *list.List
	timePeriod           int
}

// NewRocR100WithoutStorage creates a Rate of Change Ratio 100 Scale Indicator (RocR100) without storage
func NewRocR100WithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *RocR100WithoutStorage, err error) {

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
	ind := RocR100WithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        (timePeriod * -1),
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, err
}

// A Rate of Change Ratio Indicator (RocR)
type RocR100 struct {
	*RocR100WithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewRocR100 creates a Rate of Change Ratio 100 Scale Indicator (RocR100) for online usage
func NewRocR100(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RocR100, err error) {
	newRocR100 := RocR100{selectData: selectData}
	newRocR100.RocR100WithoutStorage, err = NewRocR100WithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newRocR100.Data = append(newRocR100.Data, dataItem)
		})

	newRocR100.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newRocR100.Data = append(newRocR100.Data, dataItem)
	}
	return &newRocR100, err
}

/// NewDefaultRocR100 creates a Rate of Change Ratio 100 Scale Indicator (RocR100) for online usage with default parameters
//	- timePeriod: 10
func NewDefaultRocR100() (indicator *RocR100, err error) {
	timePeriod := 10
	return NewRocR100(timePeriod, gotrade.UseClosePrice)
}

// NewRocR100WithSrcLen creates a Rate of Change Ratio 100 Scale Indicator (RocR100) for offline usage
func NewRocR100WithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RocR100, err error) {
	ind, err := NewRocR100(timePeriod, selectData)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultRocR100WithSrcLen creates a Rate of Change Ratio 100 Scale Indicator (RocR100) for offline usage with default parameters
func NewDefaultRocR100WithSrcLen(sourceLength int) (indicator *RocR100, err error) {
	ind, err := NewDefaultRocR100()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewRocR100ForStream creates a Rate of Change Ratio 100 Scale Indicator (RocR100) for online usage with a source data stream
func NewRocR100ForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RocR100, err error) {
	ind, err := NewRocR100(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultRocR100ForStream creates a Rate of Change Ratio 100 Scale Indicator (RocR100) for online usage with a source data stream
func NewDefaultRocR100ForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *RocR100, err error) {
	ind, err := NewDefaultRocR100()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewRocR100ForStreamWithSrcLen creates a Rate of Change Ratio 100 Scale Indicator (RocR100) for offline usage with a source data stream
func NewRocR100ForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RocR100, err error) {
	ind, err := NewRocR100WithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultRocR100ForStreamWithSrcLen creates a Rate of Change Ratio 100 Scale Indicator (RocR100) for offline usage with a source data stream
func NewDefaultRocR100ForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *RocR100, err error) {
	ind, err := NewDefaultRocR100WithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *RocR100) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *RocR100WithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

		//    RocR100 = (price/previousPrice - 1) * 100
		previousPrice := ind.periodHistory.Front().Value.(float64)

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}
		var result float64
		if previousPrice != 0 {
			result = (tickData / previousPrice) * 100.0
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
