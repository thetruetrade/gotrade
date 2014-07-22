package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// An Accumulation Distribution Line Indicator (ADL), no storage, for use in other indicators
type ADLWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	previousADL          float64
}

// NewADLWithoutStorage creates an Accumulation Distribution Line Indicator (ADL) without storage
func NewADLWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *ADLWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	ind := ADLWithoutStorage{
		baseIndicator:        newBaseIndicator(0),
		baseFloatBounds:      newBaseFloatBounds(),
		previousADL:          float64(0.0),
		valueAvailableAction: valueAvailableAction,
	}

	return &ind, nil
}

// An Accumulation Distribution Line Indicator (ADL)
type ADL struct {
	*ADLWithoutStorage

	// public variables
	Data []float64
}

// NewADL creates an Accumulation Distribution Line Indicator (ADL) for online usage
func NewADL() (indicator *ADL, err error) {
	ind := ADL{}
	ind.ADLWithoutStorage, err = NewADLWithoutStorage(func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewADLWithKnownSourceLength creates an Accumulation Distribution Line Indicator (ADL) for offline usage
func NewADLWithKnownSourceLength(sourceLength int) (indicator *ADL, err error) {
	ind, err := NewADL()
	ind.Data = make([]float64, 0, sourceLength)

	return ind, err
}

// NewADLForStream creates an Accumulation Distribution Line Indicator (ADL) for online usage with a source data stream
func NewADLForStream(priceStream *gotrade.DOHLCVStream) (indicator *ADL, err error) {
	ind, err := NewADL()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewADLForStreamWithKnownSourceLength creates an Accumulation Distribution Line Indicator (ADL) for offline usage with a source data stream
func NewADLForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *ADL, err error) {
	ind, err := NewADLWithKnownSourceLength(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *ADLWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	// increment the number of results this indicator can be expected to return
	ind.dataLength += 1

	moneyFlowMultiplier := ((tickData.C() - tickData.L()) - (tickData.H() - tickData.C())) / (tickData.H() - tickData.L())
	moneyFlowVolume := moneyFlowMultiplier * tickData.V()
	result := ind.previousADL + moneyFlowVolume

	if ind.validFromBar == -1 {
		// set the streamBarIndex from which this indicator returns valid results
		ind.validFromBar = streamBarIndex
	}

	// update the maximum result value
	if result > ind.maxValue {
		ind.maxValue = result
	}

	// update the minimum result value
	if result < ind.minValue {
		ind.minValue = result
	}

	// notify of a new result value though the value available action
	ind.valueAvailableAction(result, streamBarIndex)

	ind.previousADL = result
}
