package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type RSIWithoutStorage struct {
	*baseIndicator
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableAction
	periodCounter        int
	previousClose        float64
	previousGain         float64
	previousLoss         float64
}

func NewRSIWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *RSIWithoutStorage, err error) {
	newRSI := RSIWithoutStorage{baseIndicator: newBaseIndicator(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod * -1) - 1,
		previousClose:               0.0,
		previousGain:                0.0,
		previousLoss:                0.0}

	newRSI.selectData = selectData
	newRSI.valueAvailableAction = valueAvailableAction

	return &newRSI, err
}

// A Relative Strength Indicator
type RSI struct {
	*RSIWithoutStorage

	// public variables
	Data []float64
}

// NewRSI returns a new Relative Strength Indicator(RSI) configured with the
// specified timePeriod. The RSI results are stored in the DATA field.
func NewRSI(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RSI, err error) {
	newRSI := RSI{}
	newRSI.RSIWithoutStorage, err = NewRSIWithoutStorage(timePeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newRSI.Data = append(newRSI.Data, dataItem)
		})

	newRSI.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newRSI.Data = append(newRSI.Data, dataItem)
	}
	return &newRSI, err
}

func NewRSIForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *RSI, err error) {
	newRSI, err := NewRSI(timePeriod, selectData)
	priceStream.AddTickSubscription(newRSI)
	return newRSI, err
}

func (ind *RSIWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *RSIWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1

	if ind.periodCounter > ind.GetTimePeriod()*-1 {

		if ind.periodCounter <= 0 {

			if tickData > ind.previousClose {
				ind.previousGain += (tickData - ind.previousClose)
			} else {
				ind.previousLoss -= (tickData - ind.previousClose)
			}
		}

		if ind.periodCounter == 0 {
			ind.previousGain /= float64(ind.GetTimePeriod())
			ind.previousLoss /= float64(ind.GetTimePeriod())

			//    RSI = 100 * (prevGain/(prevGain+prevLoss))
			ind.dataLength += 1
			if ind.validFromBar == -1 {
				ind.validFromBar = streamBarIndex
			}

			var result float64
			if ind.previousGain+ind.previousLoss == 0.0 {
				result = 0.0
			} else {
				result = 100.0 * (ind.previousGain / (ind.previousGain + ind.previousLoss))
			}

			if result > ind.maxValue {
				ind.maxValue = result
			}

			if result < ind.minValue {
				ind.minValue = result
			}

			ind.valueAvailableAction(result, streamBarIndex)
		}

		if ind.periodCounter > 0 {
			ind.previousGain *= float64(ind.GetTimePeriod() - 1)
			ind.previousLoss *= float64(ind.GetTimePeriod() - 1)

			if tickData > ind.previousClose {
				ind.previousGain += (tickData - ind.previousClose)
			} else {
				ind.previousLoss -= (tickData - ind.previousClose)
			}

			ind.previousGain /= float64(ind.GetTimePeriod())
			ind.previousLoss /= float64(ind.GetTimePeriod())

			//    RSI = 100 * (prevGain/(prevGain+prevLoss))
			ind.dataLength += 1

			var result float64
			if ind.previousGain+ind.previousLoss == 0.0 {
				result = 0.0
			} else {
				result = 100.0 * (ind.previousGain / (ind.previousGain + ind.previousLoss))
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
	ind.previousClose = tickData
}
