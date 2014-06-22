package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type LinearRegSlope struct {
	*LinearRegWithoutStorage

	// public variables
	Data []float64
}

func NewLinearRegSlope(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *LinearRegSlope, err error) {
	newInd := LinearRegSlope{}
	newInd.LinearRegWithoutStorage, err = NewLinearRegWithoutStorage(timePeriod, selectData,
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
