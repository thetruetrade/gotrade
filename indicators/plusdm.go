package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Plus Directional Movement Indicator (PlusDm), no storage, for use in other indicators
type PlusDmWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodCounter  int
	previousHigh   float64
	previousLow    float64
	previousPlusDm float64
	timePeriod     int
}

// NewPlusDmWithoutStorage creates a Plus Directional Movement Indicator (PlusDm) without storage
func NewPlusDmWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *PlusDmWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for this indicator is 1
	if timePeriod < 1 {
		return nil, errors.New("timePeriod is less than the minimum (1)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := 1
	if timePeriod > 1 {
		lookback = timePeriod - 1
	}
	newPlusDm := PlusDmWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		periodCounter:                -1,
		previousPlusDm:               0.0,
		timePeriod:                   timePeriod,
	}

	return &newPlusDm, nil
}

// A Plus Directional Movement Indicator (PlusDm)
type PlusDm struct {
	*PlusDmWithoutStorage

	// public variables
	Data []float64
}

// NewPlusDm creates a Plus Directional Movement Indicator (PlusDm) for online usage
func NewPlusDm(timePeriod int) (indicator *PlusDm, err error) {
	ind := PlusDm{}
	ind.PlusDmWithoutStorage, err = NewPlusDmWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultPlusDm creates a Plus Directional Movement Indicator (PlusDm) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultPlusDm() (indicator *PlusDm, err error) {
	timePeriod := 14
	return NewPlusDm(timePeriod)
}

// NewPlusDmWithSrcLen creates a Plus Directional Movement Indicator (PlusDm) for offline usage
func NewPlusDmWithSrcLen(sourceLength uint, timePeriod int) (indicator *PlusDm, err error) {
	ind, err := NewPlusDm(timePeriod)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultPlusDmWithSrcLen creates a Plus Directional Movement Indicator (PlusDm) for offline usage with default parameters
func NewDefaultPlusDmWithSrcLen(sourceLength uint) (indicator *PlusDm, err error) {
	ind, err := NewDefaultPlusDm()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewPlusDmForStream creates a Plus Directional Movement Indicator (PlusDm) for online usage with a source data stream
func NewPlusDmForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *PlusDm, err error) {
	ind, err := NewPlusDm(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultPlusDmForStream creates a Plus Directional Movement Indicator (PlusDm) for online usage with a source data stream
func NewDefaultPlusDmForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *PlusDm, err error) {
	ind, err := NewDefaultPlusDm()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewPlusDmForStreamWithSrcLen creates a Plus Directional Movement Indicator (PlusDm) for offline usage with a source data stream
func NewPlusDmForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int) (indicator *PlusDm, err error) {
	ind, err := NewPlusDmWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultPlusDmForStreamWithSrcLen creates a Plus Directional Movement Indicator (PlusDm) for offline usage with a source data stream
func NewDefaultPlusDmForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *PlusDm, err error) {
	ind, err := NewDefaultPlusDmWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *PlusDmWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	high := tickData.H()
	low := tickData.L()
	diffP := high - ind.previousHigh
	diffM := ind.previousLow - low

	if ind.lookbackPeriod == 1 {
		if ind.periodCounter > 0 {

			var result float64
			if (diffP > 0) && (diffP > diffM) {
				result = diffP
			} else {
				result = 0
			}

			ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
		}
	} else {
		if ind.periodCounter > 0 {
			if ind.periodCounter < ind.timePeriod {
				if (diffP > 0) && (diffP > diffM) {
					ind.previousPlusDm += diffP
				}

				if ind.periodCounter == ind.timePeriod-1 {

					result := ind.previousPlusDm

					ind.UpdateIndicatorWithNewValue(result, streamBarIndex)

				}
			} else {
				var result float64
				if (diffP > 0) && (diffP > diffM) {
					result = ind.previousPlusDm - (ind.previousPlusDm / float64(ind.timePeriod)) + diffP
				} else {
					result = ind.previousPlusDm - (ind.previousPlusDm / float64(ind.timePeriod))
				}

				ind.UpdateIndicatorWithNewValue(result, streamBarIndex)

				ind.previousPlusDm = result
			}
		}
	}

	ind.previousHigh = high
	ind.previousLow = low
}
