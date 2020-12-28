package teamworkapi

import (
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

func TestInitLog(t *testing.T) {

	conf, err := LoadLogConfig("./testdata/conf.json")
	if err != nil {
		t.Errorf(err.Error())
	}

	err = InitLog(conf)

	log.Info("This is a test at INFO level.  Mary had a little lamb.")
	log.Debug("This is a test at DEBUG level.  Its fleece was white as snow.")
}