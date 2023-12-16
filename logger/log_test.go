package logger

import (
	"log"
	"os"
	"testing"
)

func TestInfof(t *testing.T) {
	Infof("This's is %s", "pandan")
}

func TestInfo(t *testing.T) {
	Info("one", "pandan")

	log.New(os.Stdout, "", log.LstdFlags).Println("one", "pandan")
}

func TestErrorf(t *testing.T) {
	Errorf("This's is %s", "pandan")
}

func TestWarnf(t *testing.T) {
	Warnf("This's is %s", "pandan")
}

func TestLogLevel(t *testing.T) {

	SetLogLevel(LogLevelSlient)
	Infof("This's is %s", "pandan")
	Warnf("This's is %s", "pandan")
	Errorf("This's is %s", "pandan")

	SetLogLevel(LogLevelInfo)
	Infof("This's is %s", "pandan")
	Warnf("This's is %s", "pandan")
	Errorf("This's is %s", "pandan")
}
