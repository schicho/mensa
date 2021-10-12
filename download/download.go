package download

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/schicho/mensa/csvutil"
)

const (
	// exampleUrl = "https://www.stwno.de/infomax/daten-extern/csv/UNI-P/16.csv"
	urlPrefix  = "https://www.stwno.de/infomax/daten-extern/csv"
	urlPostfix = ".csv"
)

// GetCSV gets the CSV via the internet. Fixes the formatting errors in the csv.
// Finally provides a new io.Reader to read correct CSVs from.
func GetCSV(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("could not download mensa data")
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Could not read response.")
	}

	responseString := csvutil.Windows1252ToUTF8(responseBytes)

	// Fix formatting early, so we don't need to bother later.
	responseString = csvutil.FixCSVFormatting(responseString)
	responseString = csvutil.RemoveBrackets(responseString)

	return []byte(responseString), nil
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
