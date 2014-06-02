//Advance Decline Line (ADL)
package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// An Accumulation Distribution Line Indicator
type ADLWithoutStorage struct {
	*baseIndicator

	// private variables
	valueAvailableAction ValueAvailableAction
	previousADL          float64
}

// NewADLWithoutStorage returns a new Accumulation Distribution Line (ADL)
// This version is intended for use by other indicators.
// The ADL results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewADLWithoutStorage(valueAvailableAction ValueAvailableAction) (indicator *ADLWithoutStorage, err error) {
	newADL := ADLWithoutStorage{baseIndicator: newBaseIndicator(0),
		previousADL: float64(0.0)}
	newADL.valueAvailableAction = valueAvailableAction

	return &newADL, nil
}

// An Accumulation Distribution Line Indicator
type ADL struct {
	*ADLWithoutStorage

	// public variables
	Data []float64
}

// NewADL returns a new Accumulation Distribution Line (ADL)
// The ADL results are stored in the Data field.
func NewADL() (indicator *ADL, err error) {
	newADL := ADL{}
	newADL.ADLWithoutStorage, err = NewADLWithoutStorage(func(dataItem float64, streamBarIndex int) {
		newADL.Data = append(newADL.Data, dataItem)
	})

	return &newADL, err
}

func NewADLForStream(priceStream *gotrade.DOHLCVStream) (indicator *ADL, err error) {
	newADL, err := NewADL()
	priceStream.AddTickSubscription(newADL)
	return newADL, err
}

func (ind *ADLWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.dataLength += 1

	moneyFlowMultiplier := ((tickData.C() - tickData.L()) - (tickData.H() - tickData.C())) / (tickData.H() - tickData.L())
	moneyFlowVolume := moneyFlowMultiplier * tickData.V()
	ADL := ind.previousADL + moneyFlowVolume

	if ind.validFromBar == -1 {
		ind.validFromBar = streamBarIndex
	}

	if ADL > ind.maxValue {
		ind.maxValue = ADL
	}

	if ADL < ind.minValue {
		ind.minValue = ADL
	}

	ind.valueAvailableAction(ADL, streamBarIndex)

	ind.previousADL = ADL
}
