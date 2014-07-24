package indicators

import (
	"container/list"
	"errors"
	"github.com/thetruetrade/gotrade"
)

// An Average Directional Index Rating (Adxr), no storage
type AdxrWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	periodCounter        int
	periodHistory        *list.List
	adx                  *AdxWithoutStorage
	valueAvailableAction ValueAvailableActionFloat
	timePeriod           int
}

// NewAdxrWithoutStorage creates an Average Directional Index Rating (Adxr) without storage
func NewAdxrWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *AdxrWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// // the minimum timeperiod for an Adxr indicator is 2
	// if timePeriod < 2 {
	// 	return nil, errors.New("timePeriod is less than the minimum (2)")
	// }

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	ind := AdxrWithoutStorage{
		timePeriod:           timePeriod,
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        0,
		periodHistory:        list.New(),
		valueAvailableAction: valueAvailableAction,
	}

	ind.adx, err = NewAdxWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
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

// A Directional Movement Indicator Rating (Adxr)
type Adxr struct {
	*AdxrWithoutStorage

	// public variables
	Data []float64
}

// NewAdxr creates an Average Directional Index Rating (Adxr) for online usage
func NewAdxr(timePeriod int) (indicator *Adxr, err error) {
	ind := Adxr{}
	ind.AdxrWithoutStorage, err = NewAdxrWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultAdxr creates an Average Directional Index Rating (Adxr) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultAdxr() (indicator *Adxr, err error) {
	timePeriod := 14
	return NewAdxr(timePeriod)
}

// NewAdxrWithKnownSourceLength creates an Average Directional Index Rating (Adxr) for offline usage
func NewAdxrWithKnownSourceLength(sourceLength int, timePeriod int) (indicator *Adxr, err error) {
	ind, err := NewAdxr(timePeriod)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultAdxrWithKnownSourceLength creates an Average Directional Index Rating (Adxr) for offline usage with default parameters
func NewDefaultAdxrWithKnownSourceLength(sourceLength int) (indicator *Adxr, err error) {

	ind, err := NewDefaultAdxr()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

func NewAdxrForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *Adxr, err error) {
	newAdxr, err := NewAdxr(timePeriod)
	priceStream.AddTickSubscription(newAdxr)
	return newAdxr, err
}

// NewAdxrForStreamWithKnownSourceLength creates an Average Directional Index Rating (Adxr) for offline usage with a source data stream
func NewAdxrForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *Adxr, err error) {
	ind, err := NewAdxrWithKnownSourceLength(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAdxrForStreamWithKnownSourceLength creates an Average Directional Index Rating (Adxr) for offline usage with a source data stream
func NewDefaultAdxrForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Adxr, err error) {
	ind, err := NewDefaultAdxrWithKnownSourceLength(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// GetTimePeriod returns the configured ADdxrtimePeriod
func (ind *AdxrWithoutStorage) GetTimePeriod() int {
	return ind.timePeriod
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *AdxrWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	ind.adx.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
