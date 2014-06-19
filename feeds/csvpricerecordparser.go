package feeds

import (
	"github.com/thetruetrade/gotrade"
	"strconv"
	"time"
)

type CSVDOHLCVRecordParser struct {
}

func (csvFPSP *CSVDOHLCVRecordParser) ParseRecord(csvRecord []string,
	dateColumnIndex int,
	openPriceColumnIndex int,
	highPriceColumnIndex int,
	lowPriceColumnIndex int,
	closePriceColumnIndex int,
	volumeColumnIndex int,
	dateParser TextDateParser) (dholcv gotrade.DOHLCV, err error) {

	// parse the record based on the supplied indexes and date func
	recordLength := len(csvRecord)

	// date
	var date time.Time
	if dateColumnIndex != -1 && recordLength > dateColumnIndex {

		date, err = dateParser(csvRecord[dateColumnIndex])
	}

	// open
	var open float64
	if openPriceColumnIndex != -1 && recordLength > openPriceColumnIndex {
		open, err = strconv.ParseFloat(csvRecord[openPriceColumnIndex], 64)
	}

	// high
	var high float64
	if recordLength > highPriceColumnIndex {
		high, err = strconv.ParseFloat(csvRecord[highPriceColumnIndex], 64)
	}

	// low
	var low float64
	if recordLength > lowPriceColumnIndex {
		low, err = strconv.ParseFloat(csvRecord[lowPriceColumnIndex], 64)
	}

	// close
	var close float64
	if recordLength > closePriceColumnIndex {
		close, err = strconv.ParseFloat(csvRecord[closePriceColumnIndex], 64)
	}

	// volume
	var volume float64
	if recordLength > volumeColumnIndex {
		volume, err = strconv.ParseFloat(csvRecord[volumeColumnIndex], 64)
	}

	dohlcv := gotrade.NewDOHLCVDataItem(date, open, high, low, close, volume)

	return dohlcv, err
}
