package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	"github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when creating a macd", func() {
	var (
		fastTimePeriod   int = 3
		slowTimePeriod   int = 6
		signalTimePeriod int = 2
		indicator        *indicators.Macd
		indicatorError   error
	)

	Context("and the indicator was given a fastTimePeriod below the minimum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewMacd(1, slowTimePeriod, signalTimePeriod, gotrade.UseClosePrice)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a fastTimePeriod above the maximum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewMacd(indicators.MaximumLookbackPeriod+1, slowTimePeriod, signalTimePeriod, gotrade.UseClosePrice)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a slowTimePeriod below the minimum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewMacd(fastTimePeriod, 1, signalTimePeriod, gotrade.UseClosePrice)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a slowTimePeriod above the maximum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewMacd(fastTimePeriod, indicators.MaximumLookbackPeriod+1, signalTimePeriod, gotrade.UseClosePrice)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a signalTimePeriod below the minimum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewMacd(fastTimePeriod, slowTimePeriod, 0, gotrade.UseClosePrice)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a signalTimePeriod above the maximum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewMacd(fastTimePeriod, slowTimePeriod, indicators.MaximumLookbackPeriod+1, gotrade.UseClosePrice)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a nil data selection func", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewMacd(fastTimePeriod, slowTimePeriod, signalTimePeriod, nil)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
			Expect(indicatorError).To(Equal(indicators.ErrDOHLCVDataSelectFuncIsNil))
		})
	})
})

var _ = Describe("when calculating a moving average convergence divergence (macd) with DOHLCV source data", func() {
	var (
		fastTimePeriod   int = 3
		slowTimePeriod   int = 6
		signalTimePeriod int = 2
		indicator        *indicators.Macd
		inputs           IndicatorWithFloatBoundsSharedSpecInputs
		stream           *fakeDOHLCVStreamSubscriber
	)

	Context("given the indicator is created via the standard constructor", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewMacd(fastTimePeriod, slowTimePeriod, signalTimePeriod, gotrade.UseClosePrice)

			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				},
				func() float64 {
					return GetDataMinMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
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
			indicator, _ = indicators.NewDefaultMacd()
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				},
				func() float64 {
					return GetDataMinMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
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
			indicator, _ = indicators.NewMacdWithSrcLen(uint(len(sourceDOHLCVData)), fastTimePeriod, slowTimePeriod, signalTimePeriod, gotrade.UseClosePrice)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				},
				func() float64 {
					return GetDataMinMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Macd)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.Signal)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.Histogram)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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
				Expect(len(indicator.Macd)).To(Equal(cap(indicator.Macd)))
				Expect(len(indicator.Signal)).To(Equal(cap(indicator.Signal)))
				Expect(len(indicator.Histogram)).To(Equal(cap(indicator.Histogram)))
			})
		})
	})

	Context("given the indicator is created via the constructor with defaulted parameters and fixed source length", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewDefaultMacdWithSrcLen(uint(len(sourceDOHLCVData)))
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				},
				func() float64 {
					return GetDataMinMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Macd)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.Signal)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.Histogram)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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
				Expect(len(indicator.Macd)).To(Equal(cap(indicator.Macd)))
				Expect(len(indicator.Signal)).To(Equal(cap(indicator.Signal)))
				Expect(len(indicator.Histogram)).To(Equal(cap(indicator.Histogram)))
			})
		})
	})

	Context("given the indicator is created via the constructor for use with a price stream", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewMacdForStream(stream, fastTimePeriod, slowTimePeriod, signalTimePeriod, gotrade.UseClosePrice)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				},
				func() float64 {
					return GetDataMinMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
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
			indicator, _ = indicators.NewDefaultMacdForStream(stream)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				},
				func() float64 {
					return GetDataMinMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
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
			indicator, _ = indicators.NewMacdForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream, fastTimePeriod, slowTimePeriod, signalTimePeriod, gotrade.UseClosePrice)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				},
				func() float64 {
					return GetDataMinMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Macd)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.Signal)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.Histogram)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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
				Expect(len(indicator.Macd)).To(Equal(cap(indicator.Macd)))
				Expect(len(indicator.Signal)).To(Equal(cap(indicator.Signal)))
				Expect(len(indicator.Histogram)).To(Equal(cap(indicator.Histogram)))
			})
		})
	})

	Context("given the indicator is created via the constructor for use with a price stream with fixed source length with defaulted parmeters", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewDefaultMacdForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetDataMaxMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				},
				func() float64 {
					return GetDataMinMacd(indicator.Macd, indicator.Signal, indicator.Histogram)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Macd)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.Signal)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
			Expect(cap(indicator.Histogram)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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
				Expect(len(indicator.Macd)).To(Equal(cap(indicator.Macd)))
				Expect(len(indicator.Signal)).To(Equal(cap(indicator.Signal)))
				Expect(len(indicator.Histogram)).To(Equal(cap(indicator.Histogram)))
			})
		})
	})

})
