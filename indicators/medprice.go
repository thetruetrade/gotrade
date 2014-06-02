package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type MedPriceWithoutStorage struct {
	*baseIndicator

	// private variables
	valueAvailableAction ValueAvailableAction
}

func NewMedPriceWithoutStorage(valueAvailableAction ValueAvailableAction) (indicator *MedPriceWithoutStorage, err error) {
	newVar := MedPriceWithoutStorage{baseIndicator: newBaseIndicator(0)}

	newVar.valueAvailableAction = valueAvailableAction

	return &newVar, nil
}

type MedPrice struct {
	*MedPriceWithoutStorage

	// public variables
	Data []float64
}

func NewMedPrice() (indicator *MedPrice, err error) {
	newVar := MedPrice{}
	newVar.MedPriceWithoutStorage, err = NewMedPriceWithoutStorage(func(dataItem float64, streamBarIndex int) {
		newVar.Data = append(newVar.Data, dataItem)
	})

	return &newVar, err
}

func NewMedPriceForStream(priceStream *gotrade.DOHLCVStream) (indicator *MedPrice, err error) {
	newVar, err := NewMedPrice()
	priceStream.AddTickSubscription(newVar)
	return newVar, err
}

func (ind *MedPriceWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.dataLength += 1

	if ind.validFromBar == -1 {
		ind.validFromBar = streamBarIndex
	}

	result := (tickData.H() + tickData.L()) / float64(2.0)

	if result > ind.maxValue {
		ind.maxValue = result
	}

	if result < ind.minValue {
		ind.minValue = result
	}

	ind.valueAvailableAction(result, streamBarIndex)
}
