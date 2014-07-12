package indicators

import (
	"github.com/thetruetrade/gotrade"
)

type BollingerBands struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableActionBollinger
	sma                  *SMAWithoutStorage
	stdDev               *StdDeviation
	currentSMA           float64

	UpperBand  []float64
	MiddleBand []float64
	LowerBand  []float64
}

func NewBollingerBands(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
	newBB := BollingerBands{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod)}
	newBB.currentSMA = 0.0
	newBB.selectData = selectData
	newBB.valueAvailableAction = func(dataItemUpperBand float64, dataItemMiddleBand float64, dataItemLowerBand float64, streamBarIndex int) {
		newBB.UpperBand = append(newBB.UpperBand, dataItemUpperBand)
		newBB.MiddleBand = append(newBB.MiddleBand, dataItemMiddleBand)
		newBB.LowerBand = append(newBB.LowerBand, dataItemLowerBand)
	}
	newBB.sma, _ = NewSMAWithoutStorage(timePeriod, selectData, func(dataItem float64, streamBarIndex int) {
		newBB.currentSMA = dataItem
	})

	newBB.stdDev, _ = NewStdDeviation(timePeriod, selectData)
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

		newBB.valueAvailableAction(upperBand, newBB.currentSMA, lowerBand, streamBarIndex)
	}

	return &newBB, nil
}

func NewBollingerBandsForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *BollingerBands, err error) {
	bb, err := NewBollingerBands(timePeriod, selectData)
	priceStream.AddTickSubscription(bb)
	return bb, err
}

func (bb *BollingerBands) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData float64 = bb.selectData(tickData)
	bb.RecieveTick(selectedData, streamBarIndex)
}

// http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance - Knuth
func (bb *BollingerBands) RecieveTick(tickData float64, streamBarIndex int) {
	bb.sma.ReceiveTick(tickData, streamBarIndex)
	bb.stdDev.ReceiveTick(tickData, streamBarIndex)
}
