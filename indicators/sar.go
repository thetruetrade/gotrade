package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Stop and Reverse Indicator (Sar), no storage, for use in other indicators
type SarWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction  ValueAvailableActionFloat
	periodCounter         int
	isLong                bool
	extremePoint          float64
	accelerationFactor    float64
	acceleration          float64
	accelerationFactorMax float64
	previousSar           float64
	previousHigh          float64
	previousLow           float64
	minusDM               *MinusDmWithoutStorage
	hasInitialDirection   bool
}

func NewSarWithoutStorage(accelerationFactor float64, accelerationFactorMax float64, valueAvailableAction ValueAvailableActionFloat) (indicator *SarWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum accelerationFactor for this indicator is 0
	if accelerationFactor < 0 {
		return nil, errors.New("accelerationFactor is less than the minimum (0)")
	}

	// check the maximum accelerationFactor
	if accelerationFactor > math.MaxFloat64 {
		return nil, errors.New("accelerationFactor is greater than the maximum float64 size")
	}

	// the minimum accelerationFactorMax for this indicator is 0
	if accelerationFactorMax < 0 {
		return nil, errors.New("accelerationFactorMax is less than the minimum (0)")
	}

	// check the maximum accelerationFactorMax
	if accelerationFactorMax > math.MaxFloat64 {
		return nil, errors.New("accelerationFactorMax is greater than the maximum float64 size")
	}

	lookback := 1
	ind := SarWithoutStorage{
		baseIndicator:         newBaseIndicator(lookback),
		baseFloatBounds:       newBaseFloatBounds(),
		periodCounter:         -2,
		isLong:                false,
		hasInitialDirection:   false,
		accelerationFactor:    accelerationFactor,
		accelerationFactorMax: accelerationFactorMax,
		extremePoint:          0.0,
		previousSar:           0.0,
		previousHigh:          0.0,
		previousLow:           0.0,
		acceleration:          accelerationFactor,
		valueAvailableAction:  valueAvailableAction,
	}

	ind.minusDM, err = NewMinusDmWithoutStorage(1, func(dataItem float64, streamBarIndex int) {
		if dataItem > 0 {
			ind.isLong = false
		} else {
			ind.isLong = true
		}
		ind.hasInitialDirection = true
	})

	return &ind, err
}

// A Stop and Reverse Indicator (Sar)
type Sar struct {
	*SarWithoutStorage

	// public variables
	Data []float64
}

// NewSar creates a Stop and Reverse Indicator (Sar) for online usage
func NewSar(accelerationFactor float64, accelerationFactorMax float64) (indicator *Sar, err error) {
	ind := Sar{}
	ind.SarWithoutStorage, err = NewSarWithoutStorage(accelerationFactor, accelerationFactorMax, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultSar creates a Stop and Reverse Indicator (Sar) for online usage with default parameters
//	- accelerationFactor: 0.02
//  - accelerationFactorMax: 0.2
func NewDefaultSar() (indicator *Sar, err error) {
	accelerationFactor := 0.02
	accelerationFactorMax := 0.2
	return NewSar(accelerationFactor, accelerationFactorMax)
}

// NewSarWithSrcLen creates a Stop and Reverse Indicator (Sar) for offline usage
func NewSarWithSrcLen(sourceLength int, accelerationFactor float64, accelerationFactorMax float64) (indicator *Sar, err error) {
	ind, err := NewSar(accelerationFactor, accelerationFactorMax)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultSarWithSrcLen creates a Stop and Reverse Indicator (Sar) for offline usage with default parameters
func NewDefaultSarWithSrcLen(sourceLength int) (indicator *Sar, err error) {
	ind, err := NewDefaultSar()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewSarForStream creates a Stop and Reverse Indicator (Sar) for online usage with a source data stream
func NewSarForStream(priceStream *gotrade.DOHLCVStream, accelerationFactor float64, accelerationFactorMax float64) (indicator *Sar, err error) {
	ind, err := NewSar(accelerationFactor, accelerationFactorMax)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultSarForStream creates a Stop and Reverse Indicator (Sar) for online usage with a source data stream
func NewDefaultSarForStream(priceStream *gotrade.DOHLCVStream) (indicator *Sar, err error) {
	ind, err := NewDefaultSar()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewSarForStreamWithSrcLen creates a Stop and Reverse Indicator (Sar) for offline usage with a source data stream
func NewSarForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream, accelerationFactor float64, accelerationFactorMax float64) (indicator *Sar, err error) {
	ind, err := NewSarWithSrcLen(sourceLength, accelerationFactor, accelerationFactorMax)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultSarForStreamWithSrcLen creates a Stop and Reverse Indicator (Sar) for offline usage with a source data stream
func NewDefaultSarForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Sar, err error) {
	ind, err := NewDefaultSarWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *SarWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	if ind.hasInitialDirection == false {
		ind.minusDM.ReceiveDOHLCVTick(tickData, streamBarIndex)
	}

	if ind.hasInitialDirection == true {
		if ind.periodCounter == 0 {
			if ind.isLong {
				ind.extremePoint = tickData.H()
				ind.previousSar = ind.previousLow
			} else {
				ind.extremePoint = tickData.L()
				ind.previousSar = ind.previousHigh
			}

			// this is a trick for the first iteration only,
			// the high low of the first bar will be used as the sar for the
			// second bar. According tyo TALib this is the closest to Wilders
			// originla idea of having the first entry day use the previous
			// extreme, except now that extreme is solely derived from the first
			// bar, supposedly Meta stock uses the same method.
			ind.previousHigh = tickData.H()
			ind.previousLow = tickData.L()
		}

		if ind.periodCounter >= 0 {
			var result float64 = 0.0
			if ind.isLong {
				if tickData.L() <= ind.previousSar {
					// switch to short if the low penetrates the Sar value
					ind.isLong = false
					ind.previousSar = ind.extremePoint

					// make sure the overridden Sar is within yesterdays and todays range
					if ind.previousSar < ind.previousHigh {
						ind.previousSar = ind.previousHigh
					}
					if ind.previousSar < tickData.H() {
						ind.previousSar = tickData.H()
					}

					result = ind.previousSar

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

					// adjust af and extremePoint
					ind.acceleration = ind.accelerationFactor
					ind.extremePoint = tickData.L()

					// calculate the new Sar
					var diff float64 = ind.extremePoint - ind.previousSar
					ind.previousSar = ind.previousSar + ind.acceleration*(diff)

					// make sure the overridden Sar is within yesterdays and todays range
					if ind.previousSar < ind.previousHigh {
						ind.previousSar = ind.previousHigh
					}
					if ind.previousSar < tickData.H() {
						ind.previousSar = tickData.H()
					}

				} else {
					// no switch

					// just output the current Sar
					result = ind.previousSar

					// increment the number of results this indicator can be expected to return
					ind.dataLength += 1

					if ind.validFromBar == -1 {
						// set the streamBarIndex from which this indicator returns valid results
						ind.validFromBar = streamBarIndex
					}

					if result > ind.maxValue {
						// update the maximum result value
						ind.maxValue = result
					}

					// update the minimum result value
					if result < ind.minValue {
						ind.minValue = result
					}

					// notify of a new result value though the value available action
					ind.valueAvailableAction(result, streamBarIndex)

					if tickData.H() > ind.extremePoint {
						// adjust af and extremePoint
						ind.extremePoint = tickData.H()
						ind.acceleration += ind.accelerationFactor
						if ind.acceleration > ind.accelerationFactorMax {
							ind.acceleration = ind.accelerationFactorMax
						}
					}

					// calculate the new Sar
					var diff float64 = ind.extremePoint - ind.previousSar
					ind.previousSar = ind.previousSar + ind.acceleration*(diff)

					// make sure the overridden Sar is within yesterdays and todays range
					if ind.previousSar > ind.previousLow {
						ind.previousSar = ind.previousLow
					}
					if ind.previousSar > tickData.L() {
						ind.previousSar = tickData.L()
					}
				}
			} else {
				// short
				// switch to long if the high penetrates the Sar value
				if tickData.H() >= ind.previousSar {
					ind.isLong = true
					ind.previousSar = ind.extremePoint

					// make sure the overridden Sar is within yesterdays and todays range
					if ind.previousSar > ind.previousLow {
						ind.previousSar = ind.previousLow
					}
					if ind.previousSar > tickData.L() {
						ind.previousSar = tickData.L()
					}

					result = ind.previousSar

					// increment the number of results this indicator can be expected to return
					ind.dataLength += 1

					if ind.validFromBar == -1 {
						// set the streamBarIndex from which this indicator returns valid results
						ind.validFromBar = streamBarIndex
					}

					if result > ind.maxValue {
						// update the maximum result value
						ind.maxValue = result
					}

					// update the minimum result value
					if result < ind.minValue {
						ind.minValue = result
					}

					// notify of a new result value though the value available action
					ind.valueAvailableAction(result, streamBarIndex)

					// adjust af and extremePoint
					ind.acceleration = ind.accelerationFactor
					ind.extremePoint = tickData.H()

					// calculate the new Sar
					var diff float64 = ind.extremePoint - ind.previousSar
					ind.previousSar = ind.previousSar + ind.acceleration*(diff)

					// make sure the overridden Sar is within yesterdays and todays range
					if ind.previousSar > ind.previousLow {
						ind.previousSar = ind.previousLow
					}
					if ind.previousSar > tickData.L() {
						ind.previousSar = tickData.L()
					}
				} else {
					// no switch

					// just output the current Sar
					result = ind.previousSar

					// increment the number of results this indicator can be expected to return
					ind.dataLength += 1

					if ind.validFromBar == -1 {
						// set the streamBarIndex from which this indicator returns valid results
						ind.validFromBar = streamBarIndex
					}

					if result > ind.maxValue {
						// update the maximum result value
						ind.maxValue = result
					}

					// update the minimum result value
					if result < ind.minValue {
						ind.minValue = result
					}

					// notify of a new result value though the value available action
					ind.valueAvailableAction(result, streamBarIndex)

					if tickData.L() < ind.extremePoint {
						// adjust af and extremePoint
						ind.extremePoint = tickData.L()
						ind.acceleration += ind.accelerationFactor
						if ind.acceleration > ind.accelerationFactorMax {
							ind.acceleration = ind.accelerationFactorMax
						}
					}

					// calculate the new Sar
					var diff float64 = ind.extremePoint - ind.previousSar
					ind.previousSar = ind.previousSar + ind.acceleration*(diff)

					// make sure the overridden Sar is within yesterdays and todays range
					if ind.previousSar < ind.previousHigh {
						ind.previousSar = ind.previousHigh
					}
					if ind.previousSar < tickData.H() {
						ind.previousSar = tickData.H()
					}
				}
			}
		}
	}

	ind.previousHigh = tickData.H()
	ind.previousLow = tickData.L()
}
