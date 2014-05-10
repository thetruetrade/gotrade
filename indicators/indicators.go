package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

var (
	// Indicator errors
	ErrSourceDataEmpty                      = errors.New("Source data is empty")
	ErrNotEnoughSourceDataForLookbackPeriod = errors.New("Source data does not contain enough data for the specfied lookback period")
	ErrLookbackPeriodMustBeGreaterThanZero  = errors.New("Lookback period must be greater than 0")
)

type Indicator interface {
	ValidFromBar() int
	Length() int
	MinValue() float64
	MaxValue() float64
}

type baseIndicator struct {
	validFromBar int
	dataLength   int
	selectData   gotrade.DataSelectionFunc
	minValue     float64
	maxValue     float64
}

func newBaseIndicator() *baseIndicator {
	ind := baseIndicator{validFromBar: -1, minValue: math.MaxFloat64, maxValue: math.SmallestNonzeroFloat64}
	return &ind
}

func (ind *baseIndicator) ValidFromBar() int {
	return ind.validFromBar
}

func (ind *baseIndicator) MinValue() float64 {
	return ind.minValue
}

func (ind *baseIndicator) MaxValue() float64 {
	return ind.maxValue
}

func (ind *baseIndicator) Length() int {
	return ind.dataLength
}

type baseIndicatorWithLookback struct {
	*baseIndicator
	LookbackPeriod int
}

func newBaseIndicatorWithLookback(lookbackPeriod int) *baseIndicatorWithLookback {
	ind := baseIndicatorWithLookback{baseIndicator: newBaseIndicator(),
		LookbackPeriod: lookbackPeriod}
	return &ind
}

type ValueAvailableAction func(dataItem float64, streamBarIndex int)
type ValueAvailableActionDOHLCV func(dataItem gotrade.DOHLCV, streamBarIndex int)
type ValueAvailableActionBollinger func(dataItem BollingerBand, streamBarIndex int)
type ValueAvailableActionMACD func(dataItem MACDData, streamBarIndex int)
