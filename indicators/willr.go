package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Williamns Percent R Indicator
type WillRWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodHighHistory *list.List
	periodLowHistory  *list.List
	periodCounter     int
	timePeriod        int
}

// NewWillRWithoutStorage creates a Williams Percent R Indicator (WillR) without storage
func NewWillRWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *WillRWithoutStorage, err error) {

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

	lookback := timePeriod - 1
	ind := WillRWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		periodCounter:                timePeriod * -1,
		periodHighHistory:            list.New(),
		periodLowHistory:             list.New(),
		timePeriod:                   timePeriod,
	}

	return &ind, nil
}

// A Simple Moving Average Indicator
type WillR struct {
	*WillRWithoutStorage

	// public variables
	Data []float64
}

// NewWillR creates a Williams Percent R Indicator (WillR) for online usage
func NewWillR(timePeriod int) (indicator *WillR, err error) {
	ind := WillR{}
	ind.WillRWithoutStorage, err = NewWillRWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultWillR creates a Williams Percent R Indicator (WillR) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultWillR() (indicator *WillR, err error) {
	timePeriod := 14
	return NewWillR(timePeriod)
}

// NewWillRWithSrcLen creates a Williams Percent R Indicator (WillR) for offline usage
func NewWillRWithSrcLen(sourceLength uint, timePeriod int) (indicator *WillR, err error) {
	ind, err := NewWillR(timePeriod)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultWillRWithSrcLen creates a Williams Percent R Indicator (WillR) for offline usage with default parameters
func NewDefaultWillRWithSrcLen(sourceLength uint) (indicator *WillR, err error) {
	ind, err := NewDefaultWillR()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewWillRForStream creates a Williams Percent R Indicator (WillR) for online usage with a source data stream
func NewWillRForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *WillR, err error) {
	ind, err := NewWillR(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultWillRForStream creates a Williams Percent R Indicator (WillR) for online usage with a source data stream
func NewDefaultWillRForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *WillR, err error) {
	ind, err := NewDefaultWillR()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewWillRForStreamWithSrcLen creates a Williams Percent R Indicator (WillR) for offline usage with a source data stream
func NewWillRForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *WillR, err error) {
	ind, err := NewWillRWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultWillRForStreamWithSrcLen creates a Williams Percent R Indicator (WillR) for offline usage with a source data stream
func NewDefaultWillRForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *WillR, err error) {
	ind, err := NewDefaultWillRWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *WillRWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	ind.periodCounter += 1
	ind.periodHighHistory.PushBack(tickData.H())
	ind.periodLowHistory.PushBack(tickData.L())

	highestHigh, _ := highestHighofPeriod(ind.periodHighHistory)
	lowestLow, _ := lowestLowofPeriod(ind.periodLowHistory)

	var result float64 = (highestHigh - tickData.C()) / (highestHigh - lowestLow) * -100.0
	if ind.periodCounter >= 0 {

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
	}

	if ind.periodHighHistory.Len() >= ind.timePeriod {
		var first = ind.periodHighHistory.Front()
		ind.periodHighHistory.Remove(first)
	}
	if ind.periodLowHistory.Len() >= ind.timePeriod {
		var first = ind.periodLowHistory.Front()
		ind.periodLowHistory.Remove(first)
	}
}

func highestHighofPeriod(l *list.List) (result float64, err error) {
	if l.Len() == 0 {
		err = errors.New("list is empty no high can be calculated.")
	}

	high := math.SmallestNonzeroFloat64
	for e := l.Front(); e != nil; e = e.Next() {
		value := e.Value.(float64)
		if value > high {
			high = value
		}
	}
	return high, err
}

func lowestLowofPeriod(l *list.List) (result float64, err error) {
	if l.Len() == 0 {
		err = errors.New("list is empty no low can be calculated.")
	}

	low := math.MaxFloat64
	for e := l.Front(); e != nil; e = e.Next() {
		value := e.Value.(float64)
		if value < low {
			low = value
		}
	}

	return low, err
}
