package indicators

// DX = ( (+DI)-(-DI) ) / ( (+DI) + (-DI) )

import (
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// An Directional Movement Index Indicator (Dx), no storage, for use in other indicators
type DxWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	minusDI              *MinusDi
	plusDI               *PlusDi
	currentPlusDi        float64
	currentMinusDi       float64
	timePeriod           int
}

// NewDxWithoutStorage creates a Directional Movement Index Indicator (Dx) without storage
func NewDxWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *DxWithoutStorage, err error) {

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

	lookback := 2
	if timePeriod > 1 {
		lookback = timePeriod
	}

	ind := DxWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		currentPlusDi:        0.0,
		currentMinusDi:       0.0,
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	ind.minusDI, err = NewMinusDi(timePeriod)

	ind.minusDI.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		ind.currentMinusDi = dataItem
	}

	ind.plusDI, err = NewPlusDi(timePeriod)

	ind.plusDI.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		ind.currentPlusDi = dataItem

		var result float64
		tmp := ind.currentMinusDi + ind.currentPlusDi
		if tmp != 0.0 {
			result = 100.0 * (math.Abs(ind.currentMinusDi-ind.currentPlusDi) / tmp)
		} else {
			result = 0.0
		}

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1

		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

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

	return &ind, err
}

// A Directional Movement Index Indicator (Dx)
type Dx struct {
	*DxWithoutStorage

	// public variables
	Data []float64
}

// NewDx creates a Directional Movement Index Indicator (Dx) for online usage
func NewDx(timePeriod int) (indicator *Dx, err error) {

	ind := Dx{}
	ind.DxWithoutStorage, err = NewDxWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultDx creates a Directional Movement Index (Dx) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultDx() (indicator *Dx, err error) {
	timePeriod := 14
	return NewDx(timePeriod)
}

// NewDxWithSrcLen creates a Directional Movement Index (Dx) for offline usage
func NewDxWithSrcLen(sourceLength int, timePeriod int) (indicator *Dx, err error) {
	ind, err := NewDx(timePeriod)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultDxWithSrcLen creates a Directional Movement Index (Dx) for offline usage with default parameters
func NewDefaultDxWithSrcLen(sourceLength int) (indicator *Dx, err error) {
	ind, err := NewDefaultDx()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewDxForStream creates a Directional Movement Index (Dx) for online usage with a source data stream
func NewDxForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *Dx, err error) {
	ind, err := NewDx(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultDxForStream creates a Directional Movement Index (Dx) for online usage with a source data stream
func NewDefaultDxForStream(priceStream *gotrade.DOHLCVStream) (indicator *Dx, err error) {
	ind, err := NewDefaultDx()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDxForStreamWithSrcLen creates a Directional Movement Index (Dx) for offline usage with a source data stream
func NewDxForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *Dx, err error) {
	ind, err := NewDxWithSrcLen(sourceLength, timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultDxForStreamWithSrcLen creates a Directional Movement Index (Dx) for offline usage with a source data stream
func NewDefaultDxForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *Dx, err error) {
	ind, err := NewDefaultDxWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *DxWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.minusDI.ReceiveDOHLCVTick(tickData, streamBarIndex)
	ind.plusDI.ReceiveDOHLCVTick(tickData, streamBarIndex)
}
