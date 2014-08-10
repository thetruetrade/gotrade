// Simple Moving Average (Sma)
package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Simple Moving Average Indicator (Sma), no storage, for use in other indicators
type SmaWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	periodTotal          float64
	periodHistory        *list.List
	periodCounter        int
	valueAvailableAction ValueAvailableActionFloat
	timePeriod           int
}

// NewSmaWithoutStorage creates a Simple Moving Average Indicator (Sma) without storage
func NewSmaWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *SmaWithoutStorage, err error) {

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
	ind := SmaWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        timePeriod * -1,
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	return &ind, nil
}

// A Simple Moving Average Indicator (Sma)
type Sma struct {
	*SmaWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewSma creates a Simple Moving Average Indicator (Sma) for online usage
func NewSma(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Sma, err error) {
	ind := Sma{}
	ind.SmaWithoutStorage, err = NewSmaWithoutStorage(
		timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})
	ind.selectData = selectData

	return &ind, err
}

// NewDefaultSma creates a Simple Moving Average Indicator (Sma) for online usage with default parameters
//	- timePeriod: 10
func NewDefaultSma() (indicator *Sma, err error) {
	timePeriod := 10
	return NewSma(timePeriod, gotrade.UseClosePrice)
}

// NewSmaWithSrcLen creates a Simple Moving Average Indicator (Sma) for offline usage
func NewSmaWithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Sma, err error) {
	ind, err := NewSma(timePeriod, selectData)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultSmaWithSrcLen creates a Simple Moving Average Indicator (Sma) for offline usage with default parameters
func NewDefaultSmaWithSrcLen(sourceLength int) (indicator *Sma, err error) {
	ind, err := NewDefaultSma()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewSmaForStream creates a Simple Moving Average Indicator (Sma) for online usage with a source data stream
func NewSmaForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Sma, err error) {
	ind, err := NewSma(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultSmaForStream creates a Simple Moving Average Indicator (Sma) for online usage with a source data stream
func NewDefaultSmaForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Sma, err error) {
	ind, err := NewDefaultSma()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewSmaForStreamWithSrcLen creates a Simple Moving Average Indicator (Sma) for offline usage with a source data stream
func NewSmaForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Sma, err error) {
	ind, err := NewSmaWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultSmaForStreamWithSrcLen creates a Simple Moving Average Indicator (Sma) for offline usage with a source data stream
func NewDefaultSmaForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Sma, err error) {
	ind, err := NewDefaultSmaWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Sma) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *SmaWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {
		var valueToRemove = ind.periodHistory.Front()
		ind.periodTotal -= valueToRemove.Value.(float64)
	}
	if ind.periodHistory.Len() > ind.timePeriod {
		var first = ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}
	ind.periodTotal += tickData
	var result float64 = ind.periodTotal / float64(ind.timePeriod)
	if ind.periodCounter >= 0 {

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1

		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
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
}
