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
			List<double> openPrices = new List<double>();
			List<double> closingPrices = new List<double>();
			List<double> highPrices = new List<double>();
			List<double> lowPrices = new List<double>();
			List<double> volume = new List<double>();
			// read the source data into an array to use for all the indicators
			using (var reader = new StreamReader (@"/home/eugened/Development/local/indicator-test-generator/indicator-test-generator/JSETOPI.2013.data")) 
			{
				string line = null;
				while((line = reader.ReadLine()) != null)
				{
					string[] parts = line.Split (new char[]{ ',' });
					// format is date, O, H, L, C, V
					// we will use close prices for all these tests
					openPrices.Add (Convert.ToDouble (parts [1].Replace (".", ",")));
					highPrices.Add (Convert.ToDouble (parts [2].Replace (".", ",")));
					lowPrices.Add (Convert.ToDouble (parts [3].Replace (".", ",")));
					closingPrices.Add (Convert.ToDouble (parts [4].Replace(".", ",")));
					volume.Add(Convert.ToDouble(parts[5].Replace (".", ",")));
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

			//True Range
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/truerange_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.TrueRangeLookback ();
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.TrueRange(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Average True Range
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/atr_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.AtrLookback (14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Atr(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), 14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Accumulation / Distribution line 
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/adl_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.AdLookback();
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Ad(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), volume.ToArray(), out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Chaikin Oscillator 
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/chaikinosc_3_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.AdOscLookback(3, 10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.AdOsc(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), volume.ToArray(), 3, 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// On Balance Volume 
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/obv_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.ObvLookback();
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.Obv(0, dataLength, closingPrices.ToArray(), volume.ToArray(), out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// AvgPrice
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/avgprice_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.AvgPriceLookback();
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.AvgPrice(0, dataLength, openPrices.ToArray(), highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// AvgPrice
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/medprice_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.MedPriceLookback();
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.MedPrice(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// PLUS_DM
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/plusdm_1_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.PlusDMLookback(1);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.PlusDM(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(),1, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// PLUS_DM
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/plusdm_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.PlusDMLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.PlusDM(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// MINUS_DM
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/minusdm_1_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.MinusDMLookback(1);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.MinusDM(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(),1, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// MINUS_DM
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/minusdm_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.MinusDMLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.MinusDM(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// PLUS_DI
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/plusdi_1_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.PlusDILookback(1);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.PlusDI(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), 1, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// PLUS_DI
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/plusdi_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.PlusDILookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.PlusDI(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// MINUS_DI
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/minusdi_1_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.MinusDILookback(1);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback + 1];
				talib.Core.RetCode retCode =talib.Core.MinusDI(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), 1, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// MINUS_DI
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/minusdi_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.MinusDILookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.MinusDI(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// DX
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/dx_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.DxLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Dx(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// ADX
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/adx_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.AdxLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Adx(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// ADXR
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/adxr_1_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.AdxrLookback(1);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Adxr(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(),1, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// ADXR
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/adxr_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.AdxrLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Adxr(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// TypicalPrice
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/typprice_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.TypPriceLookback();
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.TypPrice(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// RSI
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/rsi_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.RsiLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Rsi(0, dataLength, closingPrices.ToArray(), 14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// ROC
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/roc_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.RocLookback(10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Roc(0, dataLength, closingPrices.ToArray(), 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// ROCP
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/rocp_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.RocPLookback(10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.RocP(0, dataLength, closingPrices.ToArray(), 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// ROCR
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/rocr_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.RocRLookback(10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.RocR(0, dataLength, closingPrices.ToArray(), 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// ROCR
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/rocr100_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.RocR100Lookback(10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.RocR100(0, dataLength, closingPrices.ToArray(), 10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// MFI
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/mfi_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.MfiLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Mfi(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), volume.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// SAR
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/sar_002_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.SarLookback(0.02, 0.20);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Sar(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(),0.02, 0.20, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Linear Regression
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/linear_regression_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.LinearRegLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.LinearReg(0, dataLength, closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Linear Regression Slope
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/linear_regression_slope_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.LinearRegSlopeLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.LinearRegSlope(0, dataLength, closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Linear Regression Intercept
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/linear_regression_intercept_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.LinearRegInterceptLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.LinearRegIntercept(0, dataLength, closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Linear Regression Angle
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/linear_regression_angle_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.LinearRegAngleLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.LinearRegAngle(0, dataLength, closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// TSF
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/tsf_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.TsfLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Tsf(0, dataLength, closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// KAMA
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/kama_30_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.KamaLookback(30);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Kama(0, dataLength, closingPrices.ToArray(),30, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// TRIMA
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/trima_30_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.TrimaLookback(30);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Trima(0, dataLength, closingPrices.ToArray(),30, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// WILLR
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/willr_14_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.WillRLookback(14);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.WillR(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(),14, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// STOCH
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/stoch_5_3_3_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.StochLookback(5, 3, talib.Core.MAType.Sma, 3, talib.Core.MAType.Sma);
				int dataLength = closingPrices.Count - 1;
				double[] outSlowK = new double[dataLength - lookback +1];
				double[] outSlowD = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Stoch(0, dataLength, highPrices.ToArray(), lowPrices.ToArray(), closingPrices.ToArray(), 5,3, talib.Core.MAType.Sma, 3, talib.Core.MAType.Sma, out outBeginIndex, out outNBElement, outSlowK, outSlowD);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					for (var i=0;i< outSlowK.Length;i++) 
					{
						writer.WriteLine ("{0}, {1}", outSlowK[i].ToString(CultureInfo.InvariantCulture), outSlowD[i].ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// STOCHRSI
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/stochrsi_14_5_3_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.StochRsiLookback(14, 5, 3, talib.Core.MAType.Sma);
				int dataLength = closingPrices.Count - 1;
				double[] outFastK = new double[dataLength - lookback +1];
				double[] outFastD = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.StochRsi(0, dataLength, closingPrices.ToArray(), 14,5,3, talib.Core.MAType.Sma, out outBeginIndex, out outNBElement, outFastK, outFastD);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					for (var i=0;i< outFastK.Length;i++) 
					{
						writer.WriteLine ("{0}, {1}", outFastK[i].ToString(CultureInfo.InvariantCulture), outFastD[i].ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}

			// Momentum
			using (var writer = new StreamWriter (@"/home/eugened/Development/go/src/github.com/thetruetrade/gotrade/testdata/mom_10_expectedresult.data")) 
			{
				int outBeginIndex = 0;
				int outNBElement = 0;
				int lookback = talib.Core.MomLookback(10);
				int dataLength = closingPrices.Count - 1;
				double[] outData = new double[dataLength - lookback +1];
				talib.Core.RetCode retCode =talib.Core.Mom(0, dataLength, closingPrices.ToArray(),10, out outBeginIndex, out outNBElement, outData);
				if (retCode == TicTacTec.TA.Library.Core.RetCode.Success) 
				{
					foreach (var item in outData) 
					{
						writer.WriteLine (item.ToString(CultureInfo.InvariantCulture));
					}
				}
				writer.Flush ();
			}
		}
	}
}
