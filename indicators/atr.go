// Average True Range (ATR)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// An Average True Range Indicator
type ATRWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableAction
	trueRange            *TrueRangeWithoutStorage
	sma                  *SMAWithoutStorage
	previousAvgTrueRange float64
	multiplier           float64
}

// NewATRWithoutStorage returns a new Average True Range (ATR) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The ATR results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewATRWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableAction) (indicator *ATRWithoutStorage, err error) {
	newATR := ATRWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		multiplier:                  float64(timePeriod - 1),
		previousAvgTrueRange:        -1}
	newATR.valueAvailableAction = valueAvailableAction
	newATR.sma, err = NewSMAWithoutStorage(timePeriod, nil, func(dataItem float64, streamBarIndex int) {
		newATR.previousAvgTrueRange = dataItem

		newATR.dataLength += 1
		if newATR.validFromBar == -1 {
			newATR.validFromBar = streamBarIndex
		}

		if dataItem > newATR.maxValue {
			newATR.maxValue = dataItem
		}

		if dataItem < newATR.minValue {
			newATR.minValue = dataItem
		}
		newATR.valueAvailableAction(dataItem, streamBarIndex)
	})

	newATR.trueRange, err = NewTrueRangeWithoutStorage(func(dataItem float64, streamBarIndex int) {

		if newATR.previousAvgTrueRange == -1 {
			newATR.sma.ReceiveTick(dataItem, streamBarIndex)
		} else {

			newATR.dataLength += 1

			avgTrueRange := ((newATR.previousAvgTrueRange * newATR.multiplier) + dataItem) / float64(newATR.GetTimePeriod())

			if avgTrueRange > newATR.maxValue {
				newATR.maxValue = avgTrueRange
			}

			if avgTrueRange < newATR.minValue {
				newATR.minValue = avgTrueRange
			}

			newATR.valueAvailableAction(avgTrueRange, streamBarIndex)

			// update the previous true range for the next tick
			newATR.previousAvgTrueRange = avgTrueRange
		}

	})
	return &newATR, nil
}

// An Average True Range Indicator
type ATR struct {
	*ATRWithoutStorage

	// public variables
	Data []float64
}

// NewATR returns a new Average True Range (ATR) configured with the
// specified timePeriod. The ATR results are stored in the Data field.
func NewATR(timePeriod int) (indicator *ATR, err error) {
	newATR := ATR{}
	newATR.ATRWithoutStorage, err = NewATRWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		newATR.Data = append(newATR.Data, dataItem)
	})

	return &newATR, err
}

func NewATRForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *ATR, err error) {
	newATR, err := NewATR(timePeriod)
	priceStream.AddTickSubscription(newATR)
	return newATR, err
}

func (ind *ATRWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	// update the current true range
	ind.trueRange.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
