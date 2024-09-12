package system

import (
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/afero"
)

var DefaultFileSystem = afero.NewOsFs()

var CurrentUser = user.Current

var Stdout io.Writer = os.Stdout

var (
	originalFS   = DefaultFileSystem
	originalArgs = os.Args

	originalEnv     = Env()
	originalWorkDir = workdir()

	originalDefaultHTTPClient = http.DefaultClient
	originalDefaultTransport  = http.DefaultTransport
)

// Env returns the current environment variables as a usable map
func Env() map[string]string {
	r := map[string]string{}
	for _, value := range os.Environ() {
		key := strings.Split(value, "=")[0]
		r[key] = os.Getenv(key)
	}
	return r
}

func GetenvOrDefault(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func workdir() string {
	wd, err := os.Getwd()
	if err != nil {
		// explicitly ignore this error, just have a condition to please linters
	}
	return wd
}

// Reset restores the system environment as it was when the package was imported
func Reset() {
	os.Args = originalArgs
	http.DefaultClient = originalDefaultHTTPClient
	http.DefaultTransport = originalDefaultTransport
	CurrentUser = user.Current
	e := Env()
	for key, old := range originalEnv {
		if new, ok := e[key]; !ok || new != old {
			os.Setenv(key, old)
		}
	}
	for key := range e {
		if _, ok := originalEnv[key]; !ok {
			os.Unsetenv(key)
		}
	}
	DefaultFileSystem = originalFS
	if originalWorkDir != "" {
		err := os.Chdir(originalWorkDir)
		if err != nil {
			// explicitly ignore this error, just have a condition to please linters
		}
	}
	Stdout = os.Stdout
}
