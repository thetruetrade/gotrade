package indicators_test

import (
	//"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
	//"math"
	"time"
)

var _ = Describe("when calculating bollinger bands", func() {
	var (
		period     int = 3
		bb         *BollingerBands
		sourceData = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 6.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 7.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 8.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0)}
	)

	BeforeEach(func() {
		bb, _ = NewBollingerBands(period, gotrade.UseClosePrice)
	})

	Context("when the bollinger band has received less ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i < period-1; i++ {
				bb.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the bollinger band should have no result data", func() {
			Expect(len(bb.Data)).To(Equal(0))
		})
	})

	Context("when the bollinger band has received ticks equal to the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i <= period-1; i++ {
				bb.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the bollinger band should have result data with a single entry", func() {
			Expect(len(bb.Data)).To(Equal(1))
		})

		It("the bollinger band should have a single result equal to the sum of the ticks divided by the lookback period", func() {
			sumData := 0.0
			for i := 0; i <= period-1; i++ {
				sumData += gotrade.UseClosePrice(sourceData[i])
			}
			Expect(bb.Data[0].MiddleBand).To(Equal(sumData / float64(period)))
		})
	})

	Context("when the bollinger band has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceData {
				bb.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the bollinger band should have result data with entries equal to the number of ticks less the (lookback period - 1)", func() {
			Expect(len(bb.Data)).To(Equal(len(sourceData) - (period - 1)))
		})

		It("the bollinger middle band should have a result for each tick equal to the sum of the ticks divided by the lookback period", func() {
			Expect(bb.Data[0].MiddleBand).To(Equal((5.0 + 6.0 + 7.0) / float64(period)))
			Expect(bb.Data[1].MiddleBand).To(Equal((6.0 + 7.0 + 8.0) / float64(period)))
			Expect(bb.Data[2].MiddleBand).To(Equal((7.0 + 8.0 + 9.0) / float64(period)))
		})

		//It("the bollinger band should have a +standard deviation result", func() {
		//	var s1 float64 = 5.0 + 6.0 + 7.0
		//	var n float64 = float64(period)
		//	var s2 float64 = (5.0 * 5.0) + (6.0 * 6.0) + (7.0 * 7.0)

		//	var variance = (s2/(n-1.0) - (s1/(n-1.0))*(s1/(n-1.0)))
		//	fmt.Println("VAR: ", variance)
		//	var stdDev = math.Sqrt(variance)
		//	fmt.Println("SD: ", stdDev)
		//	var final = bb.Data[0].MiddleBand + 2*stdDev
		//	Expect(bb.Data[0].UpperBand).To(Equal(final))
		//	//Expect(bb.Data[0].UpperBand).To(Equal(final))
		//	//			Expect(bb.Data[1].UpperBand).To(Equal((6.0 + 7.0 + 8.0) / float64(period)))
		//	//			Expect(bb.Data[2].UpperBand).To(Equal((7.0 + 8.0 + 9.0) / float64(period)))
		//})

		//It("the bollinger band should have a +standard deviation result", func() {

		//	Expect(bb.Data[0].LowerBand).To(Equal((5.0 ^ 2 + 6.0 + 7.0) / float64(period)))
		//	Expect(bb.Data[1].LowerBand).To(Equal((6.0 + 7.0 + 8.0) / float64(period)))
		//	Expect(bb.Data[2].LowerBand).To(Equal((7.0 + 8.0 + 9.0) / float64(period)))
		//})
	})
})
