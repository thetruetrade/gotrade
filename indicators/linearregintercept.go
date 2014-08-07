package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// A Linear Regression Intercept Indicator (LinRegInt)
type LinRegInt struct {
	*LinRegWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewLinRegInt creates a Linear Regression Intercept Indicator (LinRegInt) for online usage
func NewLinRegInt(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegInt, err error) {
	ind := LinRegInt{selectData: selectData}
	ind.LinRegWithoutStorage, err = NewLinRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			result := intercept

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

// NewDefaultLinRegInt creates a Linear Regression Intercept Indicator (LinRegInt) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultLinRegInt() (indicator *LinRegInt, err error) {
	timePeriod := 14
	return NewLinRegInt(timePeriod, gotrade.UseClosePrice)
}

// NewLinRegIntWithSrcLen creates a Linear Regression Intercept Indicator (LinRegInt) for offline usage
func NewLinRegIntWithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegInt, err error) {
	ind, err := NewLinRegInt(timePeriod, selectData)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultLinRegIntWithSrcLen creates a Linear Regression Intercept Indicator (LinRegInt) for offline usage with default parameters
func NewDefaultLinRegIntWithSrcLen(sourceLength int) (indicator *LinRegInt, err error) {
	ind, err := NewDefaultLinRegInt()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewLinRegIntForStream creates a Linear Regression Intercept Indicator (LinRegInt) for online usage with a source data stream
func NewLinRegIntForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegInt, err error) {
	ind, err := NewLinRegInt(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLinRegIntForStream creates a Linear Regression Intercept Indicator (LinRegInt) for online usage with a source data stream
func NewDefaultLinRegIntForStream(priceStream *gotrade.DOHLCVStream) (indicator *LinRegInt, err error) {
	ind, err := NewDefaultLinRegInt()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewLinRegIntForStreamWithSrcLen creates a Linear Regression Intercept Indicator (LinRegInt) for offline usage with a source data stream
func NewLinRegIntForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegInt, err error) {
	ind, err := NewLinRegIntWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLinRegIntForStreamWithSrcLen creates a Linear Regression Intercept Indicator (LinRegInt) for offline usage with a source data stream
func NewDefaultLinRegIntForStreamWithSrcLen(sourceLength int, priceStream *gotrade.DOHLCVStream) (indicator *LinRegInt, err error) {
	ind, err := NewDefaultLinRegIntWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *LinRegInt) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}
