package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type TypPriceWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
}

func NewTypPriceWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *TypPriceWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	lookback := 0
	ind := TypPriceWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		valueAvailableAction: valueAvailableAction,
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
func NewTypPriceWithSrcLen(sourceLength int) (indicator *TypPrice, err error) {
	ind, err := NewTypPrice()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewTypPriceForStream creates a Typical Price Indicator (TypPrice) for online usage with a source data stream
func NewTypPriceForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *TypPrice, err error) {
	ind, err := NewTypPrice()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewTypPriceForStreamWithSrcLen creates a Typical Price Indicator (TypPrice) for offline usage with a source data stream
func NewTypPriceForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *TypPrice, err error) {
	ind, err := NewTypPriceWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *TypPriceWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	// increment the number of results this indicator can be expected to return
	ind.dataLength += 1

	if ind.validFromBar == -1 {
		// set the streamBarIndex from which this indicator returns valid results
		ind.validFromBar = streamBarIndex
	}

	result := (tickData.H() + tickData.L() + tickData.C()) / float64(3.0)

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
