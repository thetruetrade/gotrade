package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Bollinger Band Indicator (BollingerBand), no storage, for use in other indicators
type BollingerBandsWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionBollinger
	sma                  *SmaWithoutStorage
	stdDev               *StdDevWithoutStorage
	currentSma           float64
	timePeriod           int
}

// NewBollingerBandsWithoutStorage creates a Bollinger Band Indicator (BollingerBand) without storage
func NewBollingerBandsWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionBollinger) (indicator *BollingerBandsWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for a Bollinger Band indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod is less than the minimum (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := timePeriod - 1
	ind := BollingerBandsWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		currentSma:           0.0,
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	ind.sma, err = NewSmaWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentSma = dataItem
	})

	ind.stdDev, err = NewStdDevWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}

		var upperBand = ind.currentSma + 2*dataItem
		var lowerBand = ind.currentSma - 2*dataItem

		// update the maximum result value
		if upperBand > ind.maxValue {
			ind.maxValue = upperBand
		}

		// update the minimum result value
		if lowerBand < ind.minValue {
			ind.minValue = lowerBand
		}

		// notify of a new result value though the value available action
		ind.valueAvailableAction(upperBand, ind.currentSma, lowerBand, streamBarIndex)
	})

	return &ind, nil
}

// A Bollinger Band Indicator
type BollingerBands struct {
	*BollingerBandsWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	UpperBand  []float64
	MiddleBand []float64
	LowerBand  []float64
}

// NewBollingerBands creates a Bollinger Band Indicator (BollingerBand) for online usage
func NewBollingerBands(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
	ind := BollingerBands{selectData: selectData}

	ind.BollingerBandsWithoutStorage, err = NewBollingerBandsWithoutStorage(
		timePeriod,
		func(dataItemUpperBand float64, dataItemMiddleBand float64, dataItemLowerBand float64, streamBarIndex int) {
			ind.UpperBand = append(ind.UpperBand, dataItemUpperBand)
			ind.MiddleBand = append(ind.MiddleBand, dataItemMiddleBand)
			ind.LowerBand = append(ind.LowerBand, dataItemLowerBand)
		})

	return &ind, err
}

// NewDefaultBollingerBands creates a Bollinger Band Indicator (BollingerBand) for online usage with default parameters
//	- timePeriod: 5
//  - selectData: useClosePrice
func NewDefaultBollingerBands() (indicator *BollingerBands, err error) {
	timePeriod := 5
	return NewBollingerBands(timePeriod, gotrade.UseClosePrice)
}

// NewBollingerBandsWithSrcLen creates a Bollinger Band Indicator (BollingerBand) for offline usage
func NewBollingerBandsWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
	ind, err := NewBollingerBands(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.UpperBand = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
		ind.MiddleBand = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
		ind.LowerBand = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultBollingerBandsWithSrcLen creates a Bollinger Band Indicator (BollingerBand) for offline usage
func NewDefaultBollingerBandsWithSrcLen(sourceLength uint) (indicator *BollingerBands, err error) {
	ind, err := NewDefaultBollingerBands()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.UpperBand = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
		ind.MiddleBand = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
		ind.LowerBand = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewBollingerBandsForStream creates a Bollinger Bands Indicator (BollingerBand) for online usage with a source data stream
func NewBollingerBandsForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
	ind, err := NewBollingerBands(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultBollingerBandsForStream creates a Bollinger Bands Indicator (BollingerBand) for online usage with a source data stream
func NewDefaultBollingerBandsForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *BollingerBands, err error) {
	ind, err := NewDefaultBollingerBands()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewBollingerBandsForStreamWithSrcLen creates a Bollinger Bands Indicator (BollingerBand) for online usage with a source data stream
func NewBollingerBandsForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
	ind, err := NewBollingerBandsWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultBollingerBandsForStreamWithSrcLen creates a Bollinger Bands Indicator (BollingerBand) for online usage with a source data stream
func NewDefaultBollingerBandsForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *BollingerBands, err error) {
	ind, err := NewDefaultBollingerBandsWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *BollingerBands) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData float64 = ind.selectData(tickData)
	ind.RecieveTick(selectedData, streamBarIndex)
}

// ReceiveTick consumes a source data float price tick
func (ind *BollingerBandsWithoutStorage) RecieveTick(tickData float64, streamBarIndex int) {
	ind.sma.ReceiveTick(tickData, streamBarIndex)
	ind.stdDev.ReceiveTick(tickData, streamBarIndex)
}
