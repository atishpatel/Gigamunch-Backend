package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Config is the config.
type Config struct {
	ID            string `json:"id"`
	IsPowerSensor bool   `json:"is_power_sensor"`
	IntervalMins  int    `json:"interval_mins"`
	URL           string `json:"url"`
	LogFile       string `json:"log_file"`
}

// Device is a device that is reporting in for healthcheck.
type Device struct {
	ID            string    `json:"id"`
	IsPowerSensor bool      `json:"is_power_sensor"`
	LastCheckin   time.Time `json:"last_checkin"`
}

func main() {
	var err error
	// read config
	config := new(Config)
	filedata := readFile("device_config.json")
	err = json.Unmarshal(filedata, config)
	if err != nil {
		log.Println("failed to unmarshal config file device_config.json")
		log.Fatal(err)
	}
	// setup logging
	logf, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Println("Failed to create log file")
		log.Fatal(err)
	}
	log.SetOutput(logf)
	// run check every interval
	d := Device{
		ID:            config.ID,
		IsPowerSensor: config.IsPowerSensor,
	}
	for {
		d.LastCheckin = time.Now()
		err = checkin(config.URL, d)
		if err != nil {
			log.Printf("failed to checkin: %+v \n", err)
		}
		time.Sleep(time.Duration(config.IntervalMins) * time.Minute)
	}
}

func checkin(url string, body interface{}) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, b)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("failed to get correct response: %s", resp.Status)
	}
	return nil
}

func readFile(fileName string) []byte {
	filedata, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("failed to open %s", fileName)
		log.Fatal(err)
	}
	return filedata
}
