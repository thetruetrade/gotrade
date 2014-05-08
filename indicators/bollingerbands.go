package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type BollingerBandEntry struct {
	UpperBand  float64
	MiddleBand float64
	LowerBand  float64
}

type baseBollingerBands struct {
	*baseIndicatorWithLookback

	// private variables
	valueAvailableAction ValueAvailableActionBollinger
	sma                  *SMAWithoutStorage
	stdDev               *StdDeviation
	currentSMA           float64
}

func newBaseBollingerBands(lookbackPeriod int) *baseBollingerBands {
	ind := baseBollingerBands{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod), currentSMA: 0.0}
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
	newBB.sma, _ = NewSMAWithoutStorage(lookbackPeriod, selectData, func(dataItem float64, streamBarIndex int) {
		newBB.currentSMA = dataItem
	})

	newBB.stdDev, _ = NewStdDeviation(lookbackPeriod, selectData)
	newBB.stdDev.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newBB.dataLength += 1
		if newBB.validFromBar == -1 {
			newBB.validFromBar = streamBarIndex
		}

		var upperBand = newBB.currentSMA + 2*dataItem
		var lowerBand = newBB.currentSMA - 2*dataItem

		if upperBand > newBB.maxValue {
			newBB.maxValue = upperBand
		}

		if lowerBand < newBB.minValue {
			newBB.minValue = lowerBand
		}

		newBB.valueAvailableAction(BollingerBandEntry{UpperBand: upperBand, MiddleBand: newBB.currentSMA, LowerBand: lowerBand}, streamBarIndex)
	}

	return &newBB, nil
}

func NewBollingerBandsForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
	bb, err := NewBollingerBands(lookbackPeriod, selectData)
	priceStream.AddTickSubscription(bb)
	return bb, err
}

func (bb *baseBollingerBands) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData float64 = bb.selectData(tickData)
	bb.RecieveTick(selectedData, streamBarIndex)
}

// http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance - Knuth
func (bb *baseBollingerBands) RecieveTick(tickData float64, streamBarIndex int) {
	bb.sma.ReceiveTick(tickData, streamBarIndex)
	bb.stdDev.ReceiveTick(tickData, streamBarIndex)
}

type BollingerDataSelectionFunc func(dataItem BollingerBandEntry) float64

func UseUpperBand(dataItem BollingerBandEntry) float64 {
	return dataItem.UpperBand
}

func UseMiddleBand(dataItem BollingerBandEntry) float64 {
	return dataItem.MiddleBand
}

func UseLowerBand(dataItem BollingerBandEntry) float64 {
	return dataItem.LowerBand
}
