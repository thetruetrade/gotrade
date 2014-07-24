package indicators

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

type LinearRegAngle struct {
	*LinearRegWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

func NewLinearRegAngle(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinearRegAngle, err error) {
	newInd := LinearRegAngle{selectData: selectData}
	newInd.LinearRegWithoutStorage, err = NewLinearRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			result := math.Atan(slope) * (180.0 / math.Pi)

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

func NewLinearRegAngleForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinearRegAngle, err error) {
	newInd, err := NewLinearRegAngle(timePeriod, selectData)
	priceStream.AddTickSubscription(newInd)
	return newInd, err
}

func (ind *LinearRegAngle) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}
