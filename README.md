# gotrade

Golang stock market technical analysis library

[![Build Status](https://travis-ci.org/thetruetrade/gotrade.svg?branch=dev)](https://travis-ci.org/thetruetrade/gotrade)


[![Stories in Ready](https://badge.waffle.io/thetruetrade/gotrade.png?label=ready&title=Ready)](https://waffle.io/thetruetrade/gotrade)


GoTrade is in early design and development

Below is a look at the basic API so far

```go
	csvFeed := feeds.NewCSVFileFeedWithDOHLCVFormat("../github.com/thetruetrade/gotrade/testdata/JSETOPI.2013.data",
		feeds.DashedYearDayMonthDateParserForLocation(time.Local))

	priceStream := gotrade.NewDailyDOHLCVStream()
	sma, _ := indicators.NewSMAForStream(priceStream, 20, gotrade.UseClosePrice)
	ema, _ := indicators.NewEMAForStream(priceStream, 20, gotrade.UseClosePrice)
	bb, _ := indicators.NewBollingerBandsForStream(priceStream, 20, gotrade.UseClosePrice)

	csvFeed.FillDOHLCVStream(priceStream)

	fmt.Println("price stream has data of length: ", len(priceStream.Data))
	fmt.Println("price stream has min date: ", priceStream.MinDate())
	fmt.Println("price stream has max date: ", priceStream.MaxDate())

	fmt.Println("sma has data of length: ", len(sma.Data))
	fmt.Println("sma is valid from price stream bar number: ", sma.ValidFromBar())
	fmt.Println("sma max: ", sma.MaxValue(), " sma min: ", sma.MinValue())

	fmt.Println("ema has data of length: ", len(ema.Data))
	fmt.Println("ema is valid from price stream bar number: ", ema.ValidFromBar())
	fmt.Println("ema max: ", ema.MaxValue(), " ema min: ", ema.MinValue())

	fmt.Println("bollinger bands has data of length: ", len(bb.Data))
	fmt.Println("bollinger bands is valid from price stream bar number: ", bb.ValidFromBar())
	fmt.Println("bollinger bands max: ", bb.MaxValue(), " sma min: ", bb.MinValue())
```

Tasks for the near future include:

 * Complete a basic set of indicators
 * Add point and figure price streams
 * Basic pattern matching
   * Candlesticks
   * Point and figure patterns
 * Visualisation in [gotrade-plot](https://github.com/thetruetrade/gotrade-plot)
 * Operator support like Crosses etc.
 * Script engine in [gotrade-script](https://github.com/thetruetrade/gotrade-script)