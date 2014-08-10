package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// An On Balance Volume Indicator (Obv), no storage, for use in other indicators
type ObvWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	periodCounter        int
	previousObv          float64
	previousClose        float64
	valueAvailableAction ValueAvailableActionFloat
}

// NewObvWithoutStorage creates an On Balance Volume Indicator (Obv) without storage
func NewObvWithoutStorage(valueAvailableAction ValueAvailableActionFloat) (indicator *ObvWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	lookback := 0
	ind := ObvWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		periodCounter:        -1,
		previousObv:          0.0,
		previousClose:        0.0,
		valueAvailableAction: valueAvailableAction,
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
func NewObvWithSrcLen(sourceLength int) (indicator *Obv, err error) {
	ind, err := NewObv()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewObvForStream creates an On Balance Volume (Obv) for online usage with a source data stream
func NewObvForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Obv, err error) {
	ind, err := NewObv()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewObvForStreamWithSrcLen creates an On Balance Volume (Obv) for offline usage with a source data stream
func NewObvForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Obv, err error) {
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

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

		result := ind.previousObv

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

	if ind.periodCounter > 0 {
		closePrice := tickData.C()
		if closePrice > ind.previousClose {
			ind.previousObv += tickData.V()
		} else if closePrice < ind.previousClose {
			ind.previousObv -= tickData.V()
		}

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1

		result := ind.previousObv

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
		ind.previousClose = tickData.C()
	}
}
