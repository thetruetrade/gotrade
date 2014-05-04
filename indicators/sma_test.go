package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
	"time"
)

var _ = Describe("when calculating a simple moving average (sma)", func() {
	var (
		period int = 3
	)

	Describe("given the sma target data structure is an array of floats ", func() {
		var (
			sma        *SMA
			sourceData = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 6.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 7.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 8.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0)}
		)

		BeforeEach(func() {
			sma, _ = NewSMA(period, gotrade.UseClosePrice)
		})

		Context("and the sma has received less ticks than the lookback period", func() {

			BeforeEach(func() {
				for i := 0; i < period-1; i++ {
					sma.RecieveOrderedTick(sourceData[i], i+1)
				}
			})

			It("the sma should have no result data", func() {
				Expect(len(sma.Data)).To(Equal(0))
			})
		})

		Context("and the sma has received ticks equal to the lookback period", func() {

			BeforeEach(func() {
				for i := 0; i <= period-1; i++ {
					sma.RecieveOrderedTick(sourceData[i], i+1)
				}
			})

			It("the sma should have result data with a single entry", func() {
				Expect(len(sma.Data)).To(Equal(1))
			})

			It("the sma should have a single result equal to the sum of the ticks divided by the lookback period", func() {
				sumData := 0.0
				for i := 0; i <= period-1; i++ {
					sumData += gotrade.UseClosePrice(sourceData[i])
				}
				Expect(sma.Data[0]).To(Equal(sumData / float64(period)))
			})
		})

		Context("and the sma has received more ticks than the lookback period", func() {

			BeforeEach(func() {
				for i := range sourceData {
					sma.RecieveOrderedTick(sourceData[i], i+1)
				}
			})

			It("the sma should have result data with entries equal to the number of ticks less the (lookback period - 1)", func() {
				Expect(len(sma.Data)).To(Equal(len(sourceData) - (period - 1)))
			})

			It("the sma should have a result for each tick equal to the sum of the ticks divided by the lookback period", func() {
				Expect(sma.Data[0]).To(Equal((5.0 + 6.0 + 7.0) / float64(period)))
				Expect(sma.Data[1]).To(Equal((6.0 + 7.0 + 8.0) / float64(period)))
				Expect(sma.Data[2]).To(Equal((7.0 + 8.0 + 9.0) / float64(period)))
			})
		})
	})

	Describe("given the sma target data structure is an array of bollinger band data items ", func() {
		var (
			sma        *SMAForAttachment
			sourceData = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 6.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 7.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 8.0, 0.0),
				gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0)}
			targetData []BollingerBandEntry
		)

		BeforeEach(func() {
			targetData = []BollingerBandEntry{}
			sma, _ = NewAttachedSMA(period, gotrade.UseClosePrice, func(dataItem float64, streamBarIndex int) {
				targetData = append(targetData, BollingerBandEntry{MiddleBand: dataItem})
			})
		})

		Context("and the sma has received less ticks than the lookback period", func() {

			BeforeEach(func() {
				for i := 0; i < period-1; i++ {
					sma.RecieveOrderedTick(sourceData[i], i+1)
				}
			})

			It("the sma should have no result data", func() {
				Expect(len(targetData)).To(Equal(0))
			})
		})

		Context("and the sma has received ticks equal to the lookback period", func() {

			BeforeEach(func() {
				for i := 0; i <= period-1; i++ {
					sma.RecieveOrderedTick(sourceData[i], i+1)
				}
			})

			It("the sma should have a single entry in the result data", func() {
				Expect(len(targetData)).To(Equal(1))
			})

			It("the sma should have a single result equal to the sum of the ticks divided by the lookbac period", func() {
				sumData := 0.0
				for i := 0; i <= period-1; i++ {
					sumData += gotrade.UseClosePrice(sourceData[i])
				}
				Expect(targetData[0].MiddleBand).To(Equal(sumData / float64(period)))
			})
		})

		Context("and the sma has received more ticks than the lookback period", func() {

			BeforeEach(func() {
				targetData = []BollingerBandEntry{}
				for i := 0; i < len(sourceData); i++ {
					sma.RecieveOrderedTick(sourceData[i], i+1)
				}
			})

			It("the sma should have result data with entries equal to the number of ticks less the (lookback period - 1)", func() {
				Expect(len(targetData)).To(Equal(len(sourceData) - (period - 1)))
			})

			It("the sma should have a result for each tick equal to the sum of the ticks divided by the lookback period", func() {

				Expect(targetData[0].MiddleBand).To(Equal((5.0 + 6.0 + 7.0) / float64(period)))
				Expect(targetData[1].MiddleBand).To(Equal((6.0 + 7.0 + 8.0) / float64(period)))
				Expect(targetData[2].MiddleBand).To(Equal((7.0 + 8.0 + 9.0) / float64(period)))
			})
		})
	})

})
