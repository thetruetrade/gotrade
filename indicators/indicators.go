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
	MinimumLookbackPeriod int = 0
	MaximumLookbackPeriod int = 200
)

type Indicator interface {
	ValidFromBar() int
	GetLookbackPeriod() int
	Length() int
}

type IndicatorWithTimePeriod interface {
	GetTimePeriod() int
}

type IndicatorWithFloatBounds interface {
	MinValue() float64
	MaxValue() float64
}

type IndicatorWithIntBounds interface {
	MinValue() int64
	MaxValue() int64
}

type baseFloatBounds struct {
	minValue float64
	maxValue float64
}

func newBaseFloatBounds() *baseFloatBounds {
	ind := baseFloatBounds{minValue: math.MaxFloat64, maxValue: math.SmallestNonzeroFloat64}
	return &ind
}

func (ind *baseFloatBounds) MinValue() float64 {
	return ind.minValue
}

func (ind *baseFloatBounds) MaxValue() float64 {
	return ind.maxValue
}

type baseIntBounds struct {
	minValue int64
	maxValue int64
}

func newBaseIntBounds() *baseIntBounds {
	ind := baseIntBounds{minValue: math.MaxInt64, maxValue: math.MinInt64}
	return &ind
}

func (ind *baseIntBounds) MinValue() int64 {
	return ind.minValue
}

func (ind *baseIntBounds) MaxValue() int64 {
	return ind.maxValue
}

type baseIndicator struct {
	validFromBar   int
	dataLength     int
	selectData     gotrade.DataSelectionFunc
	lookbackPeriod int
}

func newBaseIndicator(lookbackPeriod int) *baseIndicator {
	ind := baseIndicator{lookbackPeriod: lookbackPeriod, validFromBar: -1}
	return &ind
}

func (ind *baseIndicator) ValidFromBar() int {
	return ind.validFromBar
}

func (ind *baseIndicator) GetLookbackPeriod() int {
	return ind.lookbackPeriod
}

func (ind *baseIndicator) Length() int {
	return ind.dataLength
}

type baseIndicatorWithTimePeriod struct {
	timePeriod int
}

type baseIndicatorWithFloatBounds struct {
	*baseIndicator
	*baseFloatBounds
}

func newBaseIndicatorWithFloatBounds(lookbackPeriod int) *baseIndicatorWithFloatBounds {
	ind := baseIndicatorWithFloatBounds{
		baseIndicator:   newBaseIndicator(lookbackPeriod),
		baseFloatBounds: newBaseFloatBounds()}
	return &ind
}

type baseIndicatorWithIntBounds struct {
	*baseIndicator
	*baseIntBounds
}

func newBaseIndicatorWithIntBounds(lookbackPeriod int) *baseIndicatorWithIntBounds {
	ind := baseIndicatorWithIntBounds{
		baseIndicator: newBaseIndicator(lookbackPeriod),
		baseIntBounds: newBaseIntBounds()}
	return &ind
}

func newBaseIndicatorWithTimePeriod(timePeriod int) *baseIndicatorWithTimePeriod {
	ind := baseIndicatorWithTimePeriod{timePeriod: timePeriod}
	return &ind
}

func (ind *baseIndicatorWithTimePeriod) GetTimePeriod() int {
	return ind.timePeriod
}

type ValueAvailableAction func(dataItem float64, streamBarIndex int)
type ValueAvailableActionInt func(dataItem int64, streamBarIndex int)
type ValueAvailableActionDOHLCV func(dataItem gotrade.DOHLCV, streamBarIndex int)
type ValueAvailableActionBollinger func(dataItemUpperBand float64, dataItemMiddleBand float64, dataItemLowerBand float64, streamBarIndex int)
type ValueAvailableActionMACD func(dataItemMACD float64, dataItemSignal float64, dataItemHistogram float64, streamBarIndex int)
type ValueAvailableActionAroon func(dataItemAroonUp float64, dataItemAroonDown float64, streamBarIndex int)
type ValueAvailableActionLinearReg func(dataItem float64, slope float64, intercept float64, streamBarIndex int)
