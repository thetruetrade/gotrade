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

	// lookback minimum
	MinimumLookbackPeriod int = 1
	MaximumLookbackPeriod int = 200
)

type Indicator interface {
	ValidFromBar() int
	Length() int
	MinValue() float64
	MaxValue() float64
}

type IndicatorWithLookback interface {
	GetLookbackPeriod() int
}

type IndicatorWithTimePeriod interface {
	GetTimePeriod() int
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
	lookbackPeriod int
}

func newBaseIndicatorWithLookback(lookbackPeriod int) *baseIndicatorWithLookback {
	ind := baseIndicatorWithLookback{lookbackPeriod: lookbackPeriod}
	return &ind
}

func (ind *baseIndicatorWithLookback) GetLookbackPeriod() int {
	return ind.lookbackPeriod
}

type baseIndicatorWithTimePeriod struct {
	timePeriod int
}

func newBaseIndicatorWithTimePeriod(timePeriod int) *baseIndicatorWithTimePeriod {
	ind := baseIndicatorWithTimePeriod{timePeriod: timePeriod}
	return &ind
}

func (ind *baseIndicatorWithTimePeriod) GetTimePeriod() int {
	return ind.timePeriod
}

type ValueAvailableAction func(dataItem float64, streamBarIndex int)
type ValueAvailableActionDOHLCV func(dataItem gotrade.DOHLCV, streamBarIndex int)
type ValueAvailableActionBollinger func(dataItemUpperBand float64, dataItemMiddleBand float64, dataItemLowerBand float64, streamBarIndex int)
type ValueAvailableActionMACD func(dataItemMACD float64, dataItemSignal float64, dataItemHistogram float64, streamBarIndex int)
type ValueAvailableActionAroon func(dataItemAroonUp float64, dataItemAroonDown float64, streamBarIndex int)
