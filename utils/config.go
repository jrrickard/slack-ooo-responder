package utils

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jrrickard/slack-ooo-responder/common"
)

func getLocalFile(file string) ([]byte, error) {
	return ioutil.ReadFile(strings.Replace(file, "file://", "", 1))
}

func getHTTPFile(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("Couldn't get file")
	}

	return ioutil.ReadAll(response.Body)

}

func getConfigFile(fileLocation string) (*common.Config, error) {
	var config common.Config
	var contents []byte
	var err error
	if strings.HasPrefix(fileLocation, "file://") {
		contents, err = getLocalFile(fileLocation)
	}
	contents, err = getHTTPFile(fileLocation)
	if err != nil {
		return &config, err
	}
	err = json.Unmarshal(contents, &config)
	return &config, err
}

func GetConfig() *common.Config {
	configFileLocation := flag.String("config", "", "config file location")
	flag.Parse()

	var config *common.Config
	var err error
	if *configFileLocation == "" {
		config = &common.Config{}
	} else {
		config, err = getConfigFile(*configFileLocation)
		if err != nil {
			log.Fatal("Error getting config file %v", err.Error())
		}
	}
	config.InitalizedTime = time.Now()
	getEnvironmentOverrides(config)
	validateConfig(config)
	setConfigDefaults(config)
	return config
}

func setConfigDefaults(config *common.Config) {
	if config.SurpressMessages <= 0 {
		config.SurpressMessages = 5
	}

	if config.Message == "" {
		config.Message = fmt.Sprintf("Hi! This is an automated response. I'm OOO until *%v*. I've included some helpful links below that might answer your question.", config.EndTime)
	}
}

func validateConfig(config *common.Config) {

	if config.StartTime.IsZero() {
		log.Fatal("Must provide start time")
	}

	if config.EndTime.IsZero() {
		log.Fatal("Must provide an end time")
	}
	if config.StartTime.After(config.EndTime) {
		log.Fatal("OOO Start must be before End")
	}

	if config.Token == "" {
		log.Fatal("Must provide a slack token")
	}
}

func getEnvironmentOverrides(config *common.Config) {
	apiToken := os.Getenv("SLACK_TOKEN")
	if apiToken != "" {
		config.Token = apiToken
	}

	start := os.Getenv("START_DATE")
	if start != "" {
		startTime, err := time.Parse(time.RFC3339, start)
		if err != nil {
			log.Printf("Invalid date %v", start)
		}
		config.StartTime = startTime
		config.Start = start
	}
	end := os.Getenv("END_DATE")
	if end != "" {
		endTime, err := time.Parse(time.RFC3339, end)
		if err != nil {
			log.Printf("Invalid date %v", end)
		}

		config.EndTime = endTime
		config.End = end
	}

	message := os.Getenv("OOM_MESSAGE")
	if message != "" {
		config.Message = message
	}
}
