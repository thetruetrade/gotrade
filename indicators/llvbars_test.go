package indicators_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating a lowest low value bars (llvbars) with DOHLCV source data", func() {
	var (
		period    int = 3
		indicator *LLVBars
		inputs    IndicatorWithIntBoundsSharedSpecInputs
	)

	BeforeEach(func() {
		indicator, _ = NewLLVBars(period, gotrade.UseClosePrice)

		inputs = NewIndicatorWithIntBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
			func() int64 {
				return GetIntDataMax(indicator.Data)
			},
			func() int64 {
				return GetIntDataMin(indicator.Data)
			})
	})

	Context("and the indicator has not yet received any ticks", func() {
		ShouldBeAnInitialisedIndicator(&inputs)

		ShouldNotHaveAnyIntBoundsSetYet(&inputs)
	})

	Context("and the indicator has received less ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i < indicator.GetLookbackPeriod(); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedFewerTicksThanItsLookbackPeriod(&inputs)

		ShouldNotHaveAnyIntBoundsSetYet(&inputs)
	})

	Context("and the indicator has received ticks equal to the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i <= indicator.GetLookbackPeriod(); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedTicksEqualToItsLookbackPeriod(&inputs)

		ShouldHaveIntBoundsSetToMinMaxOfResults(&inputs)
	})

	Context("and the indicator has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceDOHLCVData {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedMoreTicksThanItsLookbackPeriod(&inputs)

		ShouldHaveIntBoundsSetToMinMaxOfResults(&inputs)
	})

	Context("and the indicator has recieved all of its ticks", func() {
		BeforeEach(func() {
			for i := 0; i < len(sourceDOHLCVData); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedAllOfItsTicks(&inputs)

		ShouldHaveIntBoundsSetToMinMaxOfResults(&inputs)
	})
})
