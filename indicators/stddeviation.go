package indicators

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

type baseStdDeviation struct {
	*baseIndicatorWithLookback

	// private variables
	valueAvailableAction ValueAvailableAction
	variance             *Variance
}

func newBaseStdDeviation(lookbackPeriod int) *baseStdDeviation {
	newStdDev := baseStdDeviation{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod)}
	return &newStdDev
}

// A Standard Deviation Indicator
type StdDeviation struct {
	*baseStdDeviation

	// public variables
	Data []float64
}

// NewStdDeviation returns a new Standard Deviation (STDEV) configured with the
// specified lookbackPeriod. The STDEV results are stored in the DATA field.
func NewStdDeviation(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *StdDeviation, err error) {
	newStdDev := StdDeviation{baseStdDeviation: newBaseStdDeviation(lookbackPeriod)}
	newStdDev.variance, _ = NewVariance(lookbackPeriod, selectData)

	newStdDev.variance.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newStdDev.dataLength += 1
		if newStdDev.validFromBar == -1 {
			newStdDev.validFromBar = streamBarIndex
		}

		standardDeviation := math.Sqrt(dataItem)

		if standardDeviation > newStdDev.maxValue {
			newStdDev.maxValue = standardDeviation
		}

		if standardDeviation < newStdDev.minValue {
			newStdDev.minValue = standardDeviation
		}

		newStdDev.valueAvailableAction(standardDeviation, streamBarIndex)
	}

	newStdDev.selectData = selectData
	newStdDev.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newStdDev.Data = append(newStdDev.Data, dataItem)
	}
	return &newStdDev, nil
}

func NewStdDeviationForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *WMA, err error) {
	newStdDev, err := NewWMA(lookbackPeriod, selectData)
	priceStream.AddTickSubscription(newStdDev)
	return newStdDev, err
}

func (stdDev *baseStdDeviation) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = stdDev.selectData(tickData)
	stdDev.ReceiveTick(selectedData, streamBarIndex)
}

func (stdDev *baseStdDeviation) ReceiveTick(tickData float64, streamBarIndex int) {
	stdDev.variance.ReceiveTick(tickData, streamBarIndex)
}
