package indicators_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating a linear regression intercept (linearregintercept) with DOHLCV source data", func() {
	var (
		period    int = 3
		indicator *LinearRegIntercept
		inputs    IndicatorWithFloatBoundsSharedSpecInputs
	)

	BeforeEach(func() {
		indicator, _ = NewLinearRegIntercept(period, gotrade.UseClosePrice)

		inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
			func() float64 {
				return GetFloatDataMax(indicator.Data)
			},
			func() float64 {
				return GetFloatDataMin(indicator.Data)
			})
	})

	Context("and the indicator has not yet received any ticks", func() {
		ShouldBeAnInitialisedIndicator(&inputs)

		ShouldNotHaveAnyFloatBoundsSetYet(&inputs)
	})

	Context("and the indicator has received less ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i < indicator.GetLookbackPeriod(); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedFewerTicksThanItsLookbackPeriod(&inputs)

		ShouldNotHaveAnyFloatBoundsSetYet(&inputs)
	})

	Context("and the indicator has received ticks equal to the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i <= indicator.GetLookbackPeriod(); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedTicksEqualToItsLookbackPeriod(&inputs)

		ShouldHaveFloatBoundsSetToMinMaxOfResults(&inputs)
	})

	Context("and the indicator has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceDOHLCVData {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedMoreTicksThanItsLookbackPeriod(&inputs)

		ShouldHaveFloatBoundsSetToMinMaxOfResults(&inputs)
	})

	Context("and the indicator has recieved all of its ticks", func() {
		BeforeEach(func() {
			for i := 0; i < len(sourceDOHLCVData); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedAllOfItsTicks(&inputs)

		ShouldHaveFloatBoundsSetToMinMaxOfResults(&inputs)
	})
})
