package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// An On Balance Volume Indicator (Obv), no storage, for use in other indicators
type ObvWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	periodCounter int
	previousObv   float64
	previousClose float64
}

// NewObvWithoutStorage creates an On Balance Volume Indicator (Obv) without storage
func NewObvWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *ObvWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	lookback := 0
	ind := ObvWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
		periodCounter:                -1,
		previousObv:                  0.0,
		previousClose:                0.0,
	}

	return &ind, nil
}

// A On Balance Volume Indicator (Obv)
type Obv struct {
	*ObvWithoutStorage

	// public variables
	Data []float64
}

// NewObv creates an On Balance Volume Indicator (Obv) for online usage
func NewObv() (indicator *Obv, err error) {
	ind := Obv{}
	ind.ObvWithoutStorage, err = NewObvWithoutStorage(func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewObvWithSrcLen creates an On Balance Volume (Obv) for offline usage
func NewObvWithSrcLen(sourceLength uint) (indicator *Obv, err error) {
	ind, err := NewObv()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewObvForStream creates an On Balance Volume (Obv) for online usage with a source data stream
func NewObvForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Obv, err error) {
	ind, err := NewObv()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewObvForStreamWithSrcLen creates an On Balance Volume (Obv) for offline usage with a source data stream
func NewObvForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Obv, err error) {
	ind, err := NewObvWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *ObvWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1

	if ind.periodCounter <= 0 {
		ind.previousObv = tickData.V()
		ind.previousClose = tickData.C()

		result := ind.previousObv

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
	}

	if ind.periodCounter > 0 {
		closePrice := tickData.C()
		if closePrice > ind.previousClose {
			ind.previousObv += tickData.V()
		} else if closePrice < ind.previousClose {
			ind.previousObv -= tickData.V()
		}

		result := ind.previousObv

		ind.UpdateIndicatorWithNewValue(result, streamBarIndex)

		ind.previousClose = tickData.C()
	}
}
