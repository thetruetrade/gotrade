package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type ROCRWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	periodHistory        *list.List
}

func NewROCRWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *ROCRWithoutStorage, err error) {
	newROCR := ROCRWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod * -1),
		periodHistory:               list.New()}

	newROCR.valueAvailableAction = valueAvailableAction

	return &newROCR, err
}

// A Relative Strength Indicator
type ROCR struct {
	*ROCRWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewROCR returns a new Rate of change ratio (ROCR) configured with the
// specified timePeriod. The ROCR results are stored in the DATA field.
func NewROCR(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *ROCR, err error) {
	newROCR := ROCR{selectData: selectData}
	newROCR.ROCRWithoutStorage, err = NewROCRWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newROCR.Data = append(newROCR.Data, dataItem)
		})

	newROCR.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newROCR.Data = append(newROCR.Data, dataItem)
	}
	return &newROCR, err
}

func NewROCRForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *ROCR, err error) {
	newROCR, err := NewROCR(timePeriod, selectData)
	priceStream.AddTickSubscription(newROCR)
	return newROCR, err
}

func (ind *ROCR) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *ROCRWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

		//    ROCR = (price/previousPrice - 1) * 100
		previousPrice := ind.periodHistory.Front().Value.(float64)
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}
		var result float64
		if previousPrice != 0 {
			result = (tickData / previousPrice)
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
