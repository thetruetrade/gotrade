// Average True Range (SAR)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// A plus DM Indicator
type SARWithoutStorage struct {
	*baseIndicator

	// private variables
	valueAvailableAction  ValueAvailableAction
	periodCounter         int
	isLong                bool
	extremePoint          float64
	accelerationFactor    float64
	acceleration          float64
	accelerationFactorMax float64
	previousSAR           float64
	previousHigh          float64
	previousLow           float64
	minusDM               *MinusDMWithoutStorage
	hasInitialDirection   bool
}

// NewSARWithoutStorage returns a new Parabolic Stop and Reverse Indicator (SAR) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The SAR results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewSARWithoutStorage(accelerationFactor float64, accelerationFactorMax float64, valueAvailableAction ValueAvailableAction) (indicator *SARWithoutStorage, err error) {
	newSAR := SARWithoutStorage{baseIndicator: newBaseIndicator(1),
		periodCounter:         -2,
		isLong:                false,
		hasInitialDirection:   false,
		accelerationFactor:    accelerationFactor,
		accelerationFactorMax: accelerationFactorMax,
		extremePoint:          0.0,
		previousSAR:           0.0,
		previousHigh:          0.0,
		previousLow:           0.0,
		acceleration:          accelerationFactor}
	newSAR.minusDM, err = NewMinusDMWithoutStorage(1, func(dataItem float64, streamBarIndex int) {
		if dataItem > 0 {
			newSAR.isLong = false
		} else {
			newSAR.isLong = true
		}
		newSAR.hasInitialDirection = true
	})

	newSAR.valueAvailableAction = valueAvailableAction

	return &newSAR, nil
}

// An Average True Range Indicator
type SAR struct {
	*SARWithoutStorage

	// public variables
	Data []float64
}

// NewSAR returns a new Parabolic Stop and Reverse Indicator (SAR) configured with the
// specified timePeriod. The SAR results are stored in the Data field.
func NewSAR(accelerationFactor float64, accelerationFactorMax float64) (indicator *SAR, err error) {
	newSAR := SAR{}
	newSAR.SARWithoutStorage, err = NewSARWithoutStorage(accelerationFactor, accelerationFactorMax, func(dataItem float64, streamBarIndex int) {
		newSAR.Data = append(newSAR.Data, dataItem)
	})

	return &newSAR, err
}

func NewSARForStream(priceStream *gotrade.DOHLCVStream, accelerationFactor float64, accelerationFactorMax float64) (indicator *SAR, err error) {
	newSAR, err := NewSAR(accelerationFactor, accelerationFactorMax)
	priceStream.AddTickSubscription(newSAR)
	return newSAR, err
}

func (ind *SARWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	if ind.hasInitialDirection == false {
		ind.minusDM.ReceiveDOHLCVTick(tickData, streamBarIndex)
	}

	if ind.hasInitialDirection == true {
		if ind.periodCounter == 0 {
			if ind.isLong {
				ind.extremePoint = tickData.H()
				ind.previousSAR = ind.previousLow
			} else {
				ind.extremePoint = tickData.L()
				ind.previousSAR = ind.previousHigh
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
				if tickData.L() <= ind.previousSAR {
					// switch to short if the low penetrates the SAR value
					ind.isLong = false
					ind.previousSAR = ind.extremePoint

					// make sure the overridden SAR is within yesterdays and todays range
					if ind.previousSAR < ind.previousHigh {
						ind.previousSAR = ind.previousHigh
					}
					if ind.previousSAR < tickData.H() {
						ind.previousSAR = tickData.H()
					}

					result = ind.previousSAR

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

					// adjust af and extremePoint
					ind.acceleration = ind.accelerationFactor
					ind.extremePoint = tickData.L()

					// calculate the new SAR
					var diff float64 = ind.extremePoint - ind.previousSAR
					ind.previousSAR = ind.previousSAR + ind.acceleration*(diff)

					// make sure the overridden SAR is within yesterdays and todays range
					if ind.previousSAR < ind.previousHigh {
						ind.previousSAR = ind.previousHigh
					}
					if ind.previousSAR < tickData.H() {
						ind.previousSAR = tickData.H()
					}

				} else {
					// no switch

					// just output the current SAR
					result = ind.previousSAR
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

					if tickData.H() > ind.extremePoint {
						// adjust af and extremePoint
						ind.extremePoint = tickData.H()
						ind.acceleration += ind.accelerationFactor
						if ind.acceleration > ind.accelerationFactorMax {
							ind.acceleration = ind.accelerationFactorMax
						}
					}

					// calculate the new SAR
					var diff float64 = ind.extremePoint - ind.previousSAR
					ind.previousSAR = ind.previousSAR + ind.acceleration*(diff)

					// make sure the overridden SAR is within yesterdays and todays range
					if ind.previousSAR > ind.previousLow {
						ind.previousSAR = ind.previousLow
					}
					if ind.previousSAR > tickData.L() {
						ind.previousSAR = tickData.L()
					}
				}
			} else {
				// short
				// switch to long if the high penetrates the SAR value
				if tickData.H() >= ind.previousSAR {
					ind.isLong = true
					ind.previousSAR = ind.extremePoint

					// make sure the overridden SAR is within yesterdays and todays range
					if ind.previousSAR > ind.previousLow {
						ind.previousSAR = ind.previousLow
					}
					if ind.previousSAR > tickData.L() {
						ind.previousSAR = tickData.L()
					}

					result = ind.previousSAR
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

					// adjust af and extremePoint
					ind.acceleration = ind.accelerationFactor
					ind.extremePoint = tickData.H()

					// calculate the new SAR
					var diff float64 = ind.extremePoint - ind.previousSAR
					ind.previousSAR = ind.previousSAR + ind.acceleration*(diff)

					// make sure the overridden SAR is within yesterdays and todays range
					if ind.previousSAR > ind.previousLow {
						ind.previousSAR = ind.previousLow
					}
					if ind.previousSAR > tickData.L() {
						ind.previousSAR = tickData.L()
					}
				} else {
					// no switch

					// just output the current SAR
					result = ind.previousSAR
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

					if tickData.L() < ind.extremePoint {
						// adjust af and extremePoint
						ind.extremePoint = tickData.L()
						ind.acceleration += ind.accelerationFactor
						if ind.acceleration > ind.accelerationFactorMax {
							ind.acceleration = ind.accelerationFactorMax
						}
					}

					// calculate the new SAR
					var diff float64 = ind.extremePoint - ind.previousSAR
					ind.previousSAR = ind.previousSAR + ind.acceleration*(diff)

					// make sure the overridden SAR is within yesterdays and todays range
					if ind.previousSAR < ind.previousHigh {
						ind.previousSAR = ind.previousHigh
					}
					if ind.previousSAR < tickData.H() {
						ind.previousSAR = tickData.H()
					}
				}
			}
		}
	}

	ind.previousHigh = tickData.H()
	ind.previousLow = tickData.L()
}
