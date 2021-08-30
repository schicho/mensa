package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/schicho/mensa/canteen"
)

const FilenameConfig = "mensa_conf.json"
const FilenameCache = "mensa_data.csv"

var FilepathConfig string
var FilepathCache string

type PriceType string

const (
	PriceStudent_t  PriceType = "Student"
	PriceEmployee_t PriceType = "Employee"
	PriceGuest_t    PriceType = "Guest"
)

// Config describes the json layout of the saved config file.
type Config struct {
	University string    `json:"university"`
	Cached     time.Time `json:"cached"`
	Price      PriceType `json:"price"`
}

// config stores the currently loaded user data configuration
var config Config

// GetConfig returns a pointer to the loaded user data
func GetConfig() *Config {
	return &config
}

// init the filepath to config and cache.
func init() {
	configDir, err := os.UserConfigDir()
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Println("Cannot access home directory.")
		log.Fatal(err)
	}
	FilepathConfig = configDir + "/" + FilenameConfig
	FilepathCache = cacheDir + "/" + FilenameCache
}

// LoadConfig checks if there exists a previous configuration and loads it, or generates a new one and saves it to disk.
func LoadConfig() {
	if Exists(FilepathConfig) {
		configFile, err := os.OpenFile(FilepathConfig, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer configFile.Close()

		buffer, err := io.ReadAll(configFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(buffer, &config)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Malformed `mensa_conf.json` file. Default to Student of Uni Passau.")
			config = Config{canteen.Canteens2Abbrev["UNI_PASSAU_CANTEEN"], time.Time{}, PriceStudent_t}
			writeConfigFile()
		}
	} else {
		// Default to known values and create config file.
		fmt.Fprintln(os.Stderr, "No `mensa_conf.json` file. Creating new file. Default to Student of Uni Passau.")
		config = Config{canteen.Canteens2Abbrev["UNI_PASSAU_CANTEEN"], time.Time{}, PriceStudent_t}
		writeConfigFile()
	}
}

// UpdateConfigFile just updates the timestamp in the configuration file, if new data was cached.
func UpdateConfigFile() {
	config = Config{config.University, time.Now(), config.Price}
	writeConfigFile()
}

// writeConfigFile reads the data stored in the config variable, marshals it and writes it to disk.
func writeConfigFile() {
	buffer, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	configFile, err := os.OpenFile(FilepathConfig, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println("Could not create configuration file.")
		return
	}
	defer configFile.Close()

	_, err = configFile.Write(buffer)
	if err != nil {
		return
	}
}

func deleteConfigCache() {
	fmt.Fprintln(os.Stderr, "Deleting config and cache.")
	err := os.Remove(FilepathConfig)
	if err != nil {
		log.Fatalln("Could not clear config.")
	}
	err = os.Remove(FilepathCache)
	if err != nil {
		log.Fatalln("Could not clear cache.")
	}
}

// BuildNewConfig deletes the config and cache files from disk and creates a new default config file.
func BuildNewConfig() {
	deleteConfigCache()
	LoadConfig()
}

// Exists checks if a file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
