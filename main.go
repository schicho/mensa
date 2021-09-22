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

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Prints the meals served this week at the canteen.\n")
		fmt.Fprint(os.Stderr, "Edit the configuration file in .config/mensa_conf.json to switch to a different\n")
		fmt.Fprint(os.Stderr, "university or to a different pricing (Students/Employees/Guests).\n\n")
		fmt.Fprint(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\nMensa CLI - https://github.com/schicho/mensa\n")
	}
}

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

	configuration := config.GetConfig()

	if clearConfigCache {
		configuration.BuildNewConfig()
		return
	}

	if printNoColor {
		colorReset = ""
		colorRed = ""
	}

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

		configuration.UpdateConfigFile()
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
	var mealsByDay [7][]*canteen.Dish

	gocsv.SetCSVReader(csvutil.NewSemicolonReader)
	err = gocsv.UnmarshalBytes(canteenData, &meals)
	if err != nil {
		log.Fatalln(err)
	}

	// Sort meals by weekdays.
	for _, meal := range meals {
		// Shorten too long entries.
		if len(meal.Name) > 60 {
			meal.Name = substring.SubstringEnd(meal.Name, 60-len(meal.MealType)) + "..."
		}
		timestamp, err := time.Parse("02.01.2006", meal.Date)
		if err != nil {
			panic(err)
		}
		weekday := timestamp.Weekday()
		mealsByDay[weekday] = append(mealsByDay[weekday], meal)
	}

	// Print dishes sorted by weekdays.
	for day, meals := range mealsByDay {
		if len(meals) > 0 {
			if !printTodayOnly || day == int(currentWeekday) {
				fmt.Printf("%s%v %v:%s\n", colorRed, meals[0].Date, time.Weekday(day), colorReset)
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
