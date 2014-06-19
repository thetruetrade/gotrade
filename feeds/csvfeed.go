package feeds

import (
	"encoding/csv"
	"github.com/thetruetrade/gotrade"
	"io"
	"os"
)

type CSVFileFeed struct {
	*CSVDOHLCVRecordParser
	fileName              string
	dateColumnIndex       int
	openPriceColumnIndex  int
	highPriceColumnIndex  int
	lowPriceColumnIndex   int
	closePriceColumnIndex int
	volumeColumnIndex     int
	dateParser            TextDateParser
}

func NewCSVFileFeedWithDOHLCVFormat(fileName string,
	dateParser TextDateParser) *CSVFileFeed {

	return &CSVFileFeed{&CSVDOHLCVRecordParser{},
		fileName,
		0,
		1,
		2,
		3,
		4,
		5,
		dateParser}
}

func NewCSVFileFeed(fileName string,
	dateColumnIndex int,
	openPriceColumnIndex int,
	highPriceColumnIndex int,
	lowPriceColumnIndex int,
	closePriceColumnIndex int,
	volumeColumnIndex int,
	dateParser TextDateParser) *CSVFileFeed {

	return &CSVFileFeed{&CSVDOHLCVRecordParser{},
		fileName,
		dateColumnIndex,
		openPriceColumnIndex,
		highPriceColumnIndex,
		lowPriceColumnIndex,
		closePriceColumnIndex,
		volumeColumnIndex,
		dateParser}
}

func (csvFPSF *CSVFileFeed) FillDOHLCVStream(priceStream gotrade.DOHLCVStreamTickReceiver) (err error) {

	file, err := os.Open(csvFPSF.fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var lineNumber int = 0
	for {

		// increment the linenumbers
		lineNumber++

		// read the record from the file
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		dohlcv, err := csvFPSF.ParseRecord(record, csvFPSF.dateColumnIndex,
			csvFPSF.openPriceColumnIndex,
			csvFPSF.highPriceColumnIndex,
			csvFPSF.lowPriceColumnIndex,
			csvFPSF.closePriceColumnIndex,
			csvFPSF.volumeColumnIndex,
			csvFPSF.dateParser)

		if err != nil {
			return err
		}
		priceStream.ReceiveTick(dohlcv)
	}
	return nil
}
