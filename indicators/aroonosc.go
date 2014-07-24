package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

type AroonOscWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	//private variables
	valueAvailableAction ValueAvailableActionFloat
	aroon                *AroonWithoutStorage
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
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		valueAvailableAction: valueAvailableAction,
	}

	ind.aroon, err = NewAroonWithoutStorage(timePeriod,
		func(dataItemAroonUp float64, dataItemAroonDown float64, streamBarIndex int) {
			// increment the number of results this indicator can be expected to return
			ind.dataLength++

			result := dataItemAroonUp - dataItemAroonDown
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

// NewAroonOscWithKnownSourceLength creates an Aroon Oscillator (AroonOsc) for offline usage
func NewAroonOscWithKnownSourceLength(sourceLength int, timePeriod int) (indicator *AroonOsc, err error) {
	ind, err := NewAroonOsc(timePeriod)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultAroonOscWithKnownSourceLength creates an Aroon Oscillator (AroonOsc) for offline usage with default parameters
func NewDefaultAroonOscWithKnownSourceLength(sourceLength int) (indicator *AroonOsc, err error) {

	ind, err := NewDefaultAroonOsc()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewAroonOscForStream creates an Aroon Oscillator (AroonOsc) for online usage with a source data stream
func NewAroonOscForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *AroonOsc, err error) {
	ind, err := NewAroonOsc(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAroonOscForStream creates an Aroon Oscillator (AroonOsc) for online usage with a source data stream
func NewDefaultAroonOscForStream(priceStream *gotrade.DOHLCVStream) (indicator *AroonOsc, err error) {
	ind, err := NewDefaultAroonOsc()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewAroonOscForStreamWithKnownSourceLength creates an Aroon Oscillator (AroonOsc) for offline usage with a source data stream
func NewAroonOscForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *AroonOsc, err error) {
	ind, err := NewAroonOscWithKnownSourceLength(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAroonOscForStreamWithKnownSourceLength creates an Aroon Oscillator (AroonOsc) for offline usage with a source data stream
func NewDefaultAroonOscForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *AroonOsc, err error) {
	ind, err := NewDefaultAroonOscWithKnownSourceLength(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *AroonOsc) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.aroon.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
