package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// An Average Directional Index Rating (ADXR), no storage
type ADXRWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	periodCounter        int
	periodHistory        *list.List
	adx                  *ADXWithoutStorage
	valueAvailableAction ValueAvailableActionFloat
	timePeriod           int
}

// NewADXRWithoutStorage creates an Average Directional Index Rating (ADXR) without storage
func NewADXRWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *ADXRWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// // the minimum timeperiod for an ADX indicator is 2
	// if timePeriod < 2 {
	// 	return nil, errors.New("timePeriod is less than the minimum (2)")
	// }

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	ind := ADXRWithoutStorage{
		timePeriod:           timePeriod,
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        0,
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
	}

	ind.adx, err = NewADXWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.periodHistory.PushBack(dataItem)

		if ind.periodCounter > ind.GetLookbackPeriod() {
			adxN := ind.periodHistory.Front().Value.(float64)
			result := (dataItem + adxN) / 2.0

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

		if ind.periodHistory.Len() >= ind.adx.GetTimePeriod() {
			first := ind.periodHistory.Front()
			ind.periodHistory.Remove(first)
		}
	})

	var lookback int = 3
	if timePeriod > 1 {
		lookback = timePeriod - 1 + ind.adx.GetLookbackPeriod()
	}
	ind.baseIndicator = newBaseIndicator(lookback)

	return &ind, nil
}

// A Directional Movement Indicator Rating (ADXR)
type ADXR struct {
	*ADXRWithoutStorage

	// public variables
	Data []float64
}

// NewADXR creates an Average Directional Index Rating (ADXR) for online usage
func NewADXR(timePeriod int) (indicator *ADXR, err error) {
	ind := ADXR{}
	ind.ADXRWithoutStorage, err = NewADXRWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultADXR creates an Average Directional Index Rating (ADXR) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultADXR() (indicator *ADXR, err error) {
	timePeriod := 14
	return NewADXR(timePeriod)
}

// NewADXRWithKnownSourceLength creates an Average Directional Index Rating (ADXR) for offline usage
func NewADXRWithKnownSourceLength(sourceLength int, timePeriod int) (indicator *ADXR, err error) {
	ind, err := NewADXR(timePeriod)
	ind.Data = make([]float64, 0, sourceLength)

	return ind, err
}

// NewDefaultADXRWithKnownSourceLength creates an Average Directional Index Rating (ADXR) for offline usage with default parameters
func NewDefaultADXRWithKnownSourceLength(sourceLength int) (indicator *ADXR, err error) {

	ind, err := NewDefaultADXR()
	ind.Data = make([]float64, 0, sourceLength)
	return ind, err
}

func NewADXRForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *ADXR, err error) {
	newADXR, err := NewADXR(timePeriod)
	priceStream.AddTickSubscription(newADXR)
	return newADXR, err
}

// NewADXRForStreamWithKnownSourceLength creates an Average Directional Index Rating (ADXR) for offline usage with a source data stream
func NewADXRForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *ADXR, err error) {
	ind, err := NewADXRWithKnownSourceLength(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultADXRForStreamWithKnownSourceLength creates an Average Directional Index Rating (ADXR) for offline usage with a source data stream
func NewDefaultADXRForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *ADXR, err error) {
	ind, err := NewDefaultADXRWithKnownSourceLength(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// GetTimePeriod returns the configured ADX timePeriod
func (ind *ADXRWithoutStorage) GetTimePeriod() int {
	return ind.timePeriod
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *ADXRWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	ind.adx.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
