package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type TypPriceWithoutStorage struct {
	*baseIndicatorWithFloatBounds
}

func NewTypPriceWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *TypPriceWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	lookback := 0
	ind := TypPriceWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
	}

	return &ind, nil
}

// A Typical Price Indicator (TypPrice)
type TypPrice struct {
	*TypPriceWithoutStorage

	// public variables
	Data []float64
}

// NewTypPrice creates a Typical Price Indicator (TypPrice) for online usage
func NewTypPrice() (indicator *TypPrice, err error) {
	ind := TypPrice{}
	ind.TypPriceWithoutStorage, err = NewTypPriceWithoutStorage(func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewTypPriceWithSrcLen creates a Typical Price Indicator (TypPrice) for offline usage
func NewTypPriceWithSrcLen(sourceLength uint) (indicator *TypPrice, err error) {
	ind, err := NewTypPrice()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewTypPriceForStream creates a Typical Price Indicator (TypPrice) for online usage with a source data stream
func NewTypPriceForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *TypPrice, err error) {
	ind, err := NewTypPrice()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewTypPriceForStreamWithSrcLen creates a Typical Price Indicator (TypPrice) for offline usage with a source data stream
func NewTypPriceForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *TypPrice, err error) {
	ind, err := NewTypPriceWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *TypPriceWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	result := (tickData.H() + tickData.L() + tickData.C()) / float64(3.0)

	ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
}
