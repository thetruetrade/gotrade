package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	"github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when creating an llvbarswithoutstorage", func() {
	var (
		indicator      *indicators.LlvBarsWithoutStorage
		indicatorError error
	)

	Context("and the indicator was not given a value available action", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewLlvBarsWithoutStorage(4, nil)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
			Expect(indicatorError).To(Equal(indicators.ErrValueAvailableActionIsNil))
		})
	})

	Context("and the indicator was given a timePeriod below the minimum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewLlvBarsWithoutStorage(0, fakeIntValAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a timePeriod above the maximum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewLlvBarsWithoutStorage(indicators.MaximumLookbackPeriod+1, fakeIntValAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})
})

var _ = Describe("when calculating a lowest low value bars (llvbars) with DOHLCV source data", func() {
	var (
		period         int = 3
		indicator      *indicators.LlvBars
		inputs         IndicatorWithIntBoundsSharedSpecInputs
		stream         *fakeDOHLCVStreamSubscriber
		indicatorError error
	)

	Context("given the indicator is created via the standard constructor", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewLlvBars(period, gotrade.UseClosePrice)

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

	Context("given the indicator is created via the standard constructor with a nil data selection func", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewLlvBars(period, nil)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
			Expect(indicatorError).To(Equal(indicators.ErrDOHLCVDataSelectFuncIsNil))
		})
	})

	Context("given the indicator is created via the constructor with defaulted parameters", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewDefaultLlvBars()
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

	Context("given the indicator is created via the constructor with fixed source length", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewLlvBarsWithSrcLen(uint(len(sourceDOHLCVData)), 4, gotrade.UseClosePrice)
			inputs = NewIndicatorWithIntBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() int64 {
					return GetIntDataMax(indicator.Data)
				},
				func() int64 {
					return GetIntDataMin(indicator.Data)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Data)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
		})

		Context("and the indicator has not yet received any ticks", func() {
			ShouldBeAnInitialisedIndicator(&inputs)

			ShouldNotHaveAnyIntBoundsSetYet(&inputs)
		})

		Context("and the indicator has recieved all of its ticks", func() {
			BeforeEach(func() {
				for i := 0; i < len(sourceDOHLCVData); i++ {
					indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
				}
			})

			ShouldBeAnIndicatorThatHasReceivedAllOfItsTicks(&inputs)

			ShouldHaveIntBoundsSetToMinMaxOfResults(&inputs)

			It("no new storage capcity should have been allocated", func() {
				Expect(len(indicator.Data)).To(Equal(cap(indicator.Data)))
			})
		})
	})

	Context("given the indicator is created via the constructor with defaulted parameters and fixed source length", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewDefaultLlvBarsWithSrcLen(uint(len(sourceDOHLCVData)))
			inputs = NewIndicatorWithIntBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() int64 {
					return GetIntDataMax(indicator.Data)
				},
				func() int64 {
					return GetIntDataMin(indicator.Data)
				})
		})

		It("should have pre-allocated storge for the output data", func() {
			Expect(cap(indicator.Data)).To(Equal(len(sourceDOHLCVData) - indicator.GetLookbackPeriod()))
		})

		Context("and the indicator has not yet received any ticks", func() {
			ShouldBeAnInitialisedIndicator(&inputs)

			ShouldNotHaveAnyIntBoundsSetYet(&inputs)
		})

		Context("and the indicator has recieved all of its ticks", func() {
			BeforeEach(func() {
				for i := 0; i < len(sourceDOHLCVData); i++ {
					indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
				}
			})

			ShouldBeAnIndicatorThatHasReceivedAllOfItsTicks(&inputs)

			ShouldHaveIntBoundsSetToMinMaxOfResults(&inputs)

			It("no new storage capcity should have been allocated", func() {
				Expect(len(indicator.Data)).To(Equal(cap(indicator.Data)))
			})
		})
	})

	Context("given the indicator is created via the constructor for use with a price stream", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewLlvBarsForStream(stream, 4, gotrade.UseClosePrice)
			inputs = NewIndicatorWithIntBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() int64 {
					return GetIntDataMax(indicator.Data)
				},
				func() int64 {
					return GetIntDataMin(indicator.Data)
				})
		})

		It("should have requested to be attached to the stream", func() {
			Expect(stream.lastCallToAddTickSubscriptionArg).To(Equal(indicator))
		})

		Context("and the indicator has not yet received any ticks", func() {
			ShouldBeAnInitialisedIndicator(&inputs)

			ShouldNotHaveAnyIntBoundsSetYet(&inputs)
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

	Context("given the indicator is created via the constructor for use with a price stream with defaulted parameters", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewDefaultLlvBarsForStream(stream)
			inputs = NewIndicatorWithIntBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() int64 {
					return GetIntDataMax(indicator.Data)
				},
				func() int64 {
					return GetIntDataMin(indicator.Data)
				})
		})

		It("should have requested to be attached to the stream", func() {
			Expect(stream.lastCallToAddTickSubscriptionArg).To(Equal(indicator))
		})

		Context("and the indicator has not yet received any ticks", func() {
			ShouldBeAnInitialisedIndicator(&inputs)

			ShouldNotHaveAnyIntBoundsSetYet(&inputs)
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

	Context("given the indicator is created via the constructor for use with a price stream with fixed source length", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewLlvBarsForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream, 4, gotrade.UseClosePrice)
			inputs = NewIndicatorWithIntBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() int64 {
					return GetIntDataMax(indicator.Data)
				},
				func() int64 {
					return GetIntDataMin(indicator.Data)
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

			ShouldNotHaveAnyIntBoundsSetYet(&inputs)
		})

		Context("and the indicator has recieved all of its ticks", func() {
			BeforeEach(func() {
				for i := 0; i < len(sourceDOHLCVData); i++ {
					indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
				}
			})

			ShouldBeAnIndicatorThatHasReceivedAllOfItsTicks(&inputs)

			ShouldHaveIntBoundsSetToMinMaxOfResults(&inputs)

			It("no new storage capcity should have been allocated", func() {
				Expect(len(indicator.Data)).To(Equal(cap(indicator.Data)))
			})
		})
	})

	Context("given the indicator is created via the constructor for use with a price stream with fixed source length with defaulted parmeters", func() {
		BeforeEach(func() {
			stream = newFakeDOHLCVStreamSubscriber()
			indicator, _ = indicators.NewDefaultLlvBarsForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream)
			inputs = NewIndicatorWithIntBoundsSharedSpecInputs(indicator, len(sourceDOHLCVData), indicator,
				func() int64 {
					return GetIntDataMax(indicator.Data)
				},
				func() int64 {
					return GetIntDataMin(indicator.Data)
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

			ShouldNotHaveAnyIntBoundsSetYet(&inputs)
		})

		Context("and the indicator has recieved all of its ticks", func() {
			BeforeEach(func() {
				for i := 0; i < len(sourceDOHLCVData); i++ {
					indicator.ReceiveDOHLCVTick(sourceDOHLCVData[i], i+1)
				}
			})

			ShouldBeAnIndicatorThatHasReceivedAllOfItsTicks(&inputs)

			ShouldHaveIntBoundsSetToMinMaxOfResults(&inputs)

			It("no new storage capcity should have been allocated", func() {
				Expect(len(indicator.Data)).To(Equal(cap(indicator.Data)))
			})
		})
	})
})
