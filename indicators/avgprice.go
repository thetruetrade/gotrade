package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// An Average Price (AvgPrice), no storage, for use in other indicators
type AvgPriceWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
}

// NewAvgPriceWithoutStorage creates an Average Price(AvgPrice) without storage
func NewAvgPriceWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *AvgPriceWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	lookback := 0
	ind := AvgPriceWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		valueAvailableAction: valueAvailableAction,
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
func NewAvgPriceWithSrcLen(sourceLength int) (indicator *AvgPrice, err error) {
	ind, err := NewAvgPrice()
	ind.Data = make([]float64, 0, sourceLength)

	return ind, err
}

// NewAvgPriceForStream creates an Avgerage Price Indicator(AvgPrice) for online usage with a source data stream
func NewAvgPriceForStream(priceStream *gotrade.DOHLCVStream) (indicator *AvgPrice, err error) {
	ind, err := NewAvgPrice()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewAvgPriceForStreamWithSrcLen creates an Avgerage Price Indicator(AvgPrice) for offline usage with a source data stream
func NewAvgPriceForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *AvgPrice, err error) {
	ind, err := NewAvgPriceWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *AvgPriceWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	// increment the number of results this indicator can be expected to return
	ind.dataLength += 1

	if ind.validFromBar == -1 {
		// set the streamBarIndex from which this indicator returns valid results
		ind.validFromBar = streamBarIndex
	}

	result := (tickData.O() + tickData.H() + tickData.L() + tickData.C()) / float64(4.0)

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
}
