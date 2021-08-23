package config

import (
	"encoding/json"
	"fmt"
	"github.com/schicho/mensa/canteen"
	"io"
	"log"
	"os"
	"time"
)

const FilenameConfig = "mensa_conf"
const FilenameCache = "mensa_data"

var FilepathConfig string
var FilepathCache string

// Config describes the json layout of the saved config file.
type Config struct {
	University string    `json:"university"`
	Cached     time.Time `json:"cached"`
}

var config Config

func GetConfig() *Config{
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

// UpdateConfigFile just updates the timestamp in the configuration file, if new data was cached.
func UpdateConfigFile() {
	config = Config{config.University, time.Now()}
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
	err := os.Remove(FilepathConfig)
	if err != nil {
		log.Fatalln("Could not clear config.")
	}
	err = os.Remove(FilepathCache)
	if err != nil {
		log.Fatalln("Could not clear cache.")
	}
}

func BuildNewConfig() {
	deleteConfigCache()
	LoadConfig()
}

// exists checks if a file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
