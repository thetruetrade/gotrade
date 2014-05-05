// Weighted Moving Average (WMA)
package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

type wmaBase struct {
	Indicator

	// public variables
	LookbackPeriod int

	// private variables
	periodTotal       float64
	periodHistory     *list.List
	periodCounter     int
	periodWeightTotal int
}

// A Simple Moving Average Indicator
type WMA struct {
	wmaBase

	// public variables
	Data []float64
}
type WMAForAttachment struct {
	wmaBase
}

// NewWMA returns a new Simple Moving Average (WMA) configured with the
// specified lookbackPeriod. The WMA results are stored in the DATA field.
func NewWMA(lookbackPeriod int, transformData gotrade.DataTransformationFunc) (indicator *WMA, err error) {
	newWMA := WMA{wmaBase: wmaBase{LookbackPeriod: lookbackPeriod,
		periodCounter: lookbackPeriod * -1,
		periodHistory: list.New(),
		Indicator:     Indicator{validFromBar: -1, transformData: transformData, minValue: math.MaxFloat64, maxValue: math.SmallestNonzeroFloat64}}}

	var weightedTotal int = 0
	for i := 1; i <= lookbackPeriod; i++ {
		weightedTotal += i
	}
	newWMA.periodWeightTotal = weightedTotal
	newWMA.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newWMA.Data = append(newWMA.Data, dataItem)
	}
	return &newWMA, nil
}

func NewWMAForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, transformData gotrade.DataTransformationFunc) (indicator *WMA, err error) {
	newSma, err := NewWMA(lookbackPeriod, transformData)
	priceStream.AddSubscription(newSma)
	return newSma, err
}

// NewAttachedWMA returns a new Simple Moving Average (WMA) configured with the
// specified lookbackPeriod, this version is intended for use by other indicators.
// The WMA results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewAttachedWMA(lookbackPeriod int,
	transformData gotrade.DataTransformationFunc,
	valueAvailableAction ValueAvailableAction) (indicator *WMAForAttachment, err error) {
	newWMA := WMAForAttachment{wmaBase{LookbackPeriod: lookbackPeriod,
		periodCounter: lookbackPeriod * -1,
		periodHistory: list.New(),
		Indicator:     Indicator{validFromBar: -1, transformData: transformData, valueAvailableAction: valueAvailableAction, minValue: math.MaxFloat64, maxValue: math.SmallestNonzeroFloat64}}}

	return &newWMA, nil
}

func (wma *wmaBase) RecieveOrderedTick(dataItem gotrade.DOHLCV, streamBarIndex int) {
	wma.periodCounter += 1
	wma.dataLength += 1
	var transformedData = wma.transformData(dataItem)

	if transformedData > wma.maxValue {
		wma.maxValue = transformedData
	}

	if transformedData < wma.minValue {
		wma.minValue = transformedData
	}

	wma.periodHistory.PushBack(transformedData)

	if wma.periodCounter > 0 {
		var valueToRemove = wma.periodHistory.Front()
		wma.periodTotal -= valueToRemove.Value.(float64)
	}
	if wma.periodHistory.Len() > wma.LookbackPeriod {
		var first = wma.periodHistory.Front()
		wma.periodHistory.Remove(first)
	}
	wma.periodTotal += transformedData

	if wma.periodCounter >= 0 {
		if wma.validFromBar == -1 {
			wma.validFromBar = streamBarIndex
		}

		// calculate the wma
		var iter int = 1
		var sum float64 = 0
		for e := wma.periodHistory.Front(); e != nil; e = e.Next() {
			var localSum float64 = 0
			for i := 1; i <= iter; i++ {
				localSum += e.Value.(float64)
			}
			sum += localSum
			iter++
		}
		var result float64 = sum / float64(wma.periodWeightTotal)
		wma.valueAvailableAction(result, streamBarIndex)
	}
}
