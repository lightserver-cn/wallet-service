package main

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	os.Setenv("TIMEZONE", "Asia/Shanghai")
	os.Setenv("ENV", "test")

	code := m.Run()
	goleak.VerifyTestMain(m)
	os.Exit(code)
}

func TestSetupEnvironment(t *testing.T) {
	defer goleak.VerifyNone(t)

	os.Setenv("TIMEZONE", "Asia/Shanghai")
	os.Setenv("ENV", "test")

	setupEnvironment()

	location, _ := time.LoadLocation("Asia/Shanghai")
	assert.Equal(t, location, time.Local)

	logOutput := captureLogOutput(func() {
		setupEnvironment()
	})
	expectedLog := "------ ENV:test TIMEZONE:Asia/Shanghai CurrentTime:"
	assert.Contains(t, logOutput, expectedLog)
}

func captureLogOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stdout)
	return buf.String()
}
