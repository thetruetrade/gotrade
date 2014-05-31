package indicators_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating an accumulation distribution line (adl)) with DOHLCV source data", func() {
	var (
		indicator       *indicators.ADL
		indicatorInputs IndicatorSharedSpecInputs
	)

	BeforeEach(func() {
		indicator, _ = indicators.NewADL()
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

	Context("and the indicator has recveived all of its ticks", func() {
		BeforeEach(func() {
			for i := 0; i < len(sourceDOHLCVData); i++ {
				indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
			}
		})

		ShouldBeAnIndicatorThatHasReceivedAllOfItsTicks(&indicatorInputs)
	})
})
