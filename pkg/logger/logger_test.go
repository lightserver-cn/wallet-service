package logger

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"os"
	"testing"
)

const logFileExt = "log"

func TestInitZapLogger(t *testing.T) {
	// Create a temporary directory for storing log files
	tempDir, err := os.MkdirTemp("", "test-logs")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	filepath := tempDir + "/logs"
	infoFilename := "info.log"
	warnFilename := "warn.log"
	errFilename := "error.log"
	fileExt := logFileExt
	callerLoc := "wallet-service"

	logger, err := InitZapLogger(filepath, infoFilename, warnFilename, errFilename, fileExt, callerLoc)
	require.NoError(t, err)
	assert.NotNil(t, logger)

	// Test log output
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
}

func TestNewZapLogger(t *testing.T) {
	// Create a temporary directory for storing log files
	tempDir, err := os.MkdirTemp("", "test-logs")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFilepath := tempDir + "/logs"
	infoFilename := "info.log"
	warnFilename := "warn.log"
	errFilename := "error.log"
	fileExt := logFileExt
	callerLoc := "wallet-service"

	cfg := initConf(logFilepath, callerLoc)

	logger, err := newZapLogger(logFilepath, infoFilename, warnFilename, errFilename, fileExt, cfg)
	require.NoError(t, err)
	assert.NotNil(t, logger)

	// Test log output
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
}

func TestGetLogWriter(t *testing.T) {
	// Create a temporary directory for storing log files
	tempDir, err := os.MkdirTemp("", "test-logs")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	filePath := tempDir + "/test.log"
	fileExt := logFileExt

	writer, err := getLogWriter(filePath, fileExt)
	require.NoError(t, err)
	assert.NotNil(t, writer)
}

func TestGetWriter(t *testing.T) {
	// Create a temporary directory for storing log files
	tempDir, err := os.MkdirTemp("", "test-logs")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	filename := tempDir + "/test.log"
	fileExt := logFileExt

	writer, err := getWriter(filename, fileExt)
	require.NoError(t, err)
	assert.NotNil(t, writer)

	// Write a test log message
	_, _ = writer.Write([]byte("Test log message"))
}
