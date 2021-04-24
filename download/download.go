package download

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/schicho/mensa/csvutil"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// exampleUrl = "https://www.stwno.de/infomax/daten-extern/csv/UNI-P/16.csv"
	urlPrefix  = "https://www.stwno.de/infomax/daten-extern/csv"
	urlPostfix = ".csv"
)

// GetCSV gets the CSV via the internet. Fixes the formatting errors in the csv.
// Finally provides a new io.Reader to read correct CSVs from.
func GetCSV(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("could not get file")
	}
	defer resp.Body.Close()

	responseString := csvutil.Windows1252ToUTF8(ReaderToByte(resp.Body))

	// Fix formatting early, so we don't need to bother later.
	responseString = csvutil.FixCSVFormatting(responseString)

	return strings.NewReader(responseString), nil
}

// GenerateURL for the asked university and current/next week.
func GenerateURL(universityAbbrev string) string {
	now := time.Now()
	day := now.Weekday()
	_, week := time.Time.ISOWeek(now)

	// Get URL to the data of the next week on the weekend.
	if day == time.Saturday || day == time.Sunday {
		week++
	}
	return fmt.Sprintf("%s/%s/%d%s", urlPrefix, universityAbbrev, week, urlPostfix)
}

// ReaderToByte reads an io.Reader into a byte slice and returns it.
func ReaderToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}