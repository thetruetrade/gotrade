// Average True Range (MinusDI)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// A plus DM Indicator
type MinusDIWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	previousHigh         float64
	previousLow          float64
	previousMinusDM      float64
	previousTrueRange    float64
	currentTrueRange     float64
	trueRange            *TrueRange
}

// NewMinusDIWithoutStorage returns a new Minus Directional Indicator (MinusDI) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The MinusDI results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewMinusDIWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *MinusDIWithoutStorage, err error) {
	var lookback int = 1
	if timePeriod > 1 {
		lookback = timePeriod
	}
	newMinusDI := MinusDIWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               -1,
		previousMinusDM:             0.0,
		previousTrueRange:           0.0,
		currentTrueRange:            0.0}
	newMinusDI.trueRange, err = NewTrueRange()

	newMinusDI.trueRange.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newMinusDI.currentTrueRange = dataItem
	}

	newMinusDI.valueAvailableAction = valueAvailableAction

	return &newMinusDI, nil
}

// An Average True Range Indicator
type MinusDI struct {
	*MinusDIWithoutStorage

	// public variables
	Data []float64
}

// NewMinusDI returns a new Average True Range (MinusDI) configured with the
// specified timePeriod. The MinusDI results are stored in the Data field.
func NewMinusDI(timePeriod int) (indicator *MinusDI, err error) {
	newMinusDI := MinusDI{}
	newMinusDI.MinusDIWithoutStorage, err = NewMinusDIWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		newMinusDI.Data = append(newMinusDI.Data, dataItem)
	})

	return &newMinusDI, err
}

func NewMinusDIForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *MinusDI, err error) {
	newMinusDI, err := NewMinusDI(timePeriod)
	priceStream.AddTickSubscription(newMinusDI)
	return newMinusDI, err
}

func (ind *MinusDIWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	// forward to the true range indicator first using previous data
	ind.trueRange.ReceiveDOHLCVTick(tickData, streamBarIndex)

	ind.periodCounter += 1
	high := tickData.H()
	low := tickData.L()
	diffP := high - ind.previousHigh
	diffM := ind.previousLow - low

	if ind.lookbackPeriod == 1 {
		if ind.periodCounter > 0 {

			// forward to the true range indicator first using previous data
			ind.trueRange.ReceiveDOHLCVTick(tickData, streamBarIndex)

			var result float64
			if (diffM > 0) && (diffP < diffM) && ind.currentTrueRange != 0.0 {
				result = diffM / ind.currentTrueRange
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
				ind.previousTrueRange += ind.currentTrueRange
			} else {
				var result float64
				ind.previousTrueRange = ind.previousTrueRange - (ind.previousTrueRange / float64(ind.GetTimePeriod())) + ind.currentTrueRange
				if (diffM > 0) && (diffP < diffM) {
					ind.previousMinusDM = ind.previousMinusDM - (ind.previousMinusDM / float64(ind.GetTimePeriod())) + diffM
				} else {
					ind.previousMinusDM = ind.previousMinusDM - (ind.previousMinusDM / float64(ind.GetTimePeriod()))
				}

				if ind.previousTrueRange != 0.0 {
					result = float64(100.0) * ind.previousMinusDM / ind.previousTrueRange
				} else {
					result = 0.0
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
		}
	}

	ind.previousHigh = high
	ind.previousLow = low
}
