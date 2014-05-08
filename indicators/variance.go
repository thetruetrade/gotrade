package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
)

type baseVariance struct {
	*baseIndicatorWithLookback

	// private variables
	periodCounter        int
	periodHistory        *list.List
	mean                 float64
	variance             float64
	valueAvailableAction ValueAvailableAction
}

func newBaseVariance(lookbackPeriod int) *baseVariance {
	ind := baseVariance{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod),
		periodCounter: 0,
		periodHistory: list.New(),
		mean:          0.0,
		variance:      0.0}
	return &ind
}

type Variance struct {
	*baseVariance

	// public variables
	Data []float64
}

func NewVariance(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *Variance, err error) {
	newVar := Variance{baseVariance: newBaseVariance(lookbackPeriod)}
	newVar.selectData = selectData
	newVar.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newVar.Data = append(newVar.Data, dataItem)
	}
	return &newVar, nil
}

func NewVarianceForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *Variance, err error) {
	newVar, err := NewVariance(lookbackPeriod, selectData)
	priceStream.AddTickSubscription(newVar)
	return newVar, err
}

func (ind *baseVariance) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = ind.selectData(tickData)
	ind.ReceiveTick(selectedData, streamBarIndex)
}

// http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance - Knuth
func (ind *baseVariance) ReceiveTick(tickData float64, streamBarIndex int) {
	ind.periodHistory.PushBack(tickData)
	firstValue := ind.periodHistory.Front().Value.(float64)

	previousMean := ind.mean
	previousVariance := ind.variance

	if ind.periodCounter < ind.LookbackPeriod {
		ind.periodCounter += 1
		delta := tickData - previousMean
		ind.mean = previousMean + delta/float64(ind.periodCounter)

		ind.variance = previousVariance + delta*(tickData-ind.mean)
	} else {
		delta := tickData - firstValue
		dOld := firstValue - previousMean
		ind.mean = previousMean + delta/float64(ind.periodCounter)
		dNew := tickData - ind.mean
		ind.variance = previousVariance + (dOld+dNew)*(delta)
	}

	if ind.periodHistory.Len() > ind.LookbackPeriod {
		first := ind.periodHistory.Front()
		ind.periodHistory.Remove(first)
	}

	if ind.periodCounter >= ind.LookbackPeriod {
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}

		result := ind.variance / float64(ind.LookbackPeriod)

		if result > ind.maxValue {
			ind.maxValue = result
		}

		if result < ind.minValue {
			ind.minValue = result
		}

		ind.valueAvailableAction(result, streamBarIndex)
	}
}
