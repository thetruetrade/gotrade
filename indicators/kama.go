package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

// A Kaufman Adaptive Moving Average Indicator
type KAMAWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	periodTotal          float64
	periodHistory        *list.List
	periodCounter        int
	constantMax          float64
	constantDiff         float64
	sumROC               float64
	periodROC            float64
	previousClose        float64
	previousKAMA         float64
	valueAvailableAction ValueAvailableActionFloat
}

// NewKAMAWithoutStorage returns a new Kaufman Adaptive Moving Average (KAMA) configured with the
// specified timePeriod, this version is intended for use by other indicators.
// The KAMA results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewKAMAWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *KAMAWithoutStorage, err error) {
	newKAMA := KAMAWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod),
		periodCounter:               (timePeriod + 1) * -1,
		constantMax:                 float64(2.0 / (30.0 + 1.0)),
		constantDiff:                float64((2.0 / (2.0 + 1.0)) - (2.0 / (30.0 + 1.0))),
		sumROC:                      0.0,
		periodROC:                   0.0,
		periodHistory:               list.New(),
		previousClose:               math.SmallestNonzeroFloat64}
	newKAMA.valueAvailableAction = valueAvailableAction

	return &newKAMA, nil
}

// A Simple Moving Average Indicator
type KAMA struct {
	*KAMAWithoutStorage
	selectData gotrade.DataSelectionFunc

	// public variables
	Data []float64
}

// NewKAMA returns a new Simple Moving Average (KAMA) configured with the
// specified timePeriod. The KAMA results are stored in the Data field.
func NewKAMA(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *KAMA, err error) {
	newKAMA := KAMA{selectData: selectData}
	newKAMA.KAMAWithoutStorage, err = NewKAMAWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		newKAMA.Data = append(newKAMA.Data, dataItem)
	})

	return &newKAMA, err
}

func NewKAMAForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *KAMA, err error) {
	newKama, err := NewKAMA(timePeriod, selectData)
	priceStream.AddTickSubscription(newKama)
	return newKama, err
}

func (ind *KAMA) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

func (ind *KAMAWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodCounter += 1
	ind.periodHistory.PushBack(tickData)

	if ind.periodCounter <= 0 {
		if ind.previousClose > math.SmallestNonzeroFloat64 {
			ind.sumROC += math.Abs(tickData - ind.previousClose)
		}
	}
	if ind.periodCounter == 0 {
		var er float64 = 0.0
		var sc float64 = 0.0
		var closeMinusN float64 = ind.periodHistory.Front().Value.(float64)
		ind.previousKAMA = ind.previousClose
		ind.periodROC = tickData - closeMinusN

		// calculate the efficiency ratio
		if ind.sumROC <= ind.periodROC || isZero(ind.sumROC) {
			er = 1.0
		} else {
			er = math.Abs(ind.periodROC / ind.sumROC)
		}

		sc = (er * ind.constantDiff) + ind.constantMax
		sc *= sc
		ind.previousKAMA = ((tickData - ind.previousKAMA) * sc) + ind.previousKAMA

		result := ind.previousKAMA

		ind.dataLength += 1

		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}

		if result > ind.maxValue {
			ind.maxValue = result
		}

		if result < ind.minValue {
			ind.minValue = result
		}

		ind.valueAvailableAction(result, streamBarIndex)

	} else if ind.periodCounter > 0 {

		var er float64 = 0.0
		var sc float64 = 0.0
		var closeMinusN float64 = ind.periodHistory.Front().Value.(float64)
		var closeMinusN1 float64 = ind.periodHistory.Front().Next().Value.(float64)
		ind.periodROC = tickData - closeMinusN1

		ind.sumROC -= math.Abs(closeMinusN1 - closeMinusN)
		ind.sumROC += math.Abs(tickData - ind.previousClose)

		// calculate the efficiency ratio
		if ind.sumROC <= ind.periodROC || isZero(ind.sumROC) {
			er = 1.0
		} else {
			er = math.Abs(ind.periodROC / ind.sumROC)
		}

		sc = (er * ind.constantDiff) + ind.constantMax
		sc *= sc
		ind.previousKAMA = ((tickData - ind.previousKAMA) * sc) + ind.previousKAMA

		result := ind.previousKAMA

		ind.dataLength += 1

		if result > ind.maxValue {
			ind.maxValue = result
		}

		if result < ind.minValue {
			ind.minValue = result
		}

		ind.valueAvailableAction(result, streamBarIndex)
	}

	ind.previousClose = tickData

	if ind.periodHistory.Len() > (ind.GetTimePeriod() + 1) {
		var first = ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}

}

func isZero(value float64) bool {
	var epsilon float64 = 0.00000000000001
	return (((-epsilon) < value) && (value < epsilon))
}
