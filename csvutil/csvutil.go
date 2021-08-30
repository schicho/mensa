package csvutil

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/charmap"
)

var replaceNewlineSemicolon = strings.NewReplacer("\n;", ";")

// NewSemicolonReader is a stdlib CSV reader, but with a semicolon as separator, as it's needed for the files we parse.
func NewSemicolonReader(in io.Reader) gocsv.CSVReader {
	r := csv.NewReader(in)
	r.LazyQuotes = true
	r.Comma = ';'
	return r
}

// FixCSVFormatting is needed as for some reason the online provided file has line breaks in seemingly random entries.
// Fortunately they only consist of newlines followed by semicolons, (Basically one CSV entry is split over 2 lines)
// thus is easily fixable.
func FixCSVFormatting(in string) string {
	return replaceNewlineSemicolon.Replace(in)
}

// Windows1252ToUTF8 The CSV is Windows1252 encoded, but we want UTF8 to work with strings and cache the data.
func Windows1252ToUTF8(in []byte) string {
	decoder := charmap.Windows1252.NewDecoder()
	bufUTF8 := make([]byte, len(in)*3)
	n, _, err := decoder.Transform(bufUTF8, in, false)
	if err != nil {
		panic(err)
	}
	bufUTF8 = bufUTF8[:n]
	return string(bufUTF8)
}
