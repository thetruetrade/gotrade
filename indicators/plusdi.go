// Average True Range (PlusDI)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// A plus DM Indicator
type PlusDIWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	previousHigh         float64
	previousLow          float64
	previousPlusDM       float64
	previousTrueRange    float64
	currentTrueRange     float64
	trueRange            *TrueRange
}

// NewPlusDIWithoutStorage returns a new Plus Directional Movement (PlusDI) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The PlusDI results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewPlusDIWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *PlusDIWithoutStorage, err error) {
	var lookback int = 1
	if timePeriod > 1 {
		lookback = timePeriod
	}
	newPlusDI := PlusDIWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               -1,
		previousPlusDM:              0.0,
		previousTrueRange:           0.0,
		currentTrueRange:            0.0}
	newPlusDI.trueRange, err = NewTrueRange()

	newPlusDI.trueRange.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newPlusDI.currentTrueRange = dataItem
	}

	newPlusDI.valueAvailableAction = valueAvailableAction

	return &newPlusDI, nil
}

// An Average True Range Indicator
type PlusDI struct {
	*PlusDIWithoutStorage

	// public variables
	Data []float64
}

// NewPlusDI returns a new Average True Range (PlusDI) configured with the
// specified timePeriod. The PlusDI results are stored in the Data field.
func NewPlusDI(timePeriod int) (indicator *PlusDI, err error) {
	newPlusDI := PlusDI{}
	newPlusDI.PlusDIWithoutStorage, err = NewPlusDIWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		newPlusDI.Data = append(newPlusDI.Data, dataItem)
	})

	return &newPlusDI, err
}

func NewPlusDIForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *PlusDI, err error) {
	newPlusDI, err := NewPlusDI(timePeriod)
	priceStream.AddTickSubscription(newPlusDI)
	return newPlusDI, err
}

func (ind *PlusDIWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

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
			if (diffP > 0) && (diffP > diffM) && ind.currentTrueRange != 0.0 {
				result = diffP / ind.currentTrueRange
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
				if (diffP > 0) && (diffP > diffM) {
					ind.previousPlusDM += diffP
				}
				ind.previousTrueRange += ind.currentTrueRange
			} else {
				var result float64
				ind.previousTrueRange = ind.previousTrueRange - (ind.previousTrueRange / float64(ind.GetTimePeriod())) + ind.currentTrueRange
				if (diffP > 0) && (diffP > diffM) {
					ind.previousPlusDM = ind.previousPlusDM - (ind.previousPlusDM / float64(ind.GetTimePeriod())) + diffP
				} else {
					ind.previousPlusDM = ind.previousPlusDM - (ind.previousPlusDM / float64(ind.GetTimePeriod()))
				}

				if ind.previousTrueRange != 0.0 {
					result = float64(100.0) * ind.previousPlusDM / ind.previousTrueRange
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
