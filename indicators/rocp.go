package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type ROCPWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableAction
	periodCounter        int
	periodHistory        *list.List
}

func NewROCPWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *ROCPWithoutStorage, err error) {
	newROCP := ROCPWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod * -1),
		periodHistory:               list.New()}

	newROCP.selectData = selectData
	newROCP.valueAvailableAction = valueAvailableAction

	return &newROCP, err
}

// A Relative Strength Indicator
type ROCP struct {
	*ROCPWithoutStorage

	// public variables
	Data []float64
}

// NewROCP returns a new Rate of change percentage (ROCP) configured with the
// specified timePeriod. The ROCP results are stored in the DATA field.
func NewROCP(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *ROCP, err error) {
	newROCP := ROCP{}
	newROCP.ROCPWithoutStorage, err = NewROCPWithoutStorage(timePeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newROCP.Data = append(newROCP.Data, dataItem)
		})

	newROCP.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newROCP.Data = append(newROCP.Data, dataItem)
	}
	return &newROCP, err
}

func NewROCPForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *ROCP, err error) {
	newROCP, err := NewROCP(timePeriod, selectData)
	priceStream.AddTickSubscription(newROCP)
	return newROCP, err
}

func (ind *ROCPWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *ROCPWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

		//    ROCP = (price/previousPrice - 1) * 100
		previousPrice := ind.periodHistory.Front().Value.(float64)
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}
		var result float64
		if previousPrice != 0 {
			result = (tickData - previousPrice) / previousPrice
		} else {
			result = 0.0
		}

		if result > ind.maxValue {
			ind.maxValue = result
		}

		if result < ind.minValue {
			ind.minValue = result
		}

		ind.valueAvailableAction(result, streamBarIndex)
	}

	if ind.periodHistory.Len() > ind.GetTimePeriod() {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}
}
