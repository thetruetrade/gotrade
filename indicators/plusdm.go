// Average True Range (PlusDM)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// A plus DM Indicator
type PlusDMWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	previousHigh         float64
	previousLow          float64
	previousPlusDM       float64
}

// NewPlusDMWithoutStorage returns a new Plus Directional Movement (PlusDM) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The PlusDM results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewPlusDMWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *PlusDMWithoutStorage, err error) {
	var lookback int = 1
	if timePeriod > 1 {
		lookback = timePeriod - 1
	}
	newPlusDM := PlusDMWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               -1,
		previousPlusDM:              0.0}
	newPlusDM.valueAvailableAction = valueAvailableAction

	return &newPlusDM, nil
}

// An Average True Range Indicator
type PlusDM struct {
	*PlusDMWithoutStorage

	// public variables
	Data []float64
}

// NewPlusDM returns a new Average True Range (PlusDM) configured with the
// specified timePeriod. The PlusDM results are stored in the Data field.
func NewPlusDM(timePeriod int) (indicator *PlusDM, err error) {
	newPlusDM := PlusDM{}
	newPlusDM.PlusDMWithoutStorage, err = NewPlusDMWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		newPlusDM.Data = append(newPlusDM.Data, dataItem)
	})

	return &newPlusDM, err
}

func NewPlusDMForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *PlusDM, err error) {
	newPlusDM, err := NewPlusDM(timePeriod)
	priceStream.AddTickSubscription(newPlusDM)
	return newPlusDM, err
}

func (ind *PlusDMWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
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

				if ind.periodCounter == ind.GetTimePeriod()-1 {

					result := ind.previousPlusDM
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
				if (diffP > 0) && (diffP > diffM) {
					result = ind.previousPlusDM - (ind.previousPlusDM / float64(ind.GetTimePeriod())) + diffP
				} else {
					result = ind.previousPlusDM - (ind.previousPlusDM / float64(ind.GetTimePeriod()))
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

				ind.previousPlusDM = result
			}
		}
	}

	ind.previousHigh = high
	ind.previousLow = low
}
