package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// An Average Directional Index (Adx), no storage, for use in other indicators
type AdxWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodCounter int
	dx            *DxWithoutStorage
	currentDX     float64
	sumDX         float64
	previousAdx   float64
	timePeriod    int
}

// NewAdxWithoutStorage creates an Average Directional Index (Adx) without storage
func NewAdxWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *AdxWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for an Adx indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod " + ErrStrBelowMinimum + " (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod " + ErrStrAboveMaximum + " (100000)")
	}

	lookback := (2 * timePeriod) - 1
	ind := AdxWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		periodCounter:                timePeriod * -1,
		currentDX:                    0.0,
		sumDX:                        0.0,
		previousAdx:                  0.0,
		timePeriod:                   timePeriod,
	}

	ind.dx, err = NewDxWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {

		ind.currentDX = dataItem

		ind.periodCounter += 1
		if ind.periodCounter < 0 {
			ind.sumDX += ind.currentDX
		} else if ind.periodCounter == 0 {
			ind.sumDX += ind.currentDX

			result := ind.sumDX / float64(ind.timePeriod)

			ind.UpdateIndicatorWithNewValue(result, streamBarIndex)

			ind.previousAdx = result

		} else {

			result := (ind.previousAdx*float64(ind.timePeriod-1) + ind.currentDX) / float64(ind.timePeriod)

			ind.UpdateIndicatorWithNewValue(result, streamBarIndex)

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
func NewAdxWithSrcLen(sourceLength uint, timePeriod int) (indicator *Adx, err error) {
	ind, err := NewAdx(timePeriod)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultAdxWithSrcLen creates an Average Directional Index (Adx) for offline usage with default parameters
func NewDefaultAdxWithSrcLen(sourceLength uint) (indicator *Adx, err error) {
	ind, err := NewDefaultAdx()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewAdxForStream creates an Average Directional Index (Adx) for online usage with a source data stream
func NewAdxForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *Adx, err error) {
	ind, err := NewAdx(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAdxForStream creates an Average Directional Index (Adx) for online usage with a source data stream
func NewDefaultAdxForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Adx, err error) {
	ind, err := NewDefaultAdx()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewAdxForStreamWithSrcLen creates an Average Directional Index (Adx) for offline usage with a source data stream
func NewAdxForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *Adx, err error) {
	ind, err := NewAdxWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAdxForStreamWithSrcLen creates an Average Directional Index (Adx) for offline usage with a source data stream
func NewDefaultAdxForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Adx, err error) {
	ind, err := NewDefaultAdxWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *AdxWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.dx.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
