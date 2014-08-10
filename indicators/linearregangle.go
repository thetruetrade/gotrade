package indicators

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Linear Regression Angle Indicator (LinRegAng)
type LinRegAng struct {
	*LinRegWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewLinRegAng creates a Linear Regression Angle Indicator (LinRegAng) for online usage
func NewLinRegAng(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegAng, err error) {
	ind := LinRegAng{selectData: selectData}
	ind.LinRegWithoutStorage, err = NewLinRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			result := math.Atan(slope) * (180.0 / math.Pi)

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

// NewDefaultLinRegAng creates a Linear Regression Angle Indicator (LinRegAng) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultLinRegAng() (indicator *LinRegAng, err error) {
	timePeriod := 14
	return NewLinRegAng(timePeriod, gotrade.UseClosePrice)
}

// NewLinRegAngWithSrcLen creates a Linear Regression Angle Indicator (LinRegAng) for offline usage
func NewLinRegAngWithSrcLen(sourceLength int, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegAng, err error) {
	ind, err := NewLinRegAng(timePeriod, selectData)
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())

	return ind, err
}

// NewDefaultLinRegAngWithSrcLen creates a Linear Regression Angle Indicator (LinRegAng) for offline usage with default parameters
func NewDefaultLinRegAngWithSrcLen(sourceLength int) (indicator *LinRegAng, err error) {
	ind, err := NewDefaultLinRegAng()
	ind.Data = make([]float64, 0, sourceLength-ind.GetLookbackPeriod())
	return ind, err
}

// NewLinRegAngForStream creates a Linear Regression Angle Indicator (LinRegAng) for online usage with a source data stream
func NewLinRegAngForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegAng, err error) {
	ind, err := NewLinRegAng(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLinRegAngForStream creates a Linear Regression Angle Indicator (LinRegAng) for online usage with a source data stream
func NewDefaultLinRegAngForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *LinRegAng, err error) {
	ind, err := NewDefaultLinRegAng()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewLinRegAngForStreamWithSrcLen creates a Linear Regression Angle Indicator (LinRegAng) for offline usage with a source data stream
func NewLinRegAngForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegAng, err error) {
	ind, err := NewLinRegAngWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLinRegAngForStreamWithSrcLen creates a Linear Regression Angle Indicator (LinRegAng) for offline usage with a source data stream
func NewDefaultLinRegAngForStreamWithSrcLen(sourceLength int, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *LinRegAng, err error) {
	ind, err := NewDefaultLinRegAngWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *LinRegAng) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}
