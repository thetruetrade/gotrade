package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type ROCR100WithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	periodCounter        int
	periodHistory        *list.List
}

func NewROCR100WithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableActionFloat) (indicator *ROCR100WithoutStorage, err error) {
	newROCR100 := ROCR100WithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod * -1),
		periodHistory:               list.New()}

	newROCR100.selectData = selectData
	newROCR100.valueAvailableAction = valueAvailableAction

	return &newROCR100, err
}

// A Relative Strength Indicator
type ROCR100 struct {
	*ROCR100WithoutStorage

	// public variables
	Data []float64
}

// NewROCR100 returns a new Rate of change ratio (100 scale)(ROCR100) configured with the
// specified timePeriod. The ROCR100 results are stored in the DATA field.
func NewROCR100(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *ROCR100, err error) {
	newROCR100 := ROCR100{}
	newROCR100.ROCR100WithoutStorage, err = NewROCR100WithoutStorage(timePeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newROCR100.Data = append(newROCR100.Data, dataItem)
		})

	newROCR100.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newROCR100.Data = append(newROCR100.Data, dataItem)
	}
	return &newROCR100, err
}

func NewROCR100ForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *ROCR100, err error) {
	newROCR100, err := NewROCR100(timePeriod, selectData)
	priceStream.AddTickSubscription(newROCR100)
	return newROCR100, err
}

func (ind *ROCR100WithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *ROCR100WithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

		//    ROCR100 = (price/previousPrice - 1) * 100
		previousPrice := ind.periodHistory.Front().Value.(float64)
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}
		var result float64
		if previousPrice != 0 {
			result = (tickData / previousPrice) * 100.0
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
