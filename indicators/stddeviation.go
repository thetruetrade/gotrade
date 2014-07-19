package indicators

import (
	"github.com/thetruetrade/gotrade"
	"math"
)

type StdDeviationWithoutStorage struct {
	*baseIndicatorWithFloatBounds
	*baseIndicatorWithTimePeriod

	// private variables
	valueAvailableAction ValueAvailableActionFloat
	variance             *Variance
}

func NewStdDeviationWithoutStorage(timePeriod int, selectData gotrade.DataSelectionFunc, valueAvailableAction ValueAvailableActionFloat) (indicator *StdDeviationWithoutStorage, err error) {
	newStdDev := StdDeviationWithoutStorage{baseIndicatorWithFloatBounds: newBaseIndicatorWithFloatBounds(timePeriod - 1),
		baseIndicatorWithTimePeriod: newBaseIndicatorWithTimePeriod(timePeriod)}

	newStdDev.selectData = selectData
	newStdDev.valueAvailableAction = valueAvailableAction

	newStdDev.variance, err = NewVariance(timePeriod, selectData)

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
// specified timePeriod. The STDEV results are stored in the DATA field.
func NewStdDeviation(timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *StdDeviation, err error) {
	newStdDev := StdDeviation{}
	newStdDev.StdDeviationWithoutStorage, err = NewStdDeviationWithoutStorage(timePeriod, selectData,
		func(dataItem float64, streamBarIndex int) {
			newStdDev.Data = append(newStdDev.Data, dataItem)
		})

	newStdDev.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newStdDev.Data = append(newStdDev.Data, dataItem)
	}
	return &newStdDev, err
}

func NewStdDeviationForStream(priceStream *gotrade.DOHLCVStream, timePeriod int, selectData gotrade.DataSelectionFunc) (indicator *StdDeviation, err error) {
	newStdDev, err := NewStdDeviation(timePeriod, selectData)
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
