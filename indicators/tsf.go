package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// A Time Series Forecast Indicator (Tsf)
type Tsf struct {
	*LinRegWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewTsf creates a Time Series Forecast Indicator (Tsf) for online usage
func NewTsf(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Tsf, err error) {
	ind := Tsf{selectData: selectData}
	ind.LinRegWithoutStorage, err = NewLinRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			result := intercept + slope*float64(timePeriod)

			// update the maximum result value
			if result > ind.LinRegWithoutStorage.maxValue {
				ind.LinRegWithoutStorage.maxValue = result
			}

			// update the minimum result value
			if result < ind.LinRegWithoutStorage.minValue {
				ind.LinRegWithoutStorage.minValue = result
			}

			ind.Data = append(ind.Data, result)
		})

	return &ind, err
}

// NewDefaultTsf creates a Time Series Forecast Indicator (Tsf) for online usage with default parameters
//	- timePeriod: 10
func NewDefaultTsf() (indicator *Tsf, err error) {
	timePeriod := 10
	return NewTsf(timePeriod, gotrade.UseClosePrice)
}

// NewTsfWithSrcLen creates a Time Series Forecast Indicator (Tsf) for offline usage
func NewTsfWithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Tsf, err error) {
	ind, err := NewTsf(timePeriod, selectData)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultTsfWithSrcLen creates a Time Series Forecast Indicator (Tsf) for offline usage with default parameters
func NewDefaultTsfWithSrcLen(sourceLength int) (indicator *Tsf, err error) {
	ind, err := NewDefaultTsf()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewTsfForStream creates a Time Series Forecast Indicator (Tsf) for online usage with a source data stream
func NewTsfForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Tsf, err error) {
	ind, err := NewTsf(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultTsfForStream creates a Time Series Forecast Indicator (Tsf) for online usage with a source data stream
func NewDefaultTsfForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Tsf, err error) {
	ind, err := NewDefaultTsf()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewTsfForStreamWithSrcLen creates a Time Series Forecast Indicator (Tsf) for offline usage with a source data stream
func NewTsfForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Tsf, err error) {
	ind, err := NewTsfWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultTsfForStreamWithSrcLen creates a Time Series Forecast Indicator (Tsf) for offline usage with a source data stream
func NewDefaultTsfForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *Tsf, err error) {
	ind, err := NewDefaultTsfWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *Tsf) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}
