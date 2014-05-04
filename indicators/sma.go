// Simple Moving Average (SMA)
package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

type smaBase struct {
	Indicator

	// public variables
	LookbackPeriod int

	// private variables
	periodTotal   float64
	periodHistory *list.List
	periodCounter int
}

// A Simple Moving Average Indicator
type SMA struct {
	smaBase

	// public variables
	Data []float64
}
type SMAForAttachment struct {
	smaBase
}

// NewSMA returns a new Simple Moving Average (SMA) configured with the
// specified lookbackPeriod. The SMA results are stored in the DATA field.
func NewSMA(lookbackPeriod int, transformData gotrade.DataTransformationFunc) (indicator *SMA, err error) {
	newSMA := SMA{smaBase: smaBase{LookbackPeriod: lookbackPeriod,
		periodCounter: lookbackPeriod * -1,
		periodHistory: list.New(),
		Indicator:     Indicator{validFromBar: -1, transformData: transformData, minValue: math.MaxFloat64, maxValue: math.SmallestNonzeroFloat64}}}

	newSMA.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newSMA.Data = append(newSMA.Data, dataItem)
	}
	return &newSMA, nil
}

func NewSMAForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, transformData gotrade.DataTransformationFunc) (indicator *SMA, err error) {
	newSma, err := NewSMA(lookbackPeriod, transformData)
	priceStream.AddSubscription(newSma)
	return newSma, err
}

// NewAttachedSMA returns a new Simple Moving Average (SMA) configured with the
// specified lookbackPeriod, this version is intended for use by other indicators.
// The SMA results are not stored in a local field but made available though the
// configured valueAvailableAction for storage by the parent indicator.
func NewAttachedSMA(lookbackPeriod int,
	transformData gotrade.DataTransformationFunc,
	valueAvailableAction ValueAvailableAction) (indicator *SMAForAttachment, err error) {
	newSMA := SMAForAttachment{smaBase{LookbackPeriod: lookbackPeriod,
		periodCounter: lookbackPeriod * -1,
		periodHistory: list.New(),
		Indicator:     Indicator{validFromBar: -1, transformData: transformData, valueAvailableAction: valueAvailableAction, minValue: math.MaxFloat64, maxValue: math.SmallestNonzeroFloat64}}}

	return &newSMA, nil
}

func (sma *smaBase) RecieveOrderedTick(dataItem gotrade.DOHLCV, streamBarIndex int) {
	sma.periodCounter += 1
	sma.dataLength += 1
	var transformedData = sma.transformData(dataItem)

	if transformedData > sma.maxValue {
		sma.maxValue = transformedData
	}

	if transformedData < sma.minValue {
		sma.minValue = transformedData
	}

	sma.periodHistory.PushBack(transformedData)

	if sma.periodCounter > 0 {
		var valueToRemove = sma.periodHistory.Front()
		sma.periodTotal -= valueToRemove.Value.(float64)
	}
	if sma.periodHistory.Len() > sma.LookbackPeriod {
		var first = sma.periodHistory.Front()
		sma.periodHistory.Remove(first)
	}
	sma.periodTotal += transformedData
	var result float64 = sma.periodTotal / float64(sma.LookbackPeriod)
	if sma.periodCounter >= 0 {
		if sma.validFromBar == -1 {
			sma.validFromBar = streamBarIndex
		}
		sma.valueAvailableAction(result, streamBarIndex)
	}
}
