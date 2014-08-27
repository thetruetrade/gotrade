package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

type AroonOscWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	//private variables
	aroon *AroonWithoutStorage
}

func NewAroonOscWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *AroonOscWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for an AroonOsc indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod is less than the minimum (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := timePeriod
	ind := AroonOscWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
	}

	ind.aroon, err = NewAroonWithoutStorage(timePeriod,
		func(dataItemAroonUp float64, dataItemAroonDown float64, streamBarIndex int) {

			result := dataItemAroonUp - dataItemAroonDown

			ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
		})
	return &ind, nil
}

// An AroonOsc (AroonOsc)
type AroonOsc struct {
	*AroonOscWithoutStorage

	// public variables
	Data []float64
}

// NewAroonOsc creates an Aroon Oscillator (AroonOsc) for online usage
func NewAroonOsc(timePeriod int) (indicator *AroonOsc, err error) {

	ind := AroonOsc{}
	ind.AroonOscWithoutStorage, err = NewAroonOscWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultAroonOsc creates an Aroon Oscillator (AroonOsc) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultAroonOsc() (indicator *AroonOsc, err error) {
	timePeriod := 14
	return NewAroonOsc(timePeriod)
}

// NewAroonOscWithSrcLen creates an Aroon Oscillator (AroonOsc) for offline usage
func NewAroonOscWithSrcLen(sourceLength uint, timePeriod int) (indicator *AroonOsc, err error) {
	ind, err := NewAroonOsc(timePeriod)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultAroonOscWithSrcLen creates an Aroon Oscillator (AroonOsc) for offline usage with default parameters
func NewDefaultAroonOscWithSrcLen(sourceLength uint) (indicator *AroonOsc, err error) {
	ind, err := NewDefaultAroonOsc()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewAroonOscForStream creates an Aroon Oscillator (AroonOsc) for online usage with a source data stream
func NewAroonOscForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *AroonOsc, err error) {
	ind, err := NewAroonOsc(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAroonOscForStream creates an Aroon Oscillator (AroonOsc) for online usage with a source data stream
func NewDefaultAroonOscForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *AroonOsc, err error) {
	ind, err := NewDefaultAroonOsc()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewAroonOscForStreamWithSrcLen creates an Aroon Oscillator (AroonOsc) for offline usage with a source data stream
func NewAroonOscForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *AroonOsc, err error) {
	ind, err := NewAroonOscWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAroonOscForStreamWithSrcLen creates an Aroon Oscillator (AroonOsc) for offline usage with a source data stream
func NewDefaultAroonOscForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *AroonOsc, err error) {
	ind, err := NewDefaultAroonOscWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *AroonOsc) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.aroon.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
