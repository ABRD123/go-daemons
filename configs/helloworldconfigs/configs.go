// Package helloworldconfigs contains configs specific to the helloworld daemon.
package helloworldconfigs

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-daemons/configs"
)

// *****************************
// *** Environment ***
// *****************************

// AppName is the name of the daemon that gets displayed when 'ps -eaf' or similar gets run from the *nix command
// line.
const AppName = "helloworld"

// LogName is the filename of the daemon log file that is written to when the HelloWorld daemon is running.
var LogName = fmt.Sprintf("%s.log", AppName)

// PidFile is the absolute pathname/filename of the PID file that is used when the HelloWorld daemon is running.
var PidFile = filepath.Join(configs.PidPath, fmt.Sprintf("%s.pid", AppName))

// WorkingDir is the location where the HelloWorld daemon will be running from. In this case were simply using the
// current directory where the daemon was run from as the working directory.
var WorkingDir = filepath.Join("/Users/abrd/go/src/github.com/go-daemons", fmt.Sprintf("/cmd/%s", AppName))

// *****************************
// *** Timers ***
// *****************************

// CreateCheckTime is the duration to delay between successive checks for successful VM creation.
var CreateCheckTime = 60 * time.Second

// OrchestrationWaitTime is the mimimum delay between successive executions of the orchestration layer. The intent
// is to make sure the daemon is not using excessive CPU spinning and doing nothing.
var OrchestrationWaitTime = 30 * time.Second

// ************************
