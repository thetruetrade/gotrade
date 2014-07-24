package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type LinearRegSlope struct {
	*LinearRegWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

func NewLinearRegSlope(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinearRegSlope, err error) {
	newInd := LinearRegSlope{selectData: selectData}
	newInd.LinearRegWithoutStorage, err = NewLinearRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			result := slope

			if result > newInd.LinearRegWithoutStorage.maxValue {
				newInd.LinearRegWithoutStorage.maxValue = result
			}

			if result < newInd.LinearRegWithoutStorage.minValue {
				newInd.LinearRegWithoutStorage.minValue = result
			}

			newInd.Data = append(newInd.Data, result)
		})

	return &newInd, err
}

func NewLinearRegSlopeForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinearRegSlope, err error) {
	newInd, err := NewLinearRegSlope(timePeriod, selectData)
	priceStream.AddTickSubscription(newInd)
	return newInd, err
}

func (ind *LinearRegSlope) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}
