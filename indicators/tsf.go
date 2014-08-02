package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type TSF struct {
	*LinRegWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

func NewTSF(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *TSF, err error) {
	newInd := TSF{selectData: selectData}
	newInd.LinRegWithoutStorage, err = NewLinRegWithoutStorage(timePeriod,
		func(dataItem float64, slope float64, intercept float64, streamBarIndex int) {
			result := intercept + slope*float64(timePeriod)

			if result > newInd.LinRegWithoutStorage.maxValue {
				newInd.LinRegWithoutStorage.maxValue = result
			}

			if result < newInd.LinRegWithoutStorage.minValue {
				newInd.LinRegWithoutStorage.minValue = result
			}

			newInd.Data = append(newInd.Data, result)
		})

	return &newInd, err
}

func NewTSFForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *TSF, err error) {
	newInd, err := NewTSF(timePeriod, selectData)
	priceStream.AddTickSubscription(newInd)
	return newInd, err
}

func (ind *TSF) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}
