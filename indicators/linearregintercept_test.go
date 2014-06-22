package indicators_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating a linear regression intercept (linearregintercept) with DOHLCV source data", func() {
	var (
		period          int = 3
		indicator       *LinearRegIntercept
		indicatorInputs IndicatorSharedSpecInputs
	)

	BeforeEach(func() {
		indicator, _ = NewLinearRegIntercept(period, gotrade.UseClosePrice)
		indicatorInputs = IndicatorSharedSpecInputs{IndicatorUnderTest: indicator,
			SourceDataLength: len(sourceDOHLCVData),
			GetMaximum: func() float64 {
				return GetDataMax(indicator.Data)
			},
			GetMinimum: func() float64 {
				return GetDataMin(indicator.Data)
			}}
	})

	Context("and the indicator has not yet received any ticks", func() {
		ShouldBeAnInitialisedIndicator(&indicatorInputs)
	})

	Context("and the indicator has received less ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i < indicator.GetLookbackPeriod(); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedFewerTicksThanItsLookbackPeriod(&indicatorInputs)
	})

	Context("and the indicator has received ticks equal to the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i <= indicator.GetLookbackPeriod(); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedTicksEqualToItsLookbackPeriod(&indicatorInputs)
	})

	Context("and the indicator has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceDOHLCVData {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedMoreTicksThanItsLookbackPeriod(&indicatorInputs)
	})

	Context("and the indicator has recieved all of its ticks", func() {
		BeforeEach(func() {
			for i := 0; i < len(sourceDOHLCVData); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedAllOfItsTicks(&indicatorInputs)
	})
})
