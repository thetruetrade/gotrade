package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade/indicators"
	"math"
)

var _ = Describe("when creating an sarwithoutstorage", func() {
	var (
		indicator      *indicators.SarWithoutStorage
		indicatorError error
	)

	Context("and the indicator was not given a value available action", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewSarWithoutStorage(0.02, 0.2, nil)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
			Expect(indicatorError).To(Equal(indicators.ErrValueAvailableActionIsNil))
		})
	})

	Context("and the indicator was given a accelerationFactor below the minimum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewSarWithoutStorage(-1, 0.2, fakeFloatValAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a accelerationFactor above the maximum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewSarWithoutStorage(math.MaxFloat64, 0.2, fakeFloatValAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a accelerationFactorMax below the minimum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewSarWithoutStorage(0.2, -1, fakeFloatValAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})

	Context("and the indicator was given a accelerationFactorMax above the maximum", func() {
		BeforeEach(func() {
			indicator, indicatorError = indicators.NewSarWithoutStorage(0.2, math.MaxFloat64, fakeFloatValAvailable)
		})

		It("the indicator should not be created and return the appropriate error message", func() {
			Expect(indicator).To(BeNil())
		})
	})
})

var _ = Describe("when calculating the parabolic stop and reverse (Sar) with DOHLCV source data", func() {
	var (
		indicator *indicators.Sar
		inputs    IndicatorWithFloatBoundsSharedSpecInputs
		stream    *fakeDOHLCVStreamSubscriber
	)

	Context("given the indicator is created via the standard constructor", func() {
		BeforeEach(func() {
			indicator, _ = indicators.NewSar(0.02, 0.2)

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
			indicator, _ = indicators.NewDefaultSar()
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
			indicator, _ = indicators.NewSarWithSrcLen(uint(len(sourceDOHLCVData)), 0.02, 0.2)
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
			indicator, _ = indicators.NewDefaultSarWithSrcLen(uint(len(sourceDOHLCVData)))
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
			indicator, _ = indicators.NewSarForStream(stream, 0.02, 0.2)
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
			indicator, _ = indicators.NewDefaultSarForStream(stream)
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
			indicator, _ = indicators.NewSarForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream, 0.02, 0.2)
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
			indicator, _ = indicators.NewDefaultSarForStreamWithSrcLen(uint(len(sourceDOHLCVData)), stream)
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
