package indicators

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

type StdDeviationWithoutStorage struct {
	*baseIndicatorWithLookback

	// private variables
	valueAvailableAction ValueAvailableAction
	variance             *Variance
}

func NewStdDeviationWithoutStorage(lookbackPeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableAction) (indicator *StdDeviationWithoutStorage, err error) {
	newStdDev := StdDeviationWithoutStorage{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod)}

	newStdDev.selectData = selectData
	newStdDev.valueAvailableAction = valueAvailableAction

	newStdDev.variance, err = NewVariance(lookbackPeriod, selectData)

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

	return &newStdDev, err
}

// A Standard Deviation Indicator
type StdDeviation struct {
	*StdDeviationWithoutStorage

	// public variables
	Data []float64
}

// NewStdDeviation returns a new Standard Deviation (STDEV) configured with the
// specified lookbackPeriod. The STDEV results are stored in the DATA field.
func NewStdDeviation(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *StdDeviation, err error) {
	newStdDev := StdDeviation{}
	newStdDev.StdDeviationWithoutStorage, err = NewStdDeviationWithoutStorage(lookbackPeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newStdDev.Data = append(newStdDev.Data, dataItem)
		})

	newStdDev.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newStdDev.Data = append(newStdDev.Data, dataItem)
	}
	return &newStdDev, err
}

func NewStdDeviationForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *StdDeviation, err error) {
	newStdDev, err := NewStdDeviation(lookbackPeriod, selectData)
	priceStream.AddTickSubscription(newStdDev)
	return newStdDev, err
}

func (stdDev *StdDeviationWithoutStorage) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = stdDev.selectData(tickData)
	stdDev.ReceiveTick(selectedData, streamBarIndex)
}

func (stdDev *StdDeviationWithoutStorage) ReceiveTick(tickData float64, streamBarIndex int) {
	stdDev.variance.ReceiveTick(tickData, streamBarIndex)
}
