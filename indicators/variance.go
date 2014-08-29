package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Variance Indicator (Var), no storage, for use in other indicators
type VarWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodCounter int
	periodHistory *list.List
	mean          float64
	variance      float64
	timePeriod    int
}

// NewVarWithoutStorage creates a Variance Indicator (Var) without storage
func NewVarWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *VarWithoutStorage, err error) {

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

	ind := VarWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		periodCounter:                0,
		periodHistory:                list.New(),
		mean:                         0.0,
		variance:                     0.0,
		timePeriod:                   timePeriod,
	}

	return &ind, nil
}

// A Variance Indicator (Var)
type Var struct {
	*VarWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewVar creates a Variance Indicator (Var) for online usage
func NewVar(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Var, err error) {
	ind := Var{selectData: selectData}
	ind.VarWithoutStorage, err = NewVarWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultVar creates a Variance Indicator (Var) for online usage with default parameters
//	- timePeriod: 10
func NewDefaultVar() (indicator *Var, err error) {
	timePeriod := 10
	return NewVar(timePeriod, gotrade.UseClosePrice)
}

// NewVarWithSrcLen creates a Variance Indicator (Var) for offline usage
func NewVarWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Var, err error) {
	ind, err := NewVar(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultVarWithSrcLen creates a Variance Indicator (Var) for offline usage with default parameters
func NewDefaultVarWithSrcLen(sourceLength uint) (indicator *Var, err error) {
	ind, err := NewDefaultVar()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewVarForStream creates a Variance Indicator (Var) for online usage with a source data stream
func NewVarForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Var, err error) {
	ind, err := NewVar(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultVarForStream creates a Variance Indicator (Var) for online usage with a source data stream
func NewDefaultVarForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Var, err error) {
	ind, err := NewDefaultVar()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewVarForStreamWithSrcLen creates a Variance Indicator (Var) for offline usage with a source data stream
func NewVarForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Var, err error) {
	ind, err := NewVarWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultVarForStreamWithSrcLen creates a Variance Indicator (Var) for offline usage with a source data stream
func NewDefaultVarForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Var, err error) {
	ind, err := NewDefaultVarWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Var) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

// http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance - Knuth
func (ind *VarWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)
	firstValue := ind.periodHistory.Front().Value.(float64)

	previousMean := ind.mean
	previousVar := ind.variance

	if ind.periodCounter < ind.timePeriod {
		ind.periodCounter += 1
		delta := tickData - previousMean
		ind.mean = previousMean + delta/float64(ind.periodCounter)

		ind.variance = previousVar + delta*(tickData-ind.mean)
	} else {
		delta := tickData - firstValue
		dOld := firstValue - previousMean
		ind.mean = previousMean + delta/float64(ind.periodCounter)
		dNew := tickData - ind.mean
		ind.variance = previousVar + (dOld+dNew)*(delta)
	}

	if ind.periodHistory.Len() > ind.timePeriod {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}

	if ind.periodCounter >= ind.timePeriod {

		result := ind.variance / float64(ind.timePeriod)

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
	}
}
