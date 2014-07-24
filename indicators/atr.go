package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// An Average True Range (Atr), no storage, for use in other indicators
type AtrWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	trueRange            *TrueRangeWithoutStorage
	sma                  *SmaWithoutStorage
	previousAvgTrueRange float64
	multiplier           float64
	timePeriod           int
}

// NewAtrWithoutStorage creates an Average True Range (Atr) without storage
func NewAtrWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *AtrWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for an Atr indicator is 1
	if timePeriod < 1 {
		return nil, errors.New("timePeriod is less than the minimum (1)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := timePeriod
	ind := AtrWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		multiplier:           float64(timePeriod - 1),
		previousAvgTrueRange: -1,
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	ind.sma, err = NewSmaWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.previousAvgTrueRange = dataItem

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

		// update the maximum result value
		if dataItem > ind.maxValue {
			ind.maxValue = dataItem
		}

		// update the minimum result value
		if dataItem < ind.minValue {
			ind.minValue = dataItem
		}

		// notify of a new result value though the value available action
		ind.valueAvailableAction(dataItem, streamBarIndex)
	})

	ind.trueRange, err = NewTrueRangeWithoutStorage(func(dataItem float64, streamBarIndex int) {

		if ind.previousAvgTrueRange == -1 {
			ind.sma.ReceiveTick(dataItem, streamBarIndex)
		} else {

			// increment the number of results this indicator can be expected to return
			ind.dataLength += 1

			result := ((ind.previousAvgTrueRange * ind.multiplier) + dataItem) / float64(ind.timePeriod)

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

			// update the previous true range for the next tick
			ind.previousAvgTrueRange = result
		}

	})
	return &ind, nil
}

// An Average True Range Indicator
type Atr struct {
	*AtrWithoutStorage

	// public variables
	Data []float64
}

// NewAtr creates an Average True Range (Atr) for online usage
func NewAtr(timePeriod int) (indicator *Atr, err error) {
	ind := Atr{}
	ind.AtrWithoutStorage, err = NewAtrWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.Data = append(ind.Data, dataItem)
	})

	return &ind, err
}

// NewDefaultAtr creates an Average True Range (Atr) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultAtr() (indicator *Atr, err error) {
	timePeriod := 14
	return NewAtr(timePeriod)
}

// NewAtrWithKnownSourceLength creates an Average True Range (Atr) for offline usage
func NewAtrWithKnownSourceLength(sourceLength int, timePeriod int) (indicator *Atr, err error) {
	ind, err := NewAtr(timePeriod)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultAtrWithKnownSourceLength creates an Average True Range (Atr) for offline usage with default parameters
func NewDefaultAtrWithKnownSourceLength(sourceLength int) (indicator *Atr, err error) {
	ind, err := NewDefaultAtr()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewAtrForStream creates an Average True Range (Atr) for online usage with a source data stream
func NewAtrForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *Atr, err error) {
	ind, err := NewAtr(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAtrForStream creates an Average True Range (Atr) for online usage with a source data stream
func NewDefaultAtrForStream(priceStream *gotrade.DOHLCVStream) (indicator *Atr, err error) {
	ind, err := NewDefaultAtr()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewAtrForStreamWithKnownSourceLength creates an Average True Range (Atr) for offline usage with a source data stream
func NewAtrForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *Atr, err error) {
	ind, err := NewAtrWithKnownSourceLength(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultAtrForStreamWithKnownSourceLength creates an Average True Range (Atr) for offline usage with a source data stream
func NewDefaultAtrForStreamWithKnownSourceLength(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Atr, err error) {
	ind, err := NewDefaultAtrWithKnownSourceLength(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *AtrWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	// update the current true range
	ind.trueRange.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
