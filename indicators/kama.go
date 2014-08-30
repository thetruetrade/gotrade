package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Kaufman Adaptive Moving Average Indicator (Kama), no storage, for use in other indicators
type KamaWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodTotal   float64
	periodHistory *list.List
	periodCounter int
	constantMax   float64
	constantDiff  float64
	sumROC        float64
	periodROC     float64
	previousClose float64
	previousKama  float64
	timePeriod    int
}

// NewKamaWithoutStorage creates a Kaufman Adaptive Moving Average Indicator (Kama) without storage
func NewKamaWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *KamaWithoutStorage, err error) {

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
	ind := KamaWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		periodCounter:                (timePeriod + 1) * -1,
		constantMax:                  float64(2.0 / (30.0 + 1.0)),
		constantDiff:                 float64((2.0 / (2.0 + 1.0)) - (2.0 / (30.0 + 1.0))),
		sumROC:                       0.0,
		periodROC:                    0.0,
		periodHistory:                list.New(),
		previousClose:                math.SmallestNonzeroFloat64,
		timePeriod:                   timePeriod,
	}

	return &ind, nil
}

// A Kaufman Adaptive Moving Average Indicator (Kama)
type Kama struct {
	*KamaWithoutStorage
	selectData gotrade.DOHLCVDataSelectionFunc

	// public variables
	Data []float64
}

// NewKama creates a Kaufman Adaptive Moving Average Indicator (Kama) for online usage
func NewKama(timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Kama, err error) {
	if selectData == nil {
		return nil, ErrDOHLCVDataSelectFuncIsNil
	}

	ind := Kama{
		selectData: selectData,
	}

	ind.KamaWithoutStorage, err = NewKamaWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultKama creates a Kaufman Adaptive Moving Average Indicator (Kama) for online usage with default parameters
//	- timePeriod: 25
func NewDefaultKama() (indicator *Kama, err error) {
	timePeriod := 25
	return NewKama(timePeriod, gotrade.UseClosePrice)
}

// NewKamaWithSrcLen creates a Kaufman Adaptive Moving Average Indicator (Kama) for offline usage
func NewKamaWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Kama, err error) {
	ind, err := NewKama(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultKamaWithSrcLen creates a Kaufman Adaptive Moving Average Indicator (Kama) for offline usage with default parameters
func NewDefaultKamaWithSrcLen(sourceLength uint) (indicator *Kama, err error) {
	ind, err := NewDefaultKama()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewKamaForStream creates a Kaufman Adaptive Moving Average Indicator (Kama) for online usage with a source data stream
func NewKamaForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Kama, err error) {
	ind, err := NewKama(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultKamaForStream creates a Kaufman Adaptive Moving Average Indicator (Kama) for online usage with a source data stream
func NewDefaultKamaForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Kama, err error) {
	ind, err := NewDefaultKama()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewKamaForStreamWithSrcLen creates a Kaufman Adaptive Moving Average Indicator (Kama) for offline usage with a source data stream
func NewKamaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *Kama, err error) {
	ind, err := NewKamaWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultKamaForStreamWithSrcLen creates a Kaufman Adaptive Moving Average Indicator (Kama) for offline usage with a source data stream
func NewDefaultKamaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Kama, err error) {
	ind, err := NewDefaultKamaWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Kama) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *KamaWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter <= 0 {
		if ind.previousClose > math.SmallestNonzeroFloat64 {
			ind.sumROC += math.Abs(tickData - ind.previousClose)
		}
	}
	if ind.periodCounter == 0 {
		var er float64 = 0.0
		var sc float64 = 0.0
		var closeMinusN float64 = ind.periodHistory.Front().Value.(float64)
		ind.previousKama = ind.previousClose
		ind.periodROC = tickData - closeMinusN

		// calculate the efficiency ratio
		if ind.sumROC <= ind.periodROC || isZero(ind.sumROC) {
			er = 1.0
		} else {
			er = math.Abs(ind.periodROC / ind.sumROC)
		}

		sc = (er * ind.constantDiff) + ind.constantMax
		sc *= sc
		ind.previousKama = ((tickData - ind.previousKama) * sc) + ind.previousKama

		result := ind.previousKama

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)

	} else if ind.periodCounter > 0 {

		var er float64 = 0.0
		var sc float64 = 0.0
		var closeMinusN float64 = ind.periodHistory.Front().Value.(float64)
		var closeMinusN1 float64 = ind.periodHistory.Front().Next().Value.(float64)
		ind.periodROC = tickData - closeMinusN1

		ind.sumROC -= math.Abs(closeMinusN1 - closeMinusN)
		ind.sumROC += math.Abs(tickData - ind.previousClose)

		// calculate the efficiency ratio
		if ind.sumROC <= ind.periodROC || isZero(ind.sumROC) {
			er = 1.0
		} else {
			er = math.Abs(ind.periodROC / ind.sumROC)
		}

		sc = (er * ind.constantDiff) + ind.constantMax
		sc *= sc
		ind.previousKama = ((tickData - ind.previousKama) * sc) + ind.previousKama

		result := ind.previousKama

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
	}

	ind.previousClose = tickData

	if ind.periodHistory.Len() > (ind.timePeriod + 1) {
		var first = ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}

}

func isZero(value float64) bool {
	var epsilon float64 = 0.00000000000001
	return (((-epsilon) < value) && (value < epsilon))
}
