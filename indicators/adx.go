package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// An Average Directional Index (ADX), no storage
type ADXWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	dx                   *DXWithoutStorage
	currentDX            float64
	sumDX                float64
	previousADX          float64
	timePeriod           int
}

// NewADXWithoutStorage creates an Average Directional Index (ADX) without storage
func NewADXWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *ADXWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for an ADX indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod is less than the minimum (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := (2 * timePeriod) - 1
	ind := ADXWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		timePeriod:           timePeriod,
		periodCounter:        timePeriod * -1,
		currentDX:            0.0,
		sumDX:                0.0,
		previousADX:          0.0,
		valueAvailableAction: valueAvailableAction,
	}

	ind.dx, _ = NewDXWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {

		ind.currentDX = dataItem

		ind.periodCounter += 1
		if ind.periodCounter < 0 {
			ind.sumDX += ind.currentDX
		} else if ind.periodCounter == 0 {
			ind.dataLength += 1

			if ind.validFromBar == -1 {
				ind.validFromBar = streamBarIndex
			}

			ind.sumDX += ind.currentDX
			result := ind.sumDX / float64(ind.GetTimePeriod())
			if result > ind.maxValue {
				ind.maxValue = result
			}

			if result < ind.minValue {
				ind.minValue = result
			}
			ind.valueAvailableAction(result, streamBarIndex)
			ind.previousADX = result

		} else {

			ind.dataLength += 1

			result := (ind.previousADX*float64(ind.GetTimePeriod()-1) + ind.currentDX) / float64(ind.GetTimePeriod())
			if result > ind.maxValue {
				ind.maxValue = result
			}

			if result < ind.minValue {
				ind.minValue = result
			}
			ind.valueAvailableAction(result, streamBarIndex)
			ind.previousADX = result
		}

	})

	return &ind, nil
}

// A Directional Movement Indicator (ADX)
type ADX struct {
	*ADXWithoutStorage

	// public variables
	Data []float64
}

// NewADX creates an Average Directional Index (ADX) for online usage
func NewADX(timePeriod int) (indicator *ADX, err error) {

	ind := ADX{}
	ind.ADXWithoutStorage, err = NewADXWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultADX creates an Average Directional Index (ADX) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultADX() (indicator *ADX, err error) {
	timePeriod := 14
	return NewADX(timePeriod)
}

// NewADXWithKnownSourceLength creates an Average Directional Index (ADX) for offline usage
func NewADXWithKnownSourceLength(sourceLength int, timePeriod int) (indicator *ADX, err error) {
	ind, err := NewADXWithKnownSourceLength(sourceLength, timePeriod)
	ind.Data = make([]float64, 0, sourceLength)

	return ind, err
}

// NewDefaultADXWithKnownSourceLength creates an Average Directional Index (ADX) for offline usage with default parameters
func NewDefaultADXWithKnownSourceLength(sourceLength int, timePeriod int) (indicator *ADX, err error) {

	ind, err := NewDefaultADX()
	ind.Data = make([]float64, 0, sourceLength)
	return ind, err
}

// NewADXForStream creates an Average Directional Index (ADX) for online usage with a source data stream
func NewADXForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *ADX, err error) {
	ind, err := NewADX(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewADefaultDXForStream creates an Average Directional Index (ADX) for online usage with a source data stream
func NewDefaultADXForStream(priceStream *gotrade.DOHLCVStream) (indicator *ADX, err error) {
	ind, err := NewDefaultADX()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewADXForStreamWithKnownSourceLength creates an Average Directional Index (ADX) for offline usage with a source data stream
func NewADXForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *ADX, err error) {
	ind, err := NewADXWithKnownSourceLength(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultADXForStreamWithKnownSourceLength creates an Average Directional Index (ADX) for offline usage with a source data stream
func NewDefaultADXForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *ADX, err error) {
	ind, err := NewDefaultADXWithKnownSourceLength(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// GetTimePeriod returns the configured ADX timePeriod
func (ind *ADXWithoutStorage) GetTimePeriod() int {
	return ind.timePeriod
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *ADXWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.dx.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
