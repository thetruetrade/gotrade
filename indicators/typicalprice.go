package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type TypicalPriceWithoutStorage struct {
	*baseIndicator

	// private variables
	valueAvailableAction ValueAvailableAction
}

func NewTypicalPriceWithoutStorage(valueAvailableAction ValueAvailableAction) (indicator *TypicalPriceWithoutStorage, err error) {
	newVar := TypicalPriceWithoutStorage{baseIndicator: newBaseIndicator(0)}

	newVar.valueAvailableAction = valueAvailableAction

	return &newVar, nil
}

type TypicalPrice struct {
	*TypicalPriceWithoutStorage

	// public variables
	Data []float64
}

func NewTypicalPrice() (indicator *TypicalPrice, err error) {
	newVar := TypicalPrice{}
	newVar.TypicalPriceWithoutStorage, err = NewTypicalPriceWithoutStorage(func(dataItem float64, streamBarIndex int) {
		newVar.Data = append(newVar.Data, dataItem)
	})

	return &newVar, err
}

func NewTypicalPriceForStream(priceStream *gotrade.DOHLCVStream) (indicator *TypicalPrice, err error) {
	newVar, err := NewTypicalPrice()
	priceStream.AddTickSubscription(newVar)
	return newVar, err
}

func (ind *TypicalPriceWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.dataLength += 1

	if ind.validFromBar == -1 {
		ind.validFromBar = streamBarIndex
	}

	result := (tickData.H() + tickData.L() + tickData.C()) / float64(3.0)

	if result > ind.maxValue {
		ind.maxValue = result
	}

	if result < ind.minValue {
		ind.minValue = result
	}

	ind.valueAvailableAction(result, streamBarIndex)
}
