package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Chaikin Oscillator Indicator (ChaikinOsc), no storage, for use in other indicators
type ChaikinOscWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	fastTimePeriod       int
	slowTimePeriod       int
	valueAvailableAction ValueAvailableActionFloat
	adl                  *AdlWithoutStorage
	emaFast              float64
	emaSlow              float64
	emaFastMultiplier    float64
	emaSlowMultiplier    float64
	periodCounter        int
	isInitialised        bool
}

// NewChaikinOscWithoutStorage creates a Chaikin Oscillator Indicator (ChaikinOsc) without storage
// This should be as simple as EMA(Adl,3) - EMA(Adl,10), however it seems the TA-Lib emas are intialised with the
// first adl value and not offset like the macd to conincide, they are both calculated from the 2nd bar and used before their
// lookback period is reached - so the emas are calculated inline and not using the general EmaWithoutStorage
func NewChaikinOscWithoutStorage(fastTimePeriod int, slowTimePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *ChaikinOscWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum fastTimePeriod for a Chaikin Oscillator Indicator is 2
	if fastTimePeriod < 2 {
		return nil, errors.New("fastTimePeriod is less than the minimum (2)")
	}

	// the minimum slowTimePeriod for a Chaikin Oscillator Indicator is 2
	if slowTimePeriod < 2 {
		return nil, errors.New("slowTimePeriod is less than the minimum (2)")
	}

	// check the maximum fastTimePeriod
	if fastTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("fastTimePeriod is greater than the maximum (100000)")
	}

	// check the maximum slowTimePeriod
	if slowTimePeriod > MaximumLookbackPeriod {
		return nil, errors.New("slowTimePeriod is greater than the maximum (100000)")
	}

	lookback := slowTimePeriod - 1
	ind := ChaikinOscWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		slowTimePeriod:       slowTimePeriod,
		fastTimePeriod:       fastTimePeriod,
		emaFastMultiplier:    float64(2.0 / float64(fastTimePeriod+1.0)),
		emaSlowMultiplier:    float64(2.0 / float64(slowTimePeriod+1.0)),
		periodCounter:        slowTimePeriod * -1,
		isInitialised:        false,
		valueAvailableAction: valueAvailableAction,
	}

	ind.adl, err = NewAdlWithoutStorage(func(dataItem float64, streamBarIndex int) {
		ind.periodCounter += 1

		if !ind.isInitialised {
			ind.emaFast = dataItem
			ind.emaSlow = dataItem
			ind.isInitialised = true
		}
		if ind.periodCounter < 0 {
			ind.emaFast = (dataItem-ind.emaFast)*ind.emaFastMultiplier + ind.emaFast
			ind.emaSlow = (dataItem-ind.emaSlow)*ind.emaSlowMultiplier + ind.emaSlow
		}

		if ind.periodCounter >= 0 {
			// increment the number of results this indicator can be expected to return
			ind.dataLength += 1

			if ind.validFromBar == -1 {
				// set the streamBarIndex from which this indicator returns valid results
				ind.validFromBar = streamBarIndex
			}

			ind.emaFast = (dataItem-ind.emaFast)*ind.emaFastMultiplier + ind.emaFast
			ind.emaSlow = (dataItem-ind.emaSlow)*ind.emaSlowMultiplier + ind.emaSlow
			chaikinOsc := ind.emaFast - ind.emaSlow

			// update the maximum result value
			if chaikinOsc > ind.maxValue {
				ind.maxValue = chaikinOsc
			}

			// update the minimum result value
			if chaikinOsc < ind.minValue {
				ind.minValue = chaikinOsc
			}

			// notify of a new result value though the value available action
			ind.valueAvailableAction(chaikinOsc, streamBarIndex)
		}
	})

	return &ind, err
}

// A Chaikin Oscillator Indicator (ChaikinOsc)
type ChaikinOsc struct {
	*ChaikinOscWithoutStorage

	// public variables
	Data []float64
}

// NewChaikinOsc creates a Chaikin Oscillator (ChaikinOsc) for online usage
func NewChaikinOsc(fastTimePeriod int, slowTimePeriod int) (indicator *ChaikinOsc, err error) {

	newChaikinOsc := ChaikinOsc{}
	newChaikinOsc.ChaikinOscWithoutStorage, err = NewChaikinOscWithoutStorage(fastTimePeriod, slowTimePeriod,
		func(dataItem float64, streamBarIndex int) {
			newChaikinOsc.Data = append(newChaikinOsc.Data, dataItem)
		})

	return &newChaikinOsc, err
}

// NewDefaultChaikinOsc creates a Chaikin Oscillator (ChaikinOsc) for online usage with default parameters
//	- fastTimePeriod: 3
//  - slowTimePeriod: 10
func NewDefaultChaikinOsc() (indicator *ChaikinOsc, err error) {
	fastTimePeriod := 3
	slowTimePeriod := 10
	return NewChaikinOsc(fastTimePeriod, slowTimePeriod)
}

// NewChaikinOscWithSrcLen creates a Chaikin Oscillator (ChaikinOsc) for offline usage
func NewChaikinOscWithSrcLen(sourceLength int, fastTimePeriod int, slowTimePeriod int) (indicator *ChaikinOsc, err error) {
	ind, err := NewChaikinOsc(fastTimePeriod, slowTimePeriod)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultChaikinOscWithSrcLen creates a Chaikin Oscillator (ChaikinOsc) for offline usage with default parameters
func NewDefaultChaikinOscWithSrcLen(sourceLength int) (indicator *ChaikinOsc, err error) {
	ind, err := NewDefaultChaikinOsc()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewChaikinOscForStream creates a Chaikin Oscillator (ChaikinOsc) for online usage with a source data stream
func NewChaikinOscForStream(priceStream *gotrade.DOHLCVStream, fastTimePeriod int, slowTimePeriod int) (indicator *ChaikinOsc, err error) {
	newChaikinOsc, err := NewChaikinOsc(fastTimePeriod, slowTimePeriod)
	priceStream.AddTickSubscription(newChaikinOsc)
	return newChaikinOsc, err
}

// NewChaikinOscForStreamWithSrcLen creates a Chaikin Oscillator (ChaikinOsc) for offline usage with a source data stream
func NewChaikinOscForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream, fastTimePeriod int, slowTimePeriod int) (indicator *ChaikinOsc, err error) {
	ind, err := NewChaikinOscWithSrcLen(sourceLength, fastTimePeriod, slowTimePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultChaikinOscForStreamWithSrcLen creates a Chaikin Oscillator (ChaikinOsc) for offline usage with a source data stream
func NewDefaultChaikinOscForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *ChaikinOsc, err error) {
	ind, err := NewDefaultChaikinOscWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *ChaikinOsc) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.adl.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
