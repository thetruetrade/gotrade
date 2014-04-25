package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/thetruetrade/gotrade/indicators"
)

var _ = Describe("when calculating a SMA (simple moving average)", func() {
	var (
		sma             *SMA
		period          int
		results         []float64
		expectedResults []float64
		err             error
	)

	BeforeEach(func() {
		// load the expected results data
		expectedResults, _ = LoadCSVPriceDataFromFile("sma_10_expectedresult.data")
	})

	Describe("using a lookback period of 10", func() {

		BeforeEach(func() {
			period = 10
			sma, _ = NewSMA(period)
		})

		Context("given the source data has length greater than the lookback period", func() {

			BeforeEach(func() {
				results, err = sma.Calculate(TestInputData)
			})

			It("the result set should have a length equal to the source data length less the period + 1", func() {
				Expect(len(results)).To(Equal(len(TestInputData) - sma.LookbackPeriod + 1))
			})

			It("it should have correctly calculated the SMA for each item in the result set", func() {
				Expect(expectedResults).To(Equal(results))
			})

			It("it should not return any errors", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("given the source data is nil", func() {
			BeforeEach(func() {
				results, err = sma.Calculate(nil)
			})

			It("it should return the appropriate error: ErrSourceDataEmpty", func() {
				Expect(err).To(Equal(ErrSourceDataEmpty))
			})
		})

		Context("given the source data has length less than the lookback period", func() {
			BeforeEach(func() {
				results, err = sma.Calculate(TestInputData[:8])
			})

			It("it should return the appropriate error: ErrNotEnoughSourceDataForLookbackPeriod", func() {
				Expect(err).To(Equal(ErrNotEnoughSourceDataForLookbackPeriod))
			})
		})

		Context("given the lookback period is less than or equal to zero", func() {
			BeforeEach(func() {
				period = -1
				sma, _ = NewSMA(period)
			})

			JustBeforeEach(func() {
				results, err = sma.Calculate(TestInputData[:8])
			})

			It("it should return the appropriate error: ErrLookbackPeriodMustBeGreaterThanZero", func() {
				Expect(err).To(Equal(ErrLookbackPeriodMustBeGreaterThanZero))
			})
		})
	})
})
