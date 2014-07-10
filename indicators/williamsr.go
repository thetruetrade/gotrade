package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Williamns Percent R Indicator
type WILLRWithoutStorage struct {
	*baseIndicator
	*baseIndicatorWithTimePeriod

	// private variables
	periodHighHistory    *list.List
	periodLowHistory     *list.List
	periodCounter        int
	valueAvailableAction ValueAvailableAction
}

// NewWILLRWithoutStorage returns a new Simple Moving Average (WILLR) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The WILLR results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewWILLRWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableAction) (indicator *WILLRWithoutStorage, err error) {
	newWILLR := WILLRWithoutStorage{baseIndicator: newBaseIndicator(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               timePeriod * -1,
		periodHighHistory:           list.New(),
		periodLowHistory:            list.New()}
	newWILLR.valueAvailableAction = valueAvailableAction

	return &newWILLR, nil
}

// A Simple Moving Average Indicator
type WILLR struct {
	*WILLRWithoutStorage

	// public variables
	Data []float64
}

// NewWILLR returns a new Williamns Percent R  (WILLR) configured with the
// specified timePeriod. The WILLR results are stored in the Data field.
func NewWILLR(timePeriod int) (indicator *WILLR, err error) {
	newWILLR := WILLR{}
	newWILLR.WILLRWithoutStorage, err = NewWILLRWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		newWILLR.Data = append(newWILLR.Data, dataItem)
	})

	return &newWILLR, err
}

func NewWILLRForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *WILLR, err error) {
	newInd, err := NewWILLR(timePeriod)
	priceStream.AddTickSubscription(newInd)
	return newInd, err
}

func (ind *WILLRWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	ind.periodCounter += 1
	ind.periodHighHistory.PushBack(tickData.H())
	ind.periodLowHistory.PushBack(tickData.L())

	highestHigh, _ := highestHighofPeriod(ind.periodHighHistory)
	lowestLow, _ := lowestLowofPeriod(ind.periodLowHistory)

	var result float64 = (highestHigh - tickData.C()) / (highestHigh - lowestLow) * -100.0
	if ind.periodCounter >= 0 {
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

	if ind.periodHighHistory.Len() >= ind.GetTimePeriod() {
		var first = ind.periodHighHistory.Front()
		ind.periodHighHistory.Remove(first)
	}
	if ind.periodLowHistory.Len() >= ind.GetTimePeriod() {
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
