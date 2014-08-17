package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Linear Regression Indicator (LinReg), no storage, for use in other indicators
type LinRegWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	periodCounter        int
	periodHistory        *list.List
	sumX                 float64
	sumXSquare           float64
	divisor              float64
	valueAvailableAction ValueAvailableActionLinearReg
	timePeriod           int
}

// NewLinRegWithoutStorage creates a Linear Regression Indicator (LinReg) without storage
func NewLinRegWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionLinearReg) (indicator *LinRegWithoutStorage, err error) {

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
	ind := LinRegWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        (timePeriod) * -1,
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	timePeriodF := float64(timePeriod)
	timePeriodFMinusOne := timePeriodF - 1.0
	ind.sumX = timePeriodF * timePeriodFMinusOne * 0.5
	ind.sumXSquare = timePeriodF * timePeriodFMinusOne * (2.0*timePeriodF - 1.0) / 6.0
	ind.divisor = ind.sumX*ind.sumX - timePeriodF*ind.sumXSquare

	ind.valueAvailableAction = valueAvailableAction

	return &ind, nil
}

// A Linear Regression Indicator (LinReg)
type LinReg struct {
	*LinRegWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewLinReg creates a Linear Regression Indicator (LinReg) for online usage
func NewLinReg(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinReg, err error) {
	ind := LinReg{selectData: selectData}
	ind.LinRegWithoutStorage, err = NewLinRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)

			// update the maximum result value
			if dataItem > ind.LinRegWithoutStorage.maxValue {
				ind.LinRegWithoutStorage.maxValue = dataItem
			}

			// update the minimum result value
			if dataItem < ind.LinRegWithoutStorage.minValue {
				ind.LinRegWithoutStorage.minValue = dataItem
			}
		})

	return &ind, err
}

// NewDefaultLinReg creates a Linear Regression Indicator (LinReg) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultLinReg() (indicator *LinReg, err error) {
	timePeriod := 14
	return NewLinReg(timePeriod, gotrade.UseClosePrice)
}

// NewLinRegWithSrcLen creates a Linear Regression Indicator (LinReg) for offline usage
func NewLinRegWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinReg, err error) {
	ind, err := NewLinReg(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultLinRegWithSrcLen creates a Linear Regression Indicator (LinReg) for offline usage with default parameters
func NewDefaultLinRegWithSrcLen(sourceLength uint) (indicator *LinReg, err error) {
	ind, err := NewDefaultLinReg()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewLinRegForStream creates a Linear Regression Indicator (LinReg) for online usage with a source data stream
func NewLinRegForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinReg, err error) {
	ind, err := NewLinReg(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLinRegForStream creates a Linear Regression Indicator (LinReg) for online usage with a source data stream
func NewDefaultLinRegForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *LinReg, err error) {
	ind, err := NewDefaultLinReg()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewLinRegForStreamWithSrcLen creates a Linear Regression Indicator (LinReg) for offline usage with a source data stream
func NewLinRegForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinReg, err error) {
	ind, err := NewLinRegWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLinRegForStreamWithSrcLen creates a Linear Regression Indicator (LinReg) for offline usage with a source data stream
func NewDefaultLinRegForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *LinReg, err error) {
	ind, err := NewDefaultLinRegWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *LinReg) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *LinRegWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1

	if ind.periodCounter >= 0 {
		sumXY := 0.0
		sumY := 0.0
		i := ind.timePeriod
		var value float64 = 0.0
		for e := ind.periodHistory.Front(); e != nil; e = e.Next() {
			i--
			value = e.Value.(float64)
			sumY += value
			sumXY += (float64(i) * value)
		}
		sumY += tickData
		timePeriod := float64(ind.timePeriod)
		m := (timePeriod*sumXY - ind.sumX*sumY) / ind.divisor
		b := (sumY - m*ind.sumX) / timePeriod
		result := b + m*float64(timePeriod-1.0)

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

		// notify of a new result value though the value available action
		ind.valueAvailableAction(result, m, b, streamBarIndex)
	}

	ind.periodHistory.PushBack(tickData)

	if ind.periodHistory.Len() >= ind.timePeriod {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}

}
