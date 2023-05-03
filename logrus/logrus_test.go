package main

import (
	"bytes"
	"github.com/sirupsen/logrus"

	"io"
	"log"
	"os"
	"testing"
)

func TestAddCallerAndMethod(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.Info("info msg")
}
func TestAddField(t *testing.T) {
	logrus.WithFields(logrus.Fields{
		"name": "dj",
		"age":  18,
	}).Info("info msg")
}
func TestRedirect(t *testing.T) {
	writer1 := &bytes.Buffer{}
	writer2 := os.Stdout
	writer3, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(writer1, writer2, writer3))
	logrus.Info("info msg")

}
func TestJson(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.Trace("trace msg")
	logrus.Debug("debug msg")
	logrus.Info("info msg")
	logrus.Warn("warn msg")
	logrus.Error("error msg")
	logrus.Fatal("fatal msg")
	logrus.Panic("panic msg")

}
