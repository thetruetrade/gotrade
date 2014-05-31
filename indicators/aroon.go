package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

// Aroon Up = 100 x (25 - Days Since 25-day High)/25
// Aroon Down = 100 x (25 - Days Since 25-day Low)/25
type AroonWithoutStorage struct {
	*baseIndicator
	*baseIndicatorWithLookback
	*baseIndicatorWithTimePeriod

	// private variables
	periodCounter        int
	periodHighHistory    *list.List
	periodLowHistory     *list.List
	valueAvailableAction ValueAvailableActionAroon
	aroonFactor          float64
}

func NewAroonWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionAroon) (indicator *AroonWithoutStorage, err error) {
	ind := AroonWithoutStorage{baseIndicator: newBaseIndicator(),
		baseIndicatorWithLookback:   newBaseIndicatorWithLookback(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod + 1) * -1,
		periodHighHistory:           list.New(),
		periodLowHistory:            list.New()}
	ind.valueAvailableAction = valueAvailableAction
	ind.aroonFactor = 100.0 / float64(timePeriod)

	return &ind, nil
}

type Aroon struct {
	*AroonWithoutStorage

	Up   []float64
	Down []float64
}

func NewAroon(timePeriod int) (indicator *Aroon, err error) {
	newAroon := Aroon{}
	newAroon.AroonWithoutStorage, err = NewAroonWithoutStorage(timePeriod,
		func(dataItemAroonUp float64, dataItemAroonDown float64, streamBarIndex int) {
			newAroon.Up = append(newAroon.Up, dataItemAroonUp)
			newAroon.Down = append(newAroon.Down, dataItemAroonDown)
		})
	return &newAroon, err
}

func NewAroonForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *Aroon, err error) {
	ind, err := NewAroon(timePeriod)
	priceStream.AddTickSubscription(ind)
	return ind, err
}

func (ind *AroonWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHighHistory.PushBack(tickData.H())
	ind.periodLowHistory.PushBack(tickData.L())

	if ind.periodHighHistory.Len() > (1 + ind.GetLookbackPeriod()) {
		var first = ind.periodHighHistory.Front()
		ind.periodHighHistory.Remove(first)
		first = ind.periodLowHistory.Front()
		ind.periodLowHistory.Remove(first)
	}

	if ind.periodCounter >= 0 {
		ind.dataLength += 1

		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}

		var aroonUp float64
		var aroonDwn float64

		var highValue float64 = math.SmallestNonzeroFloat64
		var highIdx int = -1
		var i int = (1 + ind.GetLookbackPeriod())
		for e := ind.periodHighHistory.Front(); e != nil; e = e.Next() {
			i--
			var value float64 = e.Value.(float64)
			if highValue <= value {
				highValue = value
				highIdx = i
			}
		}
		var daysSinceHigh = highIdx

		var lowValue float64 = math.MaxFloat64
		var lowIdx int = -1
		i = (1 + ind.GetLookbackPeriod())
		for e := ind.periodLowHistory.Front(); e != nil; e = e.Next() {
			i--
			var value float64 = e.Value.(float64)
			if lowValue >= value {
				lowValue = value
				lowIdx = i
			}

		}
		var daysSinceLow = lowIdx

		aroonUp = ind.aroonFactor * float64(ind.GetLookbackPeriod()-daysSinceHigh)
		aroonDwn = ind.aroonFactor * float64(ind.GetLookbackPeriod()-daysSinceLow)
		if aroonUp > ind.maxValue {
			ind.maxValue = aroonUp
		}

		if aroonDwn < ind.minValue {
			ind.minValue = aroonDwn
		}

		ind.valueAvailableAction(aroonUp, aroonDwn, streamBarIndex)
	}

}
