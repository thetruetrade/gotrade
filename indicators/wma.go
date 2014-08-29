package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Weighted Moving Average Indicator (Wma), no storage, for use in other indicators
type WmaWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodTotal       float64
	periodHistory     *list.List
	periodCounter     int
	periodWeightTotal int
	timePeriod        int
}

// NewWmaWithoutStorage creates a Weighted Moving Average Indicator (Wma) without storage
func NewWmaWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *WmaWithoutStorage, err error) {

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
	ind := WmaWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		periodCounter:                timePeriod * -1,
		periodHistory:                list.New(),
		timePeriod:                   timePeriod,
	}

	var weightedTotal int = 0
	for i := 1; i <= timePeriod; i++ {
		weightedTotal += i
	}
	ind.periodWeightTotal = weightedTotal

	return &ind, nil
}

// A Weighted Moving Average Indicator (Wma)
type Wma struct {
	*WmaWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewWma creates a Weighted Moving Average Indicator (Wma) for online usage
func NewWma(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Wma, err error) {
	ind := Wma{selectData: selectData}
	ind.WmaWithoutStorage, err = NewWmaWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})
	return &ind, err
}

// NewDefaultWma creates a Weighted Moving Average Indicator (Wma) for online usage with default parameters
//	- timePeriod: 10
func NewDefaultWma() (indicator *Wma, err error) {
	timePeriod := 10
	return NewWma(timePeriod, gotrade.UseClosePrice)
}

// NewWmaWithSrcLen creates a Weighted Moving Average Indicator (Wma) for offline usage
func NewWmaWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Wma, err error) {
	ind, err := NewWma(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultWmaWithSrcLen creates a Weighted Moving Average Indicator (Wma) for offline usage with default parameters
func NewDefaultWmaWithSrcLen(sourceLength uint) (indicator *Wma, err error) {
	ind, err := NewDefaultWma()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewWmaForStream creates a Weighted Moving Average Indicator (Wma) for online usage with a source data stream
func NewWmaForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Wma, err error) {
	ind, err := NewWma(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultWmaForStream creates a Weighted Moving Average Indicator (Wma) for online usage with a source data stream
func NewDefaultWmaForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Wma, err error) {
	ind, err := NewDefaultWma()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewWmaForStreamWithSrcLen creates a Weighted Moving Average Indicator (Wma) for offline usage with a source data stream
func NewWmaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Wma, err error) {
	ind, err := NewWmaWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultWmaForStreamWithSrcLen creates a Weighted Moving Average Indicator (Wma) for offline usage with a source data stream
func NewDefaultWmaForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Wma, err error) {
	ind, err := NewDefaultWmaWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Wma) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *WmaWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1

	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

	}
	if ind.periodHistory.Len() > ind.timePeriod {
		var first = ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}

	if ind.periodCounter >= 0 {
		// calculate the ind
		var iter int = 1
		var sum float64 = 0
		for e := ind.periodHistory.Front(); e != nil; e = e.Next() {
			var localSum float64 = 0
			for i := 1; i <= iter; i++ {
				localSum += e.Value.(float64)
			}
			sum += localSum
			iter++
		}
		var result float64 = sum / float64(ind.periodWeightTotal)

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
	}
}
