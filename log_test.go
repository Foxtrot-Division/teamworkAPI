package teamworkapi

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestLoadLogConfig(t *testing.T) {
	
	conf, err := LoadLogConfig("./testdata/conf.json")
	if err != nil {
		t.Errorf(err.Error())
	}

	if !conf.LogToFile {
		t.Errorf("expected LogToFile to be true but got false")
	}
}

func TestInit(t *testing.T) {

	log.Info("This is a test at INFO level.  Mary had a little lamb.")
	log.Debug("This is a test at DEBUG level.  Its fleece was white as snow.")

	f, err := os.Open("./testdata/teamworkapi.log")
	defer f.Close()
	
	if err != nil {
		t.Errorf("failed to open log file (%s)", err.Error())
	}
}