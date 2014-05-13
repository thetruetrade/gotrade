package indicators_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating bollinger bands with DOHLCV source data", func() {
	var (
		period                      int = 3
		indicator                   *BollingerBands
		indicatorWithLookbackInputs IndicatorWithLookbackSharedSpecInputs
	)

	BeforeEach(func() {
		indicator, _ = NewBollingerBands(period, gotrade.UseClosePrice)
		indicatorWithLookbackInputs = IndicatorWithLookbackSharedSpecInputs{IndicatorUnderTest: indicator,
			SourceDataLength: len(sourceDOHLCVData),
			GetMaximum: func() float64 {
				return GetDataMaxBollinger(indicator.Data, UseUpperBand)
			},
			GetMinimum: func() float64 {
				return GetDataMinBollinger(indicator.Data, UseLowerBand)
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
