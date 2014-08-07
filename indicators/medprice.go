package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// A Median Price Indicator (MedPrice), no storage, for use in other indicators
type MedPriceWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
}

// NewMedPriceWithoutStorage creates a Median Price Indicator (MedPrice) without storage
func NewMedPriceWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *MedPriceWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	lookback := 0
	ind := MedPriceWithoutStorage{
		baseIndicator:   newBaseIndicator(lookback),
		baseFloatBounds: newBaseFloatBounds(),
	}

	ind.valueAvailableAction = valueAvailableAction

	return &ind, nil
}

// A Median Price Indicator (MedPrice)
type MedPrice struct {
	*MedPriceWithoutStorage

	// public variables
	Data []float64
}

func NewMedPrice() (indicator *MedPrice, err error) {
	ind := MedPrice{}
	ind.MedPriceWithoutStorage, err = NewMedPriceWithoutStorage(func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultMedPrice creates a Median Price Indicator (MedPrice) for online usage with default parameters
func NewDefaultMedPrice() (indicator *MedPrice, err error) {
	return NewMedPrice()
}

// NewMedPriceWithSrcLen creates a Median Price Indicator (MedPrice) for offline usage
func NewMedPriceWithSrcLen(sourceLength int) (indicator *MedPrice, err error) {
	ind, err := NewMedPrice()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultMedPriceWithSrcLen creates a Median Price Indicator (MedPrice) for offline usage with default parameters
func NewDefaultMedPriceWithSrcLen(sourceLength int) (indicator *MedPrice, err error) {
	ind, err := NewDefaultMedPrice()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewMedPriceForStream creates a Median Price Indicator (MedPrice) for online usage with a source data stream
func NewMedPriceForStream(priceStream *gotrade.DOHLCVStream) (indicator *MedPrice, err error) {
	ind, err := NewMedPrice()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMedPriceForStream creates a Median Price Indicator (MedPrice) for online usage with a source data stream
func NewDefaultMedPriceForStream(priceStream *gotrade.DOHLCVStream) (indicator *MedPrice, err error) {
	ind, err := NewDefaultMedPrice()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewMedPriceForStreamWithSrcLen creates a Median Price Indicator (MedPrice) for offline usage with a source data stream
func NewMedPriceForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *MedPrice, err error) {
	ind, err := NewMedPriceWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultMedPriceForStreamWithSrcLen creates a Median Price Indicator (MedPrice) for offline usage with a source data stream
func NewDefaultMedPriceForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *MedPrice, err error) {
	ind, err := NewDefaultMedPriceWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *MedPriceWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {

	// increment the number of results this indicator can be expected to return
	ind.dataLength += 1

	if ind.validFromBar == -1 {
		// set the streamBarIndex from which this indicator returns valid results
		ind.validFromBar = streamBarIndex
	}

	result := (tickData.H() + tickData.L()) / float64(2.0)

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
