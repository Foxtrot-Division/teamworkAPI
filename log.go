package teamworkapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// LogConfig models a logrus log configuration.
type LogConfig struct {
	LogToFile 	bool	`json:"logToFile"`
	LogFilePath string 	`json:"logFilePath"`
}

// LoadLogConfig loads log configuration from json file specified by path.
func LoadLogConfig(path string) (*LogConfig, error) {

	f, err := os.Open(path)
	defer f.Close()
	
	if err != nil {
		return nil, fmt.Errorf("failed to open config file at " + path)
	}

	byteValue, _ := ioutil.ReadAll(f)
	l := new(LogConfig)

	err = json.Unmarshal(byteValue, &l)
	if err != nil {
		return nil, err
	}

	if l.LogToFile && l.LogFilePath == "" {
		return nil, fmt.Errorf("failed to load log file path from file (%s)", path)
	}

	return l, nil
}

// InitLog initializes logrus instance based on conf parameters.
func InitLog(conf *LogConfig) (error) {

	if conf.LogToFile {
		f, err := os.OpenFile(conf.LogFilePath, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file (%s)", conf.LogFilePath)
		}

		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	}	
	
	// add function name to log message
	log.SetReportCaller(true)

	log.SetLevel(log.InfoLevel)
	
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	return nil
}