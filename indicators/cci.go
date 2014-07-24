package indicators

import (
	"container/list"
	"github.com/thetruetrade/gotrade"
	"math"
)

type CCIWithoutStorage struct {
	*baseIndicator
	*baseFloatBounds

	// private variables
	valueAvailableAction   ValueAvailableActionFloat
	periodCounter          int
	typicalPriceAvg        *SmaWithoutStorage
	factor                 float64
	typicalPriceHistory    *list.List
	currentAvgTypicalPrice float64
	currentTypicalPrice    float64
	timePeriod             int
}

func NewCCIWithoutStorage(timePeriod int, valueAvailableAction ValueAvailableActionFloat) (indicator *CCIWithoutStorage, err error) {
	lookbackPeriod := timePeriod - 1

	ind := CCIWithoutStorage{
		baseIndicator:        newBaseIndicator(lookbackPeriod),
		baseFloatBounds:      newBaseFloatBounds(),
		factor:               0.015,
		periodCounter:        (timePeriod * -1),
		valueAvailableAction: valueAvailableAction,
		typicalPriceHistory:  list.New(),
		timePeriod:           timePeriod,
	}

	ind.typicalPriceAvg, err = NewSmaWithoutStorage(timePeriod, func(dataItem float64, streamBarIndex int) {
		currentTypicalPriceAvg := dataItem

		var meanDeviation float64 = 0.0
		// calculate the mean deviation
		for e := ind.typicalPriceHistory.Front(); e != nil; e = e.Next() {
			value := e.Value.(float64)
			meanDeviation += math.Abs(value - currentTypicalPriceAvg)
		}
		meanDeviation /= float64(ind.timePeriod)

		result := ((ind.currentTypicalPrice - currentTypicalPriceAvg) / (ind.factor * meanDeviation))

		ind.dataLength += 1

		if ind.validFromBar == -1 {
			ind.validFromBar = streamBarIndex
		}

		if result > ind.maxValue {
			ind.maxValue = result
		}

		if result < ind.minValue {
			ind.minValue = result
		}

		ind.valueAvailableAction(result, streamBarIndex)

	})

	return &ind, err
}

// A Relative Strength Indicator
type CCI struct {
	*CCIWithoutStorage

	// public variables
	Data []float64
}

// NewCCI returns a new Relative Strength Indicator(CCI) configured with the
// specified timePeriod. The CCI results are stored in the DATA field.
func NewCCI(timePeriod int) (indicator *CCI, err error) {
	newCCI := CCI{}
	newCCI.CCIWithoutStorage, err = NewCCIWithoutStorage(timePeriod,
		func(dataItem float64, streamBarIndex int) {
			newCCI.Data = append(newCCI.Data, dataItem)
		})

	return &newCCI, err
}

func NewCCIForStream(priceStream *gotrade.DOHLCVStream, timePeriod int) (indicator *CCI, err error) {
	newCCI, err := NewCCI(timePeriod)
	priceStream.AddTickSubscription(newCCI)
	return newCCI, err
}

func (ind *CCIWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	ind.periodCounter += 1

	// calculate the typical price
	typicalPrice := (tickData.H() + tickData.L() + tickData.C()) / 3.0
	ind.currentTypicalPrice = typicalPrice

	// push it to the history
	ind.typicalPriceHistory.PushBack(typicalPrice)

	// trim the history
	if ind.typicalPriceHistory.Len() > ind.timePeriod {
		var first = ind.typicalPriceHistory.Front()
		ind.typicalPriceHistory.Remove(first)
	}

	// add it to the average
	ind.typicalPriceAvg.ReceiveTick(typicalPrice, streamBarIndex)
}
