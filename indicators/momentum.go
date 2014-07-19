package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type MomentumWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableAction
	periodCounter        int
	periodHistory        *list.List
}

func NewMomentumWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *MomentumWithoutStorage, err error) {
	newMomentum := MomentumWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod * -1),
		periodHistory:               list.New()}

	newMomentum.selectData = selectData
	newMomentum.valueAvailableAction = valueAvailableAction

	return &newMomentum, err
}

// A Momentum Indicator
type Momentum struct {
	*MomentumWithoutStorage

	// public variables
	Data []float64
}

// NewMomentum returns a new Momentum Indicator(Momentum) configured with the
// specified timePeriod. The Momentum results are stored in the DATA field.
func NewMomentum(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Momentum, err error) {
	newMomentum := Momentum{}
	newMomentum.MomentumWithoutStorage, err = NewMomentumWithoutStorage(timePeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newMomentum.Data = append(newMomentum.Data, dataItem)
		})

	newMomentum.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newMomentum.Data = append(newMomentum.Data, dataItem)
	}
	return &newMomentum, err
}

func NewMomentumForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *Momentum, err error) {
	newMomentum, err := NewMomentum(timePeriod, selectData)
	priceStream.AddTickSubscription(newMomentum)
	return newMomentum, err
}

func (ind *MomentumWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *MomentumWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter > 0 {

		// Momentum = price - previousPrice
		previousPrice := ind.periodHistory.Front().Value.(float64)
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}
		var result float64 = tickData - previousPrice

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
