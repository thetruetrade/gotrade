package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
	"time"
)

var _ = Describe("when calculating a weighted moving average (wma)", func() {
	var (
		period     int = 3
		weightsSum int = 3 + 2 + 1
	)

	Describe("given the wma target data structure is an array of floats ", func() {
		var (
			wma        *WMA
			sourceData = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 6.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 7.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 8.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0)}
		)

		BeforeEach(func() {
			wma, _ = NewWMA(period, gotrade.UseClosePrice)
		})

		Context("and the wma has received less ticks than the lookback period", func() {

			BeforeEach(func() {
				for i := 0; i < period-1; i++ {
					wma.RecieveOrderedTick(sourceData[i], i+1)
				}
			})

			It("the wma should have no result data", func() {
				Expect(len(wma.Data)).To(Equal(0))
			})
		})

		Context("and the wma has received ticks equal to the lookback period", func() {

			BeforeEach(func() {
				for i := 0; i <= period-1; i++ {
					wma.RecieveOrderedTick(sourceData[i], i+1)
				}
			})

			It("the wma should have result data with a single entry", func() {
				Expect(len(wma.Data)).To(Equal(1))
			})

			It("the wma should have a single result equal to the sum of the ticks multiplied by their weights divided by the sum of the weights", func() {
				sumData := 0.0
				for i := 0; i <= period-1; i++ {
					sumData += gotrade.UseClosePrice(sourceData[i]) * float64(i+1)
				}
				Expect(wma.Data[0]).To(Equal(sumData / float64(weightsSum)))
			})
		})

		Context("and the wma has received more ticks than the lookback period", func() {

			BeforeEach(func() {
				for i := range sourceData {
					wma.RecieveOrderedTick(sourceData[i], i+1)
				}
			})

			It("the wma should have result data with entries equal to the number of ticks less the (lookback period - 1)", func() {
				Expect(len(wma.Data)).To(Equal(len(sourceData) - (period - 1)))
			})

			It("the wma should have a result for each tick equal to the sum of the ticks multiplied by their weights divided by the sum of the weights", func() {
				Expect(wma.Data[0]).To(Equal(((5.0 * 1) + (6.0 * 2) + (7.0 * 3)) / float64(weightsSum)))
				Expect(wma.Data[1]).To(Equal(((6.0 * 1) + (7.0 * 2) + (8.0 * 3)) / float64(weightsSum)))
				Expect(wma.Data[2]).To(Equal(((7.0 * 1) + (8.0 * 2) + (9.0 * 3)) / float64(weightsSum)))
			})
		})
	})
})
