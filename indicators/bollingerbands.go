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

type baseBollingerBands struct {
	*baseIndicatorWithLookback

	// private variables
	periodCounter        int
	periodHistory        *list.List
	mean                 float64
	variance             float64
	valueAvailableAction ValueAvailableActionBollinger
}

func newBaseBollingerBands(lookbackPeriod int) *baseBollingerBands {
	ind := baseBollingerBands{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod),
		periodCounter: 0,
		periodHistory: list.New(),
		mean:          0.0,
		variance:      0.0}
	return &ind
}

type BollingerBands struct {
	*baseBollingerBands

	Data []BollingerBandEntry
}

func NewBollingerBands(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
	newBB := BollingerBands{baseBollingerBands: newBaseBollingerBands(lookbackPeriod)}
	newBB.selectData = selectData
	newBB.valueAvailableAction = func(dataItem BollingerBandEntry, streamBarIndex int) {
		newBB.Data = append(newBB.Data, dataItem)
	}
	return &newBB, nil
}

func NewBollingerBandsForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
	bb, err := NewBollingerBands(lookbackPeriod, selectData)
	priceStream.AddTickSubscription(bb)
	return bb, err
}

func (bb *baseBollingerBands) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = bb.selectData(tickData)
	bb.RecieveTick(selectedData, streamBarIndex)
}

// http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance - Knuth
func (bb *baseBollingerBands) RecieveTick(tickData float64, streamBarIndex int) {
	bb.dataLength += 1

	bb.periodHistory.PushBack(tickData)
	var firstValue = bb.periodHistory.Front().Value.(float64)

	previousMean := bb.mean
	previousVariance := bb.variance
	standardDeviation := 0.0
	if bb.periodCounter < bb.LookbackPeriod {
		bb.periodCounter += 1
		delta := tickData - previousMean
		bb.mean = previousMean + delta/float64(bb.periodCounter)

		bb.variance = previousVariance + delta*(tickData-bb.mean)
		standardDeviation = math.Sqrt(bb.variance / (float64(bb.periodCounter)))
	} else {
		delta := tickData - firstValue
		dOld := firstValue - previousMean
		bb.mean = previousMean + delta/float64(bb.periodCounter)
		dNew := tickData - bb.mean
		bb.variance = previousVariance + (dOld+dNew)*(delta)
		standardDeviation = math.Sqrt(bb.variance / (float64(bb.periodCounter)))
	}

	if bb.periodHistory.Len() > bb.LookbackPeriod {
		var first = bb.periodHistory.Front()
		bb.periodHistory.Remove(first)
	}

	if bb.periodCounter >= bb.LookbackPeriod {
		if bb.validFromBar == -1 {
			bb.validFromBar = streamBarIndex
		}

		var upperBand = bb.mean + 2*standardDeviation
		var lowerBand = bb.mean - 2*standardDeviation

		if upperBand > bb.maxValue {
			bb.maxValue = upperBand
		}

		if lowerBand < bb.minValue {
			bb.minValue = lowerBand
		}

		bb.valueAvailableAction(BollingerBandEntry{UpperBand: upperBand, MiddleBand: bb.mean, LowerBand: lowerBand}, streamBarIndex)
	}
}
