package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

// A Triangular Moving Average Indicator (Trima), no storage, for use in other indicators
type TrimaWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	sma1                 *SmaWithoutStorage
	sma2                 *SmaWithoutStorage
	currentSma           float64
	timePeriod           int
}

// NewTrimaWithoutStorage creates a Triangular Moving Average Indicator (Trima) without storage
func NewTrimaWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *TrimaWithoutStorage, err error) {

	// an indicator without storage MUST have a value available action
	if valueAvailableAction == nil {
		return nil, ErrValueAvailableActionIsNil
	}

	// the minimum timeperiod for this indicator is 2
	if timePeriod < 2 {
		return nil, errors.New("timePeriod is less than the minimum (2)")
	}

	// check the maximum timeperiod
	if timePeriod > MaximumLookbackPeriod {
		return nil, errors.New("timePeriod is greater than the maximum (100000)")
	}

	lookback := timePeriod - 1
	ind := TrimaWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	var sma1Period int
	var sma2Period int

	if timePeriod%2 == 0 {
		// even
		sma1Period = timePeriod / 2
		sma2Period = (timePeriod / 2) + 1
	} else {
		// odd
		sma1Period = (timePeriod + 1) / 2
		sma2Period = (timePeriod + 1) / 2
	}

	ind.sma1, err = NewSmaWithoutStorage(sma1Period, func(dataItem float64, streamBarIndex int) {
		ind.currentSma = dataItem
		ind.sma2.ReceiveTick(dataItem, streamBarIndex)
	})

	ind.sma2, _ = NewSmaWithoutStorage(sma2Period, func(dataItem float64, streamBarIndex int) {

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

		result := dataItem

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
	})

	return &ind, err
}

// A Triangular Moving Average Indicator (Trima)
type Trima struct {
	*TrimaWithoutStorage
	selectData gotrade.DataSelectionFunc
	// public variables
	Data []float64
}

// NewTrima creates a Triangular Moving Average Indicator (Trima) for online usage
func NewTrima(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Trima, err error) {
	ind := Trima{selectData: selectData}
	ind.TrimaWithoutStorage, err = NewTrimaWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})
	return &ind, err
}

// NewDefaultTrima creates a Triangular Moving Average Indicator (Trima) for online usage with default parameters
//	- timePeriod: 30
func NewDefaultTrima() (indicator *Trima, err error) {
	timePeriod := 30
	return NewTrima(timePeriod, gotrade.UseClosePrice)
}

// NewTrimaWithSrcLen creates a Triangular Moving Average Indicator (Trima) for offline usage
func NewTrimaWithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Trima, err error) {
	ind, err := NewTrima(timePeriod, selectData)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultTrimaWithSrcLen creates a Triangular Moving Average Indicator (Trima) for offline usage with default parameters
func NewDefaultTrimaWithSrcLen(sourceLength int) (indicator *Trima, err error) {
	ind, err := NewDefaultTrima()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewTrimaForStream creates a Triangular Moving Average Indicator (Trima) for online usage with a source data stream
func NewTrimaForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Trima, err error) {
	ind, err := NewTrima(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultTrimaForStream creates a Triangular Moving Average Indicator (Trima) for online usage with a source data stream
func NewDefaultTrimaForStream(priceStream *gotrade.DOHLCVStream) (indicator *Trima, err error) {
	ind, err := NewDefaultTrima()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewTrimaForStreamWithSrcLen creates a Triangular Moving Average Indicator (Trima) for offline usage with a source data stream
func NewTrimaForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Trima, err error) {
	ind, err := NewTrimaWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultTrimaForStreamWithSrcLen creates a Triangular Moving Average Indicator (Trima) for offline usage with a source data stream
func NewDefaultTrimaForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Trima, err error) {
	ind, err := NewDefaultTrimaWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (tema *Trima) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = tema.selectData(tickData)
	tema.ReceiveTick(selectedData, streamBarIndex)
}

func (tema *TrimaWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	tema.sma1.ReceiveTick(tickData, streamBarIndex)
}
