package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when creating a stochrsiwithoutstorage", func() {
	var (
		indicator      *indicators.StochRsiWithoutStorage
		indicatorError error
	)

	Context("and the indicator was not given a value available action", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewStochRsiWithoutStorage(8, 5, 3, nil)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
			Expect(indicatorError).To(Equal(indicators.ErrValueAvailableActionIsNil))
		})
	})

	Context("and the indicator was given a timePeriod below the minimum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewStochRsiWithoutStorage(1, 5, 3, FakeStochValueAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a timePeriod above the maximum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewStochRsiWithoutStorage(indicators.MaximumLookbackPeriod+1, 5, 3, FakeStochValueAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a fastKTimePeriod below the minimum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewStochRsiWithoutStorage(8, 0, 3, FakeStochValueAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a fastKTimePeriod above the maximum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewStochRsiWithoutStorage(8, indicators.MaximumLookbackPeriod+1, 3, FakeStochValueAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a fastDTimePeriod below the minimum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewStochRsiWithoutStorage(8, 5, 0, FakeStochValueAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a fastDTimePeriod above the maximum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewStochRsiWithoutStorage(8, 5, indicators.MaximumLookbackPeriod+1, FakeStochValueAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})
})

var _ = Describe("when calculating a stochastic relative strength (stochrsi) with DOHLCV source data", func() {
	var (
		indicator *indicators.StochRsi
		inputs    IndicatorWithFloatBoundsSharedSpecInputs
		stream    *fakeDOHLCVStreamSubscriber
	)

	Context("given the indicator is created via the standard constructor", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewStochRsi(8, 3, 2)

			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxStoch(indicator.SlowK, indicator.SlowD)
				},
				func() float64 {
					return GetDataMinStoch(indicator.SlowK, indicator.SlowD)
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

	Context("given the indicator is created via the constructor with defaulted parameters", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewDefaultStochRsi()
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxStoch(indicator.SlowK, indicator.SlowD)
				},
				func() float64 {
					return GetDataMinStoch(indicator.SlowK, indicator.SlowD)
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

	Context("given the indicator is created via the constructor with fixed source length", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewStochRsiWithSrcLen(uint(len(sourceDOHLCVData)), 8, 5, 3)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxStoch(indicator.SlowK, indicator.SlowD)
				},
				func() float64 {
					return GetDataMinStoch(indicator.SlowK, indicator.SlowD)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.SlowK)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.SlowD)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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

			It("no new storage capcity should have been allocated", func() {
				Expect(len(indicator.SlowK)).To(Equal(cap(indicator.SlowK)))
				Expect(len(indicator.SlowD)).To(Equal(cap(indicator.SlowD)))
			})
		})
	})

	Context("given the indicator is created via the constructor with defaulted parameters and fixed source length", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewDefaultStochRsiWithSrcLen(uint(len(sourceDOHLCVData)))
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxStoch(indicator.SlowK, indicator.SlowD)
				},
				func() float64 {
					return GetDataMinStoch(indicator.SlowK, indicator.SlowD)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.SlowK)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.SlowD)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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

			It("no new storage capcity should have been allocated", func() {
				Expect(len(indicator.SlowK)).To(Equal(cap(indicator.SlowK)))
				Expect(len(indicator.SlowD)).To(Equal(cap(indicator.SlowD)))
			})
		})
	})

	Context("given the indicator is created via the constructor for use with a price stream", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewStochRsiForStream(stream, 8, 5, 3)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxStoch(indicator.SlowK, indicator.SlowD)
				},
				func() float64 {
					return GetDataMinStoch(indicator.SlowK, indicator.SlowD)
				})
		})

		It("should have requested to be attached to the stream", func() {
			Expect(stream.lastCallToAddTickSubscriptionArg).To(Equal(indicator))
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

	Context("given the indicator is created via the constructor for use with a price stream with defaulted parameters", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewDefaultStochRsiForStream(stream)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxStoch(indicator.SlowK, indicator.SlowD)
				},
				func() float64 {
					return GetDataMinStoch(indicator.SlowK, indicator.SlowD)
				})
		})

		It("should have requested to be attached to the stream", func() {
			Expect(stream.lastCallToAddTickSubscriptionArg).To(Equal(indicator))
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

	Context("given the indicator is created via the constructor for use with a price stream with fixed source length", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewStochRsiForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream, 8, 5, 3)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxStoch(indicator.SlowK, indicator.SlowD)
				},
				func() float64 {
					return GetDataMinStoch(indicator.SlowK, indicator.SlowD)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.SlowK)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.SlowD)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
		})

		It("should have requested to be attached to the stream", func() {
			Expect(stream.lastCallToAddTickSubscriptionArg).To(Equal(indicator))
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

			It("no new storage capcity should have been allocated", func() {
				Expect(len(indicator.SlowK)).To(Equal(cap(indicator.SlowK)))
				Expect(len(indicator.SlowD)).To(Equal(cap(indicator.SlowD)))
			})
		})
	})

	Context("given the indicator is created via the constructor for use with a price stream with fixed source length with defaulted parmeters", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewDefaultStochRsiForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxStoch(indicator.SlowK, indicator.SlowD)
				},
				func() float64 {
					return GetDataMinStoch(indicator.SlowK, indicator.SlowD)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.SlowK)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.SlowD)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
		})

		It("should have requested to be attached to the stream", func() {
			Expect(stream.lastCallToAddTickSubscriptionArg).To(Equal(indicator))
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

			It("no new storage capcity should have been allocated", func() {
				Expect(len(indicator.SlowK)).To(Equal(cap(indicator.SlowK)))
				Expect(len(indicator.SlowD)).To(Equal(cap(indicator.SlowD)))
			})
		})
	})
})
