package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// An Accumulation Distribution Line Indicator (Adl), no storage, for use in other indicators
type AdlWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	previousAdl float64
}

// NewAdlWithoutStorage creates an Accumulation Distribution Line Indicator (Adl) without storage
func NewAdlWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *AdlWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}
	lookback := 0
	ind := AdlWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		previousAdl:                  float64(0.0),
	}

	return &ind, nil
}

// An Accumulation Distribution Line Indicator (Adl)
type Adl struct {
	*AdlWithoutStorage

	// public variables
	Data []float64
}

// NewAdl creates an Accumulation Distribution Line Indicator (Adl) for online usage
func NewAdl() (indicator *Adl, err error) {
	ind := Adl{}
	ind.AdlWithoutStorage, err = NewAdlWithoutStorage(func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewAdlWithSrcLen creates an Accumulation Distribution Line Indicator (Adl) for offline usage
func NewAdlWithSrcLen(sourceLength uint) (indicator *Adl, err error) {
	ind, err := NewAdl()
	ind.Data = make([]float64, 0, sourceLength)
	return ind, err
}

// NewAdlForStream creates an Accumulation Distribution Line Indicator (Adl) for online usage with a source data stream
func NewAdlForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Adl, err error) {
	ind, err := NewAdl()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewAdlForStreamWithSrcLen creates an Accumulation Distribution Line Indicator (Adl) for offline usage with a source data stream
func NewAdlForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Adl, err error) {
	ind, err := NewAdlWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *AdlWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	moneyFlowMultiplier := ((tickData.C() - tickData.L()) - (tickData.H() - tickData.C())) / (tickData.H() - tickData.L())
	moneyFlowVolume := moneyFlowMultiplier * tickData.V()
	result := ind.previousAdl + moneyFlowVolume

	ind.UpdateIndicatorWithNewValue(result, streamBarIndex)

	ind.previousAdl = result
}
