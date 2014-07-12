package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type ROCWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableAction
	periodCounter        int
	periodHistory        *list.List
}

func NewROCWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *ROCWithoutStorage, err error) {
	newROC := ROCWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod * -1),
		periodHistory:               list.New()}

	newROC.selectData = selectData
	newROC.valueAvailableAction = valueAvailableAction

	return &newROC, err
}

// A Relative Strength Indicator
type ROC struct {
	*ROCWithoutStorage

	// public variables
	Data []float64
}

// NewROC returns a new Rate of change (ROC) configured with the
// specified timePeriod. The ROC results are stored in the DATA field.
func NewROC(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *ROC, err error) {
	newROC := ROC{}
	newROC.ROCWithoutStorage, err = NewROCWithoutStorage(timePeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newROC.Data = append(newROC.Data, dataItem)
		})

	newROC.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newROC.Data = append(newROC.Data, dataItem)
	}
	return &newROC, err
}

func NewROCForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *ROC, err error) {
	newROC, err := NewROC(timePeriod, selectData)
	priceStream.AddTickSubscription(newROC)
	return newROC, err
}

func (ind *ROCWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *ROCWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

		//    ROC = (price/previousPrice - 1) * 100
		previousPrice := ind.periodHistory.Front().Value.(float64)
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}
		var result float64
		if previousPrice != 0 {
			result = 100.0 * ((tickData / previousPrice) - 1)
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
