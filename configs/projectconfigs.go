package configs

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// *****************************
// *** Common ***
// *****************************

// DryRun is used to protect and prevent external calls from executing while not in the Live environment.
var DryRun = GetBoolEnvVar("GO_DAEMONS_DRY_RUN", true)

// Environment is the environment we are running in: PROD | GCP | LOCAL
var Environment = getEnvironment()

// EnvironmentLocal is the environment variable for LOCAL
var EnvironmentLocal = "local"

// EnvironmentProd is the environment variable for PROD
var EnvironmentProd = "prod"

// HeartBeatTime is the time in between heart beat notifications printed in the log files.
const HeartBeatTime = 2 * time.Minute

// Host is the machine this code is running on.
var Host, _ = os.Hostname()

// Live is a boolean that is true if the code is running in the production environment.
var Live = Environment == EnvironmentProd

// ProdLogDebug is used to force debug level logs in production.
// If true and Live then logs at Debug level else if false and Live then Info level.
var ProdLogDebug = GetBoolEnvVar("GO_DAEMONS_PROD_LOG_DEBUG", false)

// LogPath is the absolute directory where the HelloWorld log files are located at.
const LogPath = "/var/log/godaemons"

// PidPath is the absolute directory where the HelloWorld PID files are located at.
const PidPath = "/var/run/godaemons"

// UIURL is the URL to the HelloWorld Site (UIs).
const UIURL = ""

// TestLogName is the name to use when creating a test log file.
const TestLogName = "testing.log"

// *** Setters ***

// GetBoolEnvVar will get the environment variable for envVarName and attempt to cast it to a bool.
// If it fails or does not exist then uses the defaultValue.
func GetBoolEnvVar(envVarName string, defaultValue bool) bool {
	result := defaultValue
	value, exists := os.LookupEnv(envVarName)
	if exists {
		tempValue, err := strconv.ParseBool(strings.ToLower(value))
		if err == nil {
			result = tempValue
		}
	}
	return result
}

// getEnvironment is a setter for Environment.
func getEnvironment() string {
	// Default to true.
	environment := EnvironmentLocal
	value, exists := os.LookupEnv("GO_DAEMONS_ENVIRON")
	if exists {
		environment = strings.ToLower(value)
	}
	return environment
}
