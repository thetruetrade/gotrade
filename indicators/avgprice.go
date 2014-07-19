package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type AvgPriceWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
}

func NewAvgPriceWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *AvgPriceWithoutStorage, err error) {
	newVar := AvgPriceWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(0)}

	newVar.valueAvailableAction = valueAvailableAction

	return &newVar, nil
}

type AvgPrice struct {
	*AvgPriceWithoutStorage

	// public variables
	Data []float64
}

func NewAvgPrice() (indicator *AvgPrice, err error) {
	newVar := AvgPrice{}
	newVar.AvgPriceWithoutStorage, err = NewAvgPriceWithoutStorage(func(dataItem float64, streamBarIndex int) {
		newVar.Data = append(newVar.Data, dataItem)
	})

	return &newVar, err
}

func NewAvgPriceForStream(priceStream *gotrade.DOHLCVStream) (indicator *AvgPrice, err error) {
	newVar, err := NewAvgPrice()
	priceStream.AddTickSubscription(newVar)
	return newVar, err
}

func (ind *AvgPriceWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.dataLength += 1

	if ind.validFromBar == -1 {
		ind.validFromBar = streamBarIndex
	}

	result := (tickData.O() + tickData.H() + tickData.L() + tickData.C()) / float64(4.0)

	if result > ind.maxValue {
		ind.maxValue = result
	}

	if result < ind.minValue {
		ind.minValue = result
	}

	ind.valueAvailableAction(result, streamBarIndex)
}
