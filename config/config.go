package main

import (
	"encoding/json"
	"fmt"
	"github.com/schicho/mensa/canteen"
	"io"
	"log"
	"os"
	"time"
)

const filenameConfig = "mensa_conf"
const filenameCache = "mensa_data"

var filepathConfig string
var filepathCache string

// Config describes the json layout of the saved config file.
type Config struct {
	University string    `json:"university"`
	Cached     time.Time `json:"cached"`
}

// package level variable to have access to the loaded or generated configuration.
var config Config

// init the filepath to config and cache.
func init() {
	configDir, err := os.UserConfigDir()
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Println("Cannot access home directory.")
		log.Fatal(err)
	}
	filepathConfig = configDir + "/" + filenameConfig
	filepathCache = cacheDir + "/" + filenameCache
}

// loadConfig checks if there exists a previous configuration and loads it, or generates a new one and saves it to disk.
func loadConfig() {
	if exists(filepathConfig) {
		configFile, err := os.OpenFile(filepathConfig, os.O_RDONLY, os.ModePerm)
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
			fmt.Println("Malformed configuration .mensa file. Defaulting to Uni Passau.")
			config = Config{canteen.Canteens2Abbrev["UNI_PASSAU_CANTEEN"], time.Time{}}
			writeConfigFile()
		}
	} else {
		// Default to known values and create config file.
		fmt.Println("No configuration .mensa file. Creating new file. Defaulting to Uni Passau.")
		config = Config{canteen.Canteens2Abbrev["UNI_PASSAU_CANTEEN"], time.Time{}}
		writeConfigFile()
	}
}

// updateConfigFile just updates the timestamp in the configuration file, if new data was cached.
func updateConfigFile() {
	config = Config{config.University, time.Now()}
	writeConfigFile()
}

// writeConfigFile reads the data stored in the config variable, marshals it and writes it to disk.
func writeConfigFile() {
	buffer, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	configFile, err := os.OpenFile(filepathConfig, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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
	err := os.Remove(filepathConfig)
	if err != nil {
		log.Fatalln("Could not clear config.")
	}
	err = os.Remove(filepathCache)
	if err != nil {
		log.Fatalln("Could not clear cache.")
	}
}

// exists checks if a file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
