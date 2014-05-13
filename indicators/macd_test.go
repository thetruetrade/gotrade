package indicators_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating a moving average convergence divergence (macd) with DOHLCV source data", func() {
	var (
		shortPeriod                 int = 3
		longPeriod                  int = 6
		signalPeriod                int = 2
		indicator                   *MACD
		indicatorWithLookbackInputs IndicatorWithLookbackSharedSpecInputs
	)

	BeforeEach(func() {
		indicator, _ = NewMACD(shortPeriod, longPeriod, signalPeriod, gotrade.UseClosePrice)
		indicatorWithLookbackInputs = IndicatorWithLookbackSharedSpecInputs{IndicatorUnderTest: indicator,
			SourceDataLength: len(sourceDOHLCVData),
			GetMaximum: func() float64 {
				return GetDataMaxMACD(indicator.Data)
			},
			GetMinimum: func() float64 {
				return GetDataMinMACD(indicator.Data)
			}}
	})

	Context("and the indicator has not yet received any ticks", func() {
		ShouldBeAnInitialisedIndicatorWithLookback(&indicatorWithLookbackInputs)
	})

	Context("and the indicator has received less ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i < indicator.GetLookbackPeriod()-1; i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedFewerTicksThanItsLookbackPeriod(&indicatorWithLookbackInputs)
	})

	Context("and the indicator has received ticks equal to the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i <= indicator.GetLookbackPeriod()-1; i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedTicksEqualToItsLookbackPeriod(&indicatorWithLookbackInputs)
	})

	Context("and the indicator has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceDOHLCVData {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedMoreTicksThanItsLookbackPeriod(&indicatorWithLookbackInputs)
	})
})
