package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// An Average Directional Index (Adx), no storage, for use in other indicators
type AdxWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	dx                   *DxWithoutStorage
	currentDX            float64
	sumDX                float64
	previousAdx          float64
	timePeriod           int
}

// NewAdxWithoutStorage creates an Average Directional Index (Adx) without storage
func NewAdxWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *AdxWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for an Adx indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod is less than the minimum (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := (2 * timePeriod) - 1
	ind := AdxWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		timePeriod:           timePeriod,
		periodCounter:        timePeriod * -1,
		currentDX:            0.0,
		sumDX:                0.0,
		previousAdx:          0.0,
		valueAvailableAction: valueAvailableAction,
	}

	ind.dx, err = NewDxWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {

		ind.currentDX = dataItem

		ind.periodCounter += 1
		if ind.periodCounter < 0 {
			ind.sumDX += ind.currentDX
		} else if ind.periodCounter == 0 {
			// increment the number of results this indicator can be expected to return
			ind.dataLength += 1

			if ind.validFromBar == -1 {
				// set the streamBarIndex from which this indicator returns valid results
				ind.validFromBar = streamBarIndex
			}

			ind.sumDX += ind.currentDX
			result := ind.sumDX / float64(ind.GetTimePeriod())

			// update the maximum result value
			if result > ind.maxValue {
				ind.maxValue = result
			}

			// update the minimum result value
			if result < ind.minValue {
				ind.minValue = result
			}
			ind.valueAvailableAction(result, streamBarIndex)
			ind.previousAdx = result

		} else {

			ind.dataLength += 1

			result := (ind.previousAdx*float64(ind.GetTimePeriod()-1) + ind.currentDX) / float64(ind.GetTimePeriod())

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

			ind.previousAdx = result
		}

	})

	return &ind, err
}

// A Directional Movement Indicator (Adx)
type Adx struct {
	*AdxWithoutStorage

	// public variables
	Data []float64
}

// NewAdx creates an Average Directional Index (Adx) for online usage
func NewAdx(timePeriod int) (indicator *Adx, err error) {

	ind := Adx{}
	ind.AdxWithoutStorage, err = NewAdxWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultAdx creates an Average Directional Index (Adx) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultAdx() (indicator *Adx, err error) {
	timePeriod := 14
	return NewAdx(timePeriod)
}

// NewAdxWithSrcLen creates an Average Directional Index (Adx) for offline usage
func NewAdxWithSrcLen(sourceLength int, timePeriod int) (indicator *Adx, err error) {
	ind, err := NewAdx(timePeriod)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultAdxWithSrcLen creates an Average Directional Index (Adx) for offline usage with default parameters
func NewDefaultAdxWithSrcLen(sourceLength int) (indicator *Adx, err error) {

	ind, err := NewDefaultAdx()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewAdxForStream creates an Average Directional Index (Adx) for online usage with a source data stream
func NewAdxForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *Adx, err error) {
	ind, err := NewAdx(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAdxForStream creates an Average Directional Index (Adx) for online usage with a source data stream
func NewDefaultAdxForStream(priceStream *gotrade.DOHLCVStream) (indicator *Adx, err error) {
	ind, err := NewDefaultAdx()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewAdxForStreamWithSrcLen creates an Average Directional Index (Adx) for offline usage with a source data stream
func NewAdxForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *Adx, err error) {
	ind, err := NewAdxWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAdxForStreamWithSrcLen creates an Average Directional Index (Adx) for offline usage with a source data stream
func NewDefaultAdxForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Adx, err error) {
	ind, err := NewDefaultAdxWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// GetTimePeriod returns the configured Adx timePeriod
func (ind *AdxWithoutStorage) GetTimePeriod() int {
	return ind.timePeriod
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *AdxWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.dx.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
