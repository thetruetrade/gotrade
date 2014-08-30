package indicators

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Linear Regression Angle Indicator (LinRegAng)
type LinRegAng struct {
	*LinRegWithoutStorage
	selectData gotrade.DOHLCVDataSelectionFunc

	// public variables
	Data []float64
}

// NewLinRegAng creates a Linear Regression Angle Indicator (LinRegAng) for online usage
func NewLinRegAng(timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *LinRegAng, err error) {
	if selectData == nil {
		return nil, ErrDOHLCVDataSelectFuncIsNil
	}

	ind := LinRegAng{
		selectData: selectData,
	}

	ind.LinRegWithoutStorage, err = NewLinRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			result := math.Atan(slope) * (180.0 / math.Pi)

			ind.UpdateMinMax(result, result)

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
func NewLinRegAngWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *LinRegAng, err error) {
	ind, err := NewLinRegAng(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultLinRegAngWithSrcLen creates a Linear Regression Angle Indicator (LinRegAng) for offline usage with default parameters
func NewDefaultLinRegAngWithSrcLen(sourceLength uint) (indicator *LinRegAng, err error) {
	ind, err := NewDefaultLinRegAng()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewLinRegAngForStream creates a Linear Regression Angle Indicator (LinRegAng) for online usage with a source data stream
func NewLinRegAngForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *LinRegAng, err error) {
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
func NewLinRegAngForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DOHLCVDataSelectionFunc) (indicator *LinRegAng, err error) {
	ind, err := NewLinRegAngWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLinRegAngForStreamWithSrcLen creates a Linear Regression Angle Indicator (LinRegAng) for offline usage with a source data stream
func NewDefaultLinRegAngForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *LinRegAng, err error) {
	ind, err := NewDefaultLinRegAngWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *LinRegAng) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}
