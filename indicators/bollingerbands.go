package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

type BollingerBandEntry struct {
	UpperBand  float64
	MiddleBand float64
	LowerBand  float64
}

type bollingerBandsBase struct {
	Indicator

	// public variables
	LookbackPeriod int

	// private variables
	periodCounter        int
	periodHistory        *list.List
	mean                 float64
	variance             float64
	validFromBarIndex    int
	dataLength           int
	valueAvailableAction ValueAvailableActionBollinger
	transformData        gotrade.DataTransformationFunc
}

type BollingerBands struct {
	bollingerBandsBase

	Data []BollingerBandEntry
}

type BollingerBandsForAttachment struct {
	bollingerBandsBase
}

func NewBollingerBands(lookbackPeriod int, transformData gotrade.DataTransformationFunc) (indicator *BollingerBands, err error) {
	newBB := BollingerBands{bollingerBandsBase: bollingerBandsBase{LookbackPeriod: lookbackPeriod,
		periodCounter: 0.0,
		mean:          0.0,
		variance:      0.0,
		transformData: transformData,
		periodHistory: list.New(),
		Indicator:     Indicator{validFromBar: -1, minValue: math.MaxFloat64, maxValue: math.SmallestNonzeroFloat64}}}
	newBB.valueAvailableAction = func(dataItem BollingerBandEntry, streamBarIndex int) {
		newBB.Data = append(newBB.Data, dataItem)
	}
	return &newBB, nil
}

func NewBollingerBandsForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, transformData gotrade.DataTransformationFunc) (indicator *BollingerBands, err error) {
	bb, err := NewBollingerBands(lookbackPeriod, transformData)
	priceStream.AddSubscription(bb)
	return bb, err
}

// http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance - Knuth
func (bb *bollingerBandsBase) RecieveOrderedTick(dataItem gotrade.DOHLCV, streamBarIndex int) {
	bb.dataLength += 1

	var transformedData = bb.transformData(dataItem)
	bb.periodHistory.PushBack(transformedData)
	var firstValue = bb.periodHistory.Front().Value.(float64)

	previousMean := bb.mean
	previousVariance := bb.variance
	standardDeviation := 0.0
	if bb.periodCounter < bb.LookbackPeriod {
		bb.periodCounter += 1
		delta := transformedData - previousMean
		bb.mean = previousMean + delta/float64(bb.periodCounter)

		bb.variance = previousVariance + delta*(transformedData-bb.mean)
		standardDeviation = math.Sqrt(bb.variance / (float64(bb.periodCounter)))
	} else {
		delta := transformedData - firstValue
		dOld := firstValue - previousMean
		bb.mean = previousMean + delta/float64(bb.periodCounter)
		dNew := transformedData - bb.mean
		bb.variance = previousVariance + (dOld+dNew)*(delta)
		standardDeviation = math.Sqrt(bb.variance / (float64(bb.periodCounter)))

	}

	var upperBand = bb.mean + 2*standardDeviation
	var lowerBand = bb.mean - 2*standardDeviation

	if upperBand > bb.maxValue {
		bb.maxValue = upperBand
	}

	if lowerBand < bb.minValue {
		bb.minValue = lowerBand
	}

	if bb.periodHistory.Len() > bb.LookbackPeriod {
		var first = bb.periodHistory.Front()
		bb.periodHistory.Remove(first)
	}

	if bb.periodCounter >= bb.LookbackPeriod {
		if bb.validFromBar == -1 {
			bb.validFromBar = streamBarIndex
		}
		bb.valueAvailableAction(BollingerBandEntry{UpperBand: upperBand, MiddleBand: bb.mean, LowerBand: lowerBand}, streamBarIndex)
	}
}
