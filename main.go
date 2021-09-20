package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/schicho/mensa/canteen"
	"github.com/schicho/mensa/config"
	"github.com/schicho/mensa/csvutil"
	"github.com/schicho/mensa/download"
	"github.com/schicho/substring"
)

var colorReset = "\033[0m"
var colorRed = "\033[31;1m"

func main() {
	var clearConfigCache bool
	var forceDownloadData bool
	var printTodayOnly bool
	var printNoColor bool
	flag.BoolVar(&clearConfigCache, "c", false, "clear config and cache")
	flag.BoolVar(&forceDownloadData, "d", false, "force downloading/updating the canteen data")
	flag.BoolVar(&printTodayOnly, "t", false, "only print the meals served today")
	flag.BoolVar(&printNoColor, "n", false, "do not color the output text (no red weekdays)")
	flag.Parse()

	if clearConfigCache {
		config.BuildNewConfig()
		return
	}

	if printNoColor {
		colorReset = ""
		colorRed = ""
	}

	// Have main package configuration equal to the one in the config package
	config.LoadConfig()
	configuration := *config.GetConfig()

	cachedYear, _ := configuration.Cached.ISOWeek()
	currentYear, _ := time.Now().ISOWeek()
	cachedDay := configuration.Cached.YearDay()
	currentDay := time.Now().YearDay()
	currentWeekday := time.Now().Weekday()

	var canteenData []byte
	var err error

	// Check if we can (still) use the cached data or need to download first and cache.
	if forceDownloadData || cachedDay < currentDay || cachedYear < currentYear || ((currentWeekday == time.
		Saturday || currentWeekday == time.Sunday) && cachedDay < currentDay) || !config.Exists(config.FilepathCache) {

		fmt.Fprintln(os.Stderr, "Downloading new data...", canteen.Abbrev2Canteens[configuration.University])

		config.UpdateConfigFile()
		canteenData, err = download.GetCSV(download.GenerateURL(configuration.University))
		if err != nil {
			panic(err)
		}

		// Cache the CSV in a file.
		cacheFile, err := os.OpenFile(config.FilepathCache, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		defer cacheFile.Close()

		_, err = cacheFile.Write(canteenData)
		if err != nil {
			log.Println(err)
		}

	} else {
		fmt.Fprintln(os.Stderr, "Using cached data of", canteen.Abbrev2Canteens[configuration.University])
		canteenDataFile, err := os.OpenFile(config.FilepathCache, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer canteenDataFile.Close()
		canteenData, err = ioutil.ReadAll(canteenDataFile)
		if err != nil {
			log.Println(err)
		}
	}

	var meals []*canteen.Dish

	gocsv.SetCSVReader(csvutil.NewSemicolonReader)
	err = gocsv.UnmarshalBytes(canteenData, &meals)
	if err != nil {
		log.Fatalln(err)
	}

	// Shorten too long entries.
	for _, meal := range meals {
		if len(meal.Name) > 60 {
			meal.Name = substring.SubstringEnd(meal.Name, 60-len(meal.MealType)) + "..."
		}
	}

	// Separate dishes by day.
	type MealDay struct {
		*canteen.Dish
		Weekday time.Weekday
	}
	var mealsByDay [7][]MealDay

	for _, meal := range meals {
		timestamp, err := time.Parse("02.01.2006", meal.Date)
		if err != nil {
			panic(err)
		}
		weekday := timestamp.Weekday()
		mealsByDay[weekday] = append(mealsByDay[weekday], MealDay{meal, weekday})
	}

	// Print dishes sorted by weekdays.
	todayWeekday := time.Now().Weekday()
	for _, meals := range mealsByDay {
		if len(meals) > 0 {
			if !printTodayOnly || meals[0].Weekday == todayWeekday {
				fmt.Printf("%s%v %v:%s\n", colorRed, meals[0].Date, meals[0].Weekday, colorReset)
				for _, meal := range meals {
					switch configuration.Price {
					case config.PriceStudent_t:
						fmt.Printf("    - %s : %s [%s]\n", meal.PriceStudent, meal.Name, meal.MealType)
					case config.PriceEmployee_t:
						fmt.Printf("    - %s : %s [%s]\n", meal.PriceEmployee, meal.Name, meal.MealType)
					case config.PriceGuest_t:
						fmt.Printf("    - %s : %s [%s]\n", meal.PriceGuest, meal.Name, meal.MealType)
					}
				}
			}
		}
	}
}
