using System;
using System.IO;
using System.Collections;
using System.Collections.Generic;
using System.Globalization;

using talib = TicTacTec.TA.Library;

namespace indicatortestgenerator
{
	class MainClass
	{
		public static void Main (string[] args)
		{
			List<double> closingPrices = new List<double>();
			List<double> highPrices = new List<double>();
			List<double> lowPrices = new List<double>();
			// read the source data into an array to use for all the indicators
			using (var reader = new StreamReader (@"/home/eugened/Development/local/indicator-test-generator/indicator-test-generator/JSETOPI.2013.data")) 
			{
				string line = null;
				while((line = reader.ReadLine()) != null)
				{
					string[] parts = line.Split (new char[]{ ',' });
					// format is date, O, H, L, C, V
					// we will use close prices for all these tests
					highPrices.Add (Convert.ToDouble (parts [2].Replace (".", ",")));
					closingPrices.Add (Convert.ToDouble (parts [4].Replace(".", ",")));
					lowPrices.Add (Convert.ToDouble (parts [3].Replace (".", ",")));
				}
			}

			// now we need to create an output file for each indicator

			// SMA
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/sma_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.SmaLookback (10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Sma(0, dataLength, closingPrices.ToArray(), 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}



			// EMA
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/ema_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.EmaLookback (10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Ema(0, dataLength, closingPrices.ToArray(), 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// WMA
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/wma_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.WmaLookback (10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Wma(0, dataLength, closingPrices.ToArray(), 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// DEMA
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/dema_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.DemaLookback (10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Dema(0, dataLength, closingPrices.ToArray(), 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// TEMA
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/tema_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.TemaLookback (10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Tema(0, dataLength, closingPrices.ToArray(), 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Variance
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/variance_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.VarianceLookback (10, 1.0);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Variance(0, dataLength, closingPrices.ToArray(), 10, 1.0, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Standard Deviation
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/stddev_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.StdDevLookback (10, 1.0);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.StdDev(0, dataLength, closingPrices.ToArray(), 10, 1.0, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Bollinger Bands
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/bb_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.SmaLookback (10);
				int dataLength = closingPrices.Count - 1;
				double[] outDataUpper = new double[dataLength - lookback + 1];
				double[] outDataMiddle = new double[dataLength - lookback + 1];
				double[] outDataLower = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Bbands(0, dataLength, closingPrices.ToArray(), 10, 2, 2, talib.Core.MAType.Sma, out outBeginIndex, out outNBElement, outDataUpper, outDataMiddle, outDataLower);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					for (int i= 0; i< outDataMiddle.Length;i++) 
					{
						writer.WriteLine ("{0}, {1}, {2}", outDataUpper[i].ToString(CultureInfo.InvariantCulture), outDataMiddle[i].ToString(CultureInfo.InvariantCulture), outDataLower[i].ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// MACD
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/macd_12_26_9_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.MacdLookback (12, 26, 9);
				int dataLength = closingPrices.Count - 1;
				double[] outMACD = new double[dataLength - lookback + 1];
				double[] outMACDSignal = new double[dataLength - lookback + 1];
				double[] outMACDHist = new double[dataLength - lookback + 1];

				talib.Core.SetUnstablePeriod (talib.Core.FuncUnstId.FuncUnstAll, 0); 
				talib.Core.RetCode retCode =talib.Core.Macd(0, dataLength, closingPrices.ToArray(), 12, 26, 9, out outBeginIndex, out outNBElement, outMACD, outMACDSignal, outMACDHist);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					for (int i= 0; i< outMACD.Length;i++) 
					{
						writer.WriteLine ("{0}, {1}, {2}", outMACD[i].ToString(CultureInfo.InvariantCulture), outMACDSignal[i].ToString(CultureInfo.InvariantCulture), outMACDHist[i].ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Aroon
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/aroon_25_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.AroonLookback (25);
				int dataLength = closingPrices.Count - 1;
				double[] outAroonDown = new double[dataLength - lookback + 1];
				double[] outAroonUp = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Aroon(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), 25, out outBeginIndex, out outNBElement, outAroonDown, outAroonUp);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					for (int i= 0; i< outAroonUp.Length;i++) 
					{
						writer.WriteLine ("{0}, {1}", outAroonUp[i].ToString(CultureInfo.InvariantCulture), outAroonDown[i].ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// AroonOsc
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/aroonosc_25_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.AroonOscLookback (25);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.AroonOsc(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), 25, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					for (int i= 0; i< outData.Length;i++) 
					{
						writer.WriteLine ("{0}", outData[i].ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

		}
	}
}
