package indicators

// import (
// 	"container/list"
// 	"github.com/thetruetrade/gotrade"
// 	"math"
// )

// type BollingerBandEntry struct {
// 	UpperBand  float64
// 	MiddleBand float64
// 	LowerBand  float64
// }

// type baseBollingerBands struct {
// 	*baseIndicatorWithLookback

// 	// private variables
// 	periodCounter        int
// 	periodHistory        *list.List
// 	previousMean         float64
// 	variance             float64
// 	valueAvailableAction ValueAvailableActionBollinger
// 	sma                  *SMAWithoutStorage
// 	mean                 float64
// }

// func newBaseBollingerBands(lookbackPeriod int) *baseBollingerBands {
// 	ind := baseBollingerBands{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod),
// 		periodCounter: 0,
// 		periodHistory: list.New(),
// 		previousMean:  0.0,
// 		variance:      0.0}
// 	return &ind
// }

// type BollingerBands struct {
// 	*baseBollingerBands

// 	Data []BollingerBandEntry
// }

// func NewBollingerBands(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
// 	newBB := BollingerBands{baseBollingerBands: newBaseBollingerBands(lookbackPeriod)}
// 	newBB.selectData = selectData
// 	newBB.valueAvailableAction = func(dataItem BollingerBandEntry, streamBarIndex int) {
// 		newBB.Data = append(newBB.Data, dataItem)
// 	}
// 	newBB.sma, _ = NewSMAWithoutStorage(lookbackPeriod, selectData, func(dataItem float64, streamBarIndex int) {
// 		if newBB.periodCounter >= newBB.LookbackPeriod {
// 			// store the previous mean once we have one
// 			newBB.previousMean = newBB.mean
// 		}

// 		newBB.mean = dataItem
// 	})
// 	return &newBB, nil
// }

// func NewBollingerBandsForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
// 	bb, err := NewBollingerBands(lookbackPeriod, selectData)
// 	priceStream.AddTickSubscription(bb)
// 	return bb, err
// }

// func (bb *baseBollingerBands) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
// 	var selectedData float64 = bb.selectData(tickData)

// 	bb.sma.ReceiveTick(selectedData, streamBarIndex)
// 	bb.RecieveTick(selectedData, streamBarIndex)

// }

// // http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance - Knuth
// func (bb *baseBollingerBands) RecieveTick(tickData float64, streamBarIndex int) {
// 	bb.dataLength += 1

// 	bb.periodHistory.PushBack(tickData)
// 	var firstValue = bb.periodHistory.Front().Value.(float64)

// 	previousMean := bb.previousMean
// 	previousVariance := bb.variance
// 	standardDeviation := 0.0
// 	bb.periodCounter++
// 	if bb.periodCounter >= bb.LookbackPeriod {
// 		delta := tickData - firstValue
// 		dOld := firstValue - previousMean

// 		dNew := tickData - bb.mean
// 		bb.variance = previousVariance + (dOld+dNew)*(delta)
// 		standardDeviation = math.Sqrt(bb.variance / (float64(bb.periodCounter)))
// 	}

// 	if bb.periodHistory.Len() > bb.LookbackPeriod {
// 		var first = bb.periodHistory.Front()
// 		bb.periodHistory.Remove(first)
// 	}

// 	if bb.periodCounter >= bb.LookbackPeriod {
// 		if bb.validFromBar == -1 {
// 			bb.validFromBar = streamBarIndex
// 		}

// 		var upperBand = bb.mean + 2*standardDeviation
// 		var lowerBand = bb.mean - 2*standardDeviation

// 		if upperBand > bb.maxValue {
// 			bb.maxValue = upperBand
// 		}

// 		if lowerBand < bb.minValue {
// 			bb.minValue = lowerBand
// 		}

// 		bb.valueAvailableAction(BollingerBandEntry{UpperBand: upperBand, MiddleBand: bb.mean, LowerBand: lowerBand}, streamBarIndex)
// 	}
// }
