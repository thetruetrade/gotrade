package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type ADXRWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	periodCounter        int
	periodHistory        *list.List
	adx                  *ADXWithoutStorage
	valueAvailableAction ValueAvailableActionFloat
}

func NewADXRWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *ADXRWithoutStorage, err error) {
	newADXR := ADXRWithoutStorage{baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter: 0,
		periodHistory: list.New()}

	newADXR.adx, err = NewADXWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		newADXR.periodHistory.PushBack(dataItem)

		if newADXR.periodCounter > newADXR.GetLookbackPeriod() {
			adxN := newADXR.periodHistory.Front().Value.(float64)
			result := (dataItem + adxN) / 2.0

			newADXR.dataLength += 1
			if newADXR.validFromBar == -1 {
				newADXR.validFromBar = streamBarIndex
			}

			if result > newADXR.maxValue {
				newADXR.maxValue = result
			}

			if result < newADXR.minValue {
				newADXR.minValue = result
			}

			newADXR.valueAvailableAction(result, streamBarIndex)
		}

		if newADXR.periodHistory.Len() >= newADXR.adx.GetTimePeriod() {
			first := newADXR.periodHistory.Front()
			newADXR.periodHistory.Remove(first)
		}
	})

	var lookback int = 3
	if timePeriod > 1 {
		lookback = timePeriod - 1 + newADXR.adx.GetLookbackPeriod()
	}
	newADXR.baseIndicatorWithFloatBounds = newBaseIndicatorWithFloatBounds(lookback)

	newADXR.valueAvailableAction = valueAvailableAction

	return &newADXR, nil
}

type ADXR struct {
	*ADXRWithoutStorage

	// public variables
	Data []float64
}

func NewADXR(timePeriod int) (indicator *ADXR, err error) {
	newADXR := ADXR{}
	newADXR.ADXRWithoutStorage, err = NewADXRWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newADXR.Data = append(newADXR.Data, dataItem)
		})

	return &newADXR, err
}

func NewADXRForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *ADXR, err error) {
	newADXR, err := NewADXR(timePeriod)
	priceStream.AddTickSubscription(newADXR)
	return newADXR, err
}

func (ind *ADXRWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	ind.adx.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
