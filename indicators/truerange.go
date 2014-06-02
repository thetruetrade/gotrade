package indicators

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

// TrueHigh = Max(High[0], Close[-1])
// TrueLow = Min(Low[0], Close[-1])
// TrueRange = TrueHigh = TrueLow

type TrueRangeWithoutStorage struct {
	*baseIndicator

	// private variables
	periodCounter        int
	previousClose        float64
	valueAvailableAction ValueAvailableAction
}

func NewTrueRangeWithoutStorage(valueAvailableAction ValueAvailableAction) (indicator *TrueRangeWithoutStorage, err error) {
	ind := TrueRangeWithoutStorage{baseIndicator: newBaseIndicator(1),
		periodCounter: -1,
		previousClose: 0.0}
	ind.valueAvailableAction = valueAvailableAction
	return &ind, nil
}

type TrueRange struct {
	*TrueRangeWithoutStorage
	Data []float64
}

func NewTrueRange() (indicator *TrueRange, err error) {
	newTrueRange := TrueRange{}
	newTrueRange.TrueRangeWithoutStorage, err = NewTrueRangeWithoutStorage(func(dataItem float64, streamBarIndex int) {
		newTrueRange.Data = append(newTrueRange.Data, dataItem)
	})
	return &newTrueRange, err
}

func NewTrueRangeForStream(priceStream *gotrade.DOHLCVStream) (indicator *TrueRange, err error) {
	ind, err := NewTrueRange()
	priceStream.AddTickSubscription(ind)
	return ind, err
}

func (ind *TrueRangeWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1

	if ind.periodCounter > 0 {
		ind.dataLength += 1

		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}
		high := math.Max(tickData.H(), ind.previousClose)
		low := math.Min(tickData.L(), ind.previousClose)
		trueRange := high - low

		if trueRange > ind.maxValue {
			ind.maxValue = trueRange
		}

		if trueRange < ind.minValue {
			ind.minValue = trueRange
		}

		ind.valueAvailableAction(trueRange, streamBarIndex)
	}
	ind.previousClose = tickData.C()

}
