package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type OBVWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodCounter        int
	previousOBV          float64
	previousClose        float64
	valueAvailableAction ValueAvailableAction
}

func NewOBVWithoutStorage(valueAvailableAction ValueAvailableAction) (indicator *OBVWithoutStorage, err error) {
	newVar := OBVWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(0),
		periodCounter: -1,
		previousOBV:   0.0,
		previousClose: 0.0}

	newVar.valueAvailableAction = valueAvailableAction

	return &newVar, nil
}

type OBV struct {
	*OBVWithoutStorage

	// public variables
	Data []float64
}

func NewOBV() (indicator *OBV, err error) {
	newVar := OBV{}
	newVar.OBVWithoutStorage, err = NewOBVWithoutStorage(func(dataItem float64, streamBarIndex int) {
		newVar.Data = append(newVar.Data, dataItem)
	})

	return &newVar, err
}

func NewOBVForStream(priceStream *gotrade.DOHLCVStream) (indicator *OBV, err error) {
	newVar, err := NewOBV()
	priceStream.AddTickSubscription(newVar)
	return newVar, err
}

func (ind *OBVWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1

	if ind.periodCounter <= 0 {
		ind.previousOBV = tickData.V()
		ind.previousClose = tickData.C()

		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}

		result := ind.previousOBV

		if result > ind.maxValue {
			ind.maxValue = result
		}

		if result < ind.minValue {
			ind.minValue = result
		}

		ind.valueAvailableAction(result, streamBarIndex)
	}

	if ind.periodCounter > 0 {
		closePrice := tickData.C()
		if closePrice > ind.previousClose {
			ind.previousOBV += tickData.V()
		} else if closePrice < ind.previousClose {
			ind.previousOBV -= tickData.V()
		}

		ind.dataLength += 1

		result := ind.previousOBV

		if result > ind.maxValue {
			ind.maxValue = result
		}

		if result < ind.minValue {
			ind.minValue = result
		}

		ind.valueAvailableAction(result, streamBarIndex)
		ind.previousClose = tickData.C()
	}
}
