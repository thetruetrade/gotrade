package indicators

import (
	"github.com/thetruetrade/gotrade"
)

// A Linear Regression Intercept Indicator (LinRegInt)
type LinRegSlp struct {
	*LinRegWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewLinRegSlp creates a Linear Regression Slope Indicator (LinRegSlp) for online usage
func NewLinRegSlp(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegSlp, err error) {
	ind := LinRegSlp{selectData: selectData}
	ind.LinRegWithoutStorage, err = NewLinRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			result := slope

			ind.UpdateMinMax(result, result)

			ind.Data = append(ind.Data, result)
		})

	return &ind, err
}

// NewDefaultLinRegSlp creates a Linear Regression Slope Indicator (LinRegSlp) for online usage with default parameters
//	- timePeriod: 14
func NewDefaultLinRegSlp() (indicator *LinRegSlp, err error) {
	timePeriod := 14
	return NewLinRegSlp(timePeriod, gotrade.UseClosePrice)
}

// NewLinRegSlpWithSrcLen creates a Linear Regression Slope Indicator (LinRegSlp) for offline usage
func NewLinRegSlpWithSrcLen(sourceLength uint, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegSlp, err error) {
	ind, err := NewLinRegSlp(timePeriod, selectData)

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewDefaultLinRegSlpWithSrcLen creates a Linear Regression Slope Indicator (LinRegSlp) for offline usage with default parameters
func NewDefaultLinRegSlpWithSrcLen(sourceLength uint) (indicator *LinRegSlp, err error) {
	ind, err := NewDefaultLinRegSlp()

	// only initialise the storage if there is enough source data to require it
	if sourceLength-uint(ind.GetLookbackPeriod()) > 1 {
		ind.Data = make([]float64, 0, sourceLength-uint(ind.GetLookbackPeriod()))
	}

	return ind, err
}

// NewLinRegSlpForStream creates a Linear Regression Slope Indicator (LinRegSlp) for online usage with a source data stream
func NewLinRegSlpForStream(priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegSlp, err error) {
	ind, err := NewLinRegSlp(timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLinRegSlpForStream creates a Linear Regression Slope Indicator (LinRegSlp) for online usage with a source data stream
func NewDefaultLinRegSlpForStream(priceStream gotrade.DOHLCVStreamSubscriber) (indicator *LinRegSlp, err error) {
	ind, err := NewDefaultLinRegSlp()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewLinRegSlpForStreamWithSrcLen creates a Linear Regression Slope Indicator (LinRegSlp) for offline usage with a source data stream
func NewLinRegSlpForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinRegSlp, err error) {
	ind, err := NewLinRegSlpWithSrcLen(sourceLength, timePeriod, selectData)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// NewDefaultLinRegSlpForStreamWithSrcLen creates a Linear Regression Slope Indicator (LinRegSlp) for offline usage with a source data stream
func NewDefaultLinRegSlpForStreamWithSrcLen(sourceLength uint, priceStream gotrade.DOHLCVStreamSubscriber) (indicator *LinRegSlp, err error) {
	ind, err := NewDefaultLinRegSlpWithSrcLen(sourceLength)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

// ReceiveDOHLCVTick consumes a source data DOHLCV price tick
func (ind *LinRegSlp) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}
