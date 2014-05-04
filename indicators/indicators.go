package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
)

var (
	// Indicator errors
	ErrSourceDataEmpty                      = errors.New("Source data is empty")
	ErrNotEnoughSourceDataForLookbackPeriod = errors.New("Source data does not contain enough data for the specfied lookback period")
	ErrLookbackPeriodMustBeGreaterThanZero  = errors.New("Lookback period must be greater than 0")
)

type Indicator struct {
	validFromBar         int
	validFromBarIndex    int
	dataLength           int
	valueAvailableAction ValueAvailableAction
	transformData        gotrade.DataTransformationFunc
	minValue             float64
	maxValue             float64
}

func (ind *Indicator) ValidFromBar() int {
	return ind.validFromBar
}

func (ind *Indicator) MinValue() float64 {
	return ind.minValue
}

func (ind *Indicator) MaxValue() float64 {
	return ind.maxValue
}

type ValueAvailableAction func(dataItem float64, streamBarIndex int)
type ValueAvailableActionDOHLCV func(dataItem gotrade.DOHLCV, streamBarIndex int)
type ValueAvailableActionBollinger func(dataItem BollingerBandEntry, streamBarIndex int)

// **************************
// Indicator helper functions
// **************************

// Ensures that the source data is not empty
func checkSourceDataIsNotEmpty(sourceData []float64) error {
	// ensure we have some data to start with
	if sourceData == nil {
		return ErrSourceDataEmpty
	}

	return nil
}

// Ensures that the source data is valid for the specified lookback period
func checkSourceValidForLookbackPeriod(sourceData []float64, lookbackPeriod int) error {
	// check that the lookbackPeriod is greater than 0
	if lookbackPeriod <= 0 {
		return ErrLookbackPeriodMustBeGreaterThanZero
	}

	// check the length of the source data is at least greater than the lookbackPeriod -1
	if len(sourceData) < (lookbackPeriod - 1) {
		return ErrNotEnoughSourceDataForLookbackPeriod
	}

	return nil
}
