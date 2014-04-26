package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating a EMA (exponential moving average)", func() {
	var (
		ema             *EMA
		period          int
		results         []float64
		expectedResults []float64
		err             error
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("ema_10_expectedresult.data")
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			ema, _ = NewEMA(period)
		})

		Context("given the source data has length greater than the lookback period", func() {

			BeforeEach(func() {
				results, err = ema.Calculate(TestInputData)
			})

			It("the result set should have a length equal to the source data length less the period + 1", func() {
				Expect(len(results)).To(Equal(len(TestInputData) - ema.LookbackPeriod + 1))
			})

			It("it should have correctly calculated the EMA for each item in the result set accurate to two decimal places", func() {
				for k := 0; k < len(results); k++ {
					Expect(expectedResults[k]).To(BeNumerically("~", results[k], 0.01))
				}
			})

			It("it should not return any errors", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("given the source data is nil", func() {
			BeforeEach(func() {
				results, err = ema.Calculate(nil)
			})

			It("it should return the appropriate error: ErrSourceDataEmpty", func() {
				Expect(err).To(Equal(ErrSourceDataEmpty))
			})
		})

		Context("given the source data has length less than the lookback period", func() {
			BeforeEach(func() {
				results, err = ema.Calculate(TestInputData[:8])
			})

			It("it should return the appropriate error: ErrNotEnoughSourceDataForLookbackPeriod", func() {
				Expect(err).To(Equal(ErrNotEnoughSourceDataForLookbackPeriod))
			})
		})

		Context("given the lookback period is less than or equal to zero", func() {
			BeforeEach(func() {
				period = -1
				ema, _ = NewEMA(period)
			})

			JustBeforeEach(func() {
				results, err = ema.Calculate(TestInputData[:8])
			})

			It("it should return the appropriate error: ErrLookbackPeriodMustBeGreaterThanZero", func() {
				Expect(err).To(Equal(ErrLookbackPeriodMustBeGreaterThanZero))
			})
		})
	})
})
