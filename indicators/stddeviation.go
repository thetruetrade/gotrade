package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Standard Deviation Indicator (StdDev), no storage, for use in other indicators
type StdDevWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	variance             *VarWithoutStorage
	timePeriod           int
}

func NewStdDevWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *StdDevWithoutStorage, err error) {

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

	ind := StdDevWithoutStorage{
		baseIndicator:        newBaseIndicator(lookback),
		baseFloatBounds:      newBaseFloatBounds(),
		valueAvailableAction: valueAvailableAction,
		timePeriod:           timePeriod,
	}

	ind.valueAvailableAction = valueAvailableAction

	ind.variance, err = NewVarWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {

		// increment the number of results this indicator can be expected to return
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			// set the streamBarIndex from which this indicator returns valid results
			ind.validFromBar = streamBarIndex
		}

		result := math.Sqrt(dataItem)

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

// A Standard Deviation Indicator (StdDev)
type StdDev struct {
	*StdDevWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewStdDev creates a Standard Deviation Indicator (StdDev) for online usage
func NewStdDev(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *StdDev, err error) {
	ind := StdDev{selectData: selectData}
	ind.StdDevWithoutStorage, err = NewStdDevWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			ind.Data = append(ind.Data, dataItem)
		})

	return &ind, err
}

// NewDefaultStdDev creates a Standard Deviation Indicator (StdDev) for online usage with default parameters
//	- timePeriod: 10
func NewDefaultStdDev() (indicator *StdDev, err error) {
	timePeriod := 10
	return NewStdDev(timePeriod, gotrade.UseClosePrice)
}

// NewStdDevWithSrcLen creates a Standard Deviation Indicator (StdDev) for offline usage
func NewStdDevWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *StdDev, err error) {
	ind, err := NewStdDev(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultStdDevWithSrcLen creates a Standard Deviation Indicator (StdDev) for offline usage with default parameters
func NewDefaultStdDevWithSrcLen(sourceLength uint) (indicator *StdDev, err error) {
	ind, err := NewDefaultStdDev()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewStdDevForStream creates a Standard Deviation Indicator (StdDev) for online usage with a source data stream
func NewStdDevForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *StdDev, err error) {
	ind, err := NewStdDev(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultStdDevForStream creates a Standard Deviation Indicator (StdDev) for online usage with a source data stream
func NewDefaultStdDevForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *StdDev, err error) {
	ind, err := NewDefaultStdDev()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewStdDevForStreamWithSrcLen creates a Standard Deviation Indicator (StdDev) for offline usage with a source data stream
func NewStdDevForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *StdDev, err error) {
	ind, err := NewStdDevWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultStdDevForStreamWithSrcLen creates a Standard Deviation Indicator (StdDev) for offline usage with a source data stream
func NewDefaultStdDevForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *StdDev, err error) {
	ind, err := NewDefaultStdDevWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (stdDev *StdDev) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = stdDev.selectData(tickData)
	stdDev.ReceiveTick(selectedData, streamBarIndex)
}

func (stdDev *StdDevWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	stdDev.variance.ReceiveTick(tickData, streamBarIndex)
}
