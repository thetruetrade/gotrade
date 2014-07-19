// Average True Range (MinusDM)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// A minus DM Indicator
type MinusDMWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	previousHigh         float64
	previousLow          float64
	previousMinusDM      float64
}

// NewMinusDMWithoutStorage returns a new Minus Directional Movement (MinusDM) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The MinusDM results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewMinusDMWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *MinusDMWithoutStorage, err error) {
	var lookback int = 1
	if timePeriod > 1 {
		lookback = timePeriod - 1
	}
	newMinusDM := MinusDMWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               -1,
		previousMinusDM:             0.0}
	newMinusDM.valueAvailableAction = valueAvailableAction

	return &newMinusDM, nil
}

// An Average True Range Indicator
type MinusDM struct {
	*MinusDMWithoutStorage

	// public variables
	Data []float64
}

// NewMinusDM returns a new Average True Range (MinusDM) configured with the
// specified timePeriod. The MinusDM results are stored in the Data field.
func NewMinusDM(timePeriod int) (indicator *MinusDM, err error) {
	newMinusDM := MinusDM{}
	newMinusDM.MinusDMWithoutStorage, err = NewMinusDMWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		newMinusDM.Data = append(newMinusDM.Data, dataItem)
	})

	return &newMinusDM, err
}

func NewMinusDMForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *MinusDM, err error) {
	newMinusDM, err := NewMinusDM(timePeriod)
	priceStream.AddTickSubscription(newMinusDM)
	return newMinusDM, err
}

func (ind *MinusDMWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	high := tickData.H()
	low := tickData.L()
	diffP := high - ind.previousHigh
	diffM := ind.previousLow - low

	if ind.lookbackPeriod == 1 {
		if ind.periodCounter > 0 {

			var result float64
			if (diffM > 0) && (diffP < diffM) {
				result = diffM
			} else {
				result = 0
			}

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
	} else {
		if ind.periodCounter > 0 {
			if ind.periodCounter < ind.GetTimePeriod() {
				if (diffM > 0) && (diffP < diffM) {
					ind.previousMinusDM += diffM
				}

				if ind.periodCounter == ind.GetTimePeriod()-1 {

					result := ind.previousMinusDM
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
			} else {
				var result float64
				if (diffM > 0) && (diffP < diffM) {
					result = ind.previousMinusDM - (ind.previousMinusDM / float64(ind.GetTimePeriod())) + diffM
				} else {
					result = ind.previousMinusDM - (ind.previousMinusDM / float64(ind.GetTimePeriod()))
				}

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

				ind.previousMinusDM = result
			}
		}
	}

	ind.previousHigh = high
	ind.previousLow = low
}
