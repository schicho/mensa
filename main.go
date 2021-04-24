package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/schicho/mensa/canteen"
	"github.com/schicho/mensa/csvutil"
	"github.com/schicho/mensa/download"
	"github.com/schicho/mensa/util"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	var clearConfigCache bool
	flag.BoolVar(&clearConfigCache, "c", false, "clear config and cache")
	flag.Parse()

	if clearConfigCache {
		deleteConfigCache()
		loadConfig()
		return
	}

	var canteenData io.Reader

	loadConfig()

	cachedYear, cachedWeek := config.Cached.ISOWeek()
	currentWeekday := time.Now().Weekday()
	currentYear, currentWeek := time.Now().ISOWeek()

	// Check if we can (still) use the cached data or need to download first and cache.
	if cachedWeek < currentWeek || cachedYear < currentYear || ((currentWeekday == time.
		Saturday || currentWeekday == time.Sunday) && config.Cached.Unix() < time.Now().Unix()) || !exists(filepathCache) {
		fmt.Println("Downloading new data...", canteen.Abbrev2Canteens[config.University])

		updateConfigFile()
		var err error
		canteenData, err = download.GetCSV(download.GenerateURL(config.University))
		if err != nil {
			panic(err)
		}

		// Cache the CSV in a file.
		cacheFile, err := os.OpenFile(filepathCache, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		defer cacheFile.Close()

		// Duplicate downloaded CSV data, store to disk and write back into the Reader.
		buffer := bytes.Buffer{}
		mw := io.MultiWriter(cacheFile, &buffer)
		_, err = mw.Write(util.ReaderToByte(canteenData))
		if err != nil {
			log.Println(err)
		}
		canteenData = bytes.NewReader(buffer.Bytes())

	} else {
		fmt.Println("Using cached data of", canteen.Abbrev2Canteens[config.University])
		canteenDataFile, err := os.OpenFile(filepathCache, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer canteenDataFile.Close()
		canteenData = bufio.NewReader(canteenDataFile)
	}

	var meals []*canteen.Dish

	gocsv.SetCSVReader(csvutil.NewSemicolonReader)
	err := gocsv.Unmarshal(canteenData, &meals)
	if err != nil {
		log.Fatalln(err)
	}

	// Shorten too long entries.
	for _, meal := range meals {
		if len(meal.Name) > 69 {
			meal.Name = meal.Name[:65] + "..."
		}
	}

	// Separate dishes by day and print.
	date := ""
	for _, meal := range meals {
		if meal.Date != date {
			fmt.Println(meal.Date, "-----------------------------")
			date = meal.Date
		}
		fmt.Printf("\t- %s\n", meal.Name)
	}
}
