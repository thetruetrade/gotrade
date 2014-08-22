package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	"github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating a time series forecast (Tsf) with DOHLCV source data", func() {
	var (
		period    int = 3
		indicator *indicators.Tsf
		inputs    IndicatorWithFloatBoundsSharedSpecInputs
		stream    *fakeDOHLCVStreamSubscriber
	)

	Context("given the indicator is created via the standard constructor", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewTsf(period, gotrade.UseClosePrice)

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

	Context("given the indicator is created via the constructor with defaulted parameters", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewDefaultTsf()
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

	Context("given the indicator is created via the constructor with fixed source length", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewTsfWithSrcLen(uint(len(sourceDOHLCVData)), 4, gotrade.UseClosePrice)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetFloatDataMax(indicator.Data)
				},
				func() float64 {
					return GetFloatDataMin(indicator.Data)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Data)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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
				Expect(len(indicator.Data)).To(Equal(cap(indicator.Data)))
			})
		})
	})

	Context("given the indicator is created via the constructor with defaulted parameters and fixed source length", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewDefaultTsfWithSrcLen(uint(len(sourceDOHLCVData)))
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetFloatDataMax(indicator.Data)
				},
				func() float64 {
					return GetFloatDataMin(indicator.Data)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Data)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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
				Expect(len(indicator.Data)).To(Equal(cap(indicator.Data)))
			})
		})
	})

	Context("given the indicator is created via the constructor for use with a price stream", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewTsfForStream(stream, 4, gotrade.UseClosePrice)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetFloatDataMax(indicator.Data)
				},
				func() float64 {
					return GetFloatDataMin(indicator.Data)
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
			indicator, _ = indicators.NewDefaultTsfForStream(stream)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetFloatDataMax(indicator.Data)
				},
				func() float64 {
					return GetFloatDataMin(indicator.Data)
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
			indicator, _ = indicators.NewTsfForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream, 4, gotrade.UseClosePrice)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetFloatDataMax(indicator.Data)
				},
				func() float64 {
					return GetFloatDataMin(indicator.Data)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Data)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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
				Expect(len(indicator.Data)).To(Equal(cap(indicator.Data)))
			})
		})
	})

	Context("given the indicator is created via the constructor for use with a price stream with fixed source length with defaulted parmeters", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewDefaultTsfForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream)
			inputs = NewIndicatorWithFloatBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() float64 {
					return GetFloatDataMax(indicator.Data)
				},
				func() float64 {
					return GetFloatDataMin(indicator.Data)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Data)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
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
				Expect(len(indicator.Data)).To(Equal(cap(indicator.Data)))
			})
		})
	})
})
