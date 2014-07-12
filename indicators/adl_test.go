package indicators_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating an accumulation distribution line (adl) with DOHLCV source data", func() {
	var (
		indicator *indicators.ADL
		inputs    IndicatorWithFloatBoundsSharedSpecInputs
	)

	BeforeEach(func() {
		indicator, _ = indicators.NewADL()
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
