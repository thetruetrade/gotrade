package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// An Average Price (AvgPrice), no storage, for use in other indicators
type AvgPriceWithoutStorage struct {
	*baseIndicatorWithFloatBounds
}

// NewAvgPriceWithoutStorage creates an Average Price(AvgPrice) without storage
func NewAvgPriceWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *AvgPriceWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	lookback := 0
	ind := AvgPriceWithoutStorage{
		baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(lookback, valueAvailableAction),
	}

	return &ind, nil
}

// An Average Price Indicator
type AvgPrice struct {
	*AvgPriceWithoutStorage

	// public variables
	Data []float64
}

// NewAvgPrice creates an Average Price (AvgPrice) for online usage
func NewAvgPrice() (indicator *AvgPrice, err error) {
	ind := AvgPrice{}
	ind.AvgPriceWithoutStorage, err = NewAvgPriceWithoutStorage(func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewAvgPriceWithSrcLen creates an Avgerage Price Indicator(AvgPrice) for offline usage
func NewAvgPriceWithSrcLen(sourceLength uint) (indicator *AvgPrice, err error) {
	ind, err := NewAvgPrice()
	ind.Data = make([]float64, 0, sourceLength)

	return ind, err
}

// NewAvgPriceForStream creates an Avgerage Price Indicator(AvgPrice) for online usage with a source data stream
func NewAvgPriceForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *AvgPrice, err error) {
	ind, err := NewAvgPrice()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewAvgPriceForStreamWithSrcLen creates an Avgerage Price Indicator(AvgPrice) for offline usage with a source data stream
func NewAvgPriceForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *AvgPrice, err error) {
	ind, err := NewAvgPriceWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *AvgPriceWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	result := (tickData.O() + tickData.H() + tickData.L() + tickData.C()) / float64(4.0)

	ind.UpdateIndicatorWithNewValue(result, streamBarIndex)
}
