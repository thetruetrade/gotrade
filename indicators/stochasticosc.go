package indicators

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

type StochasticOscWithoutStorage struct {
	*baseIndicatorWithFloatBounds

	// private variables
	valueAvailableAction ValueAvailableActionStoch
	periodCounter        int
	slowKMA              *SmaWithoutStorage
	slowDMA              *SmaWithoutStorage
	hhv                  *HHVWithoutStorage
	llv                  *LLVWithoutStorage
	currentPeriodHigh    float64
	currentPeriodLow     float64
	currentFastK         float64
	currentSlowKMA       float64
	currentSlowDMA       float64
}

func NewStochasticOscWithoutStorage(fastKTimePeriod int, slowKTimePeriod int, slowDTimePeriod int, valueAvailableAction ValueAvailableActionStoch) (indicator *StochasticOscWithoutStorage, err error) {

	ind := StochasticOscWithoutStorage{currentSlowKMA: 0.0,
		currentSlowDMA: 0.0,
		periodCounter:  (fastKTimePeriod * -1)}

	tmpSlowKMA, err := NewSmaWithoutStorage(slowKTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentSlowKMA = dataItem
		ind.slowDMA.ReceiveTick(ind.currentSlowKMA, streamBarIndex)
	})

	tmpSlowDMA, err := NewSmaWithoutStorage(slowDTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentSlowDMA = dataItem
		ind.dataLength += 1
		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}

		var max = math.Max(ind.currentSlowKMA, ind.currentSlowDMA)
		var min = math.Min(ind.currentSlowKMA, ind.currentSlowDMA)

		if max > ind.maxValue {
			ind.maxValue = max
		}

		if min < ind.minValue {
			ind.minValue = min
		}

		ind.valueAvailableAction(ind.currentSlowKMA, ind.currentSlowDMA, streamBarIndex)
	})

	timePeriod := fastKTimePeriod - 1 + tmpSlowDMA.GetLookbackPeriod() + tmpSlowKMA.GetLookbackPeriod()

	ind.baseIndicatorWithFloatBounds = newBaseIndicatorWithFloatBounds(timePeriod)
	ind.slowKMA = tmpSlowKMA
	ind.slowDMA = tmpSlowDMA
	ind.hhv, err = NewHHVWithoutStorage(fastKTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentPeriodHigh = dataItem
	})
	ind.llv, err = NewLLVWithoutStorage(fastKTimePeriod, func(dataItem float64, streamBarIndex int) {
		ind.currentPeriodLow = dataItem
	})

	ind.valueAvailableAction = valueAvailableAction

	return &ind, err
}

// A Relative Strength Indicator
type StochasticOsc struct {
	*StochasticOscWithoutStorage

	// public variables
	SlowK []float64
	SlowD []float64
}

// NewStochasticOsc returns a new Relative Strength Indicator(StochasticOsc) configured with the
// specified timePeriod. The StochasticOsc results are stored in the DATA field.
func NewStochasticOsc(fastKTimePeriod int, slowKTimePeriod int, slowDTimePeriod int) (indicator *StochasticOsc, err error) {
	newStochasticOsc := StochasticOsc{}
	newStochasticOsc.StochasticOscWithoutStorage, err = NewStochasticOscWithoutStorage(fastKTimePeriod, slowKTimePeriod, slowDTimePeriod,
		func(dataItemK float64, dataItemD float64, streamBarIndex int) {
			newStochasticOsc.SlowK = append(newStochasticOsc.SlowK, dataItemK)
			newStochasticOsc.SlowD = append(newStochasticOsc.SlowD, dataItemD)
		})

	return &newStochasticOsc, err
}

func NewStochasticOscForStream(priceStream *gotrade.DOHLCVStream, fastKTimePeriod int, slowKTimePeriod int, slowDTimePeriod int) (indicator *StochasticOsc, err error) {
	newStochasticOsc, err := NewStochasticOsc(fastKTimePeriod, slowKTimePeriod, slowDTimePeriod)
	priceStream.AddTickSubscription(newStochasticOsc)
	return newStochasticOsc, err
}

func (ind *StochasticOscWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1
	ind.hhv.ReceiveTick(tickData.H(), streamBarIndex)
	ind.llv.ReceiveTick(tickData.L(), streamBarIndex)

	if ind.periodCounter >= 0 {
		ind.currentFastK = 100.0 * ((tickData.C() - ind.currentPeriodLow) / (ind.currentPeriodHigh - ind.currentPeriodLow))
		ind.slowKMA.ReceiveTick(ind.currentFastK, streamBarIndex)
	}
}
