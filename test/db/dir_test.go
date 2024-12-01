package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/goleak"
)

func TestGetDirPath(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"),
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	tests := []struct {
		name    string
		env     map[string]string
		want    string
		wantErr bool
	}{
		{
			name:    "Success",
			env:     map[string]string{},
			want:    filepath.Dir(filepath.Join("server", "test", "db")),
			wantErr: false,
		},
		{
			name:    "Failure",
			env:     map[string]string{"MOCK_CALLER_FAILURE": "true"},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			defer func() {
				for k := range tt.env {
					os.Unsetenv(k)
				}
			}()

			got, err := GetDirPath()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDirPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				currentFileAbs, err := filepath.Abs("dir_test.go")
				if err != nil {
					t.Errorf("Failed to get absolute path of current file: %v", err)
					return
				}

				expectedAbs := filepath.Dir(currentFileAbs)

				if got != expectedAbs {
					t.Errorf("GetDirPath() = %v, want %v", got, expectedAbs)
				}
			}
		})
	}
}

func TestGetDirPath_Success(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"),
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	currentFileAbs, err := filepath.Abs("dir_test.go")
	if err != nil {
		t.Errorf("Failed to get absolute path of current file: %v", err)
		return
	}

	expectedAbs := filepath.Dir(currentFileAbs)

	dir, err := GetDirPath()
	if err != nil {
		t.Errorf("GetDirPath() error = %v", err)
		return
	}

	if dir != expectedAbs {
		t.Errorf("GetDirPath() = %v, want %v", dir, expectedAbs)
	}
}

func TestGetDirPath_Failure(t *testing.T) {
	defer goleak.VerifyNone(
		t,
		goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"),
		goleak.IgnoreTopFunction("net/http.(*Server).Serve"),
		goleak.IgnoreTopFunction("net/http/httptest.(*Server).goServe.func1"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"),
		goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
		goleak.IgnoreTopFunction("internal/poll.(*pollDesc).wait"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Accept"),
		goleak.IgnoreTopFunction("internal/poll.(*FD).Read"),
		goleak.IgnoreTopFunction("time.Sleep"),
		goleak.IgnoreTopFunction("time.AfterFunc"),
		goleak.IgnoreTopFunction("time.Ticker"),
		goleak.IgnoreTopFunction("runtime.gopark"),
		goleak.IgnoreTopFunction("runtime.forcegchelper"),
		goleak.IgnoreTopFunction("runtime.bgsweep"),
		goleak.IgnoreTopFunction("runtime.bgscavenge"),
	)

	os.Setenv("MOCK_CALLER_FAILURE", "true")
	defer os.Unsetenv("MOCK_CALLER_FAILURE")

	_, err := GetDirPath()
	if err == nil {
		t.Errorf("GetDirPath() should return an error, but got none")
		return
	}

	expectedError := "failed to get the current file path"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("GetDirPath() error = %v, want error containing %v", err, expectedError)
	}
}
