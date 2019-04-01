package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	godaemon "github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"

	"github.com/go-daemons/configs"
	"github.com/go-daemons/configs/helloworldconfigs"
	"github.com/go-daemons/internal/apps/helloworld"
	"github.com/go-daemons/internal/pkg/context"
	"github.com/go-daemons/internal/pkg/logutil"
	"github.com/go-daemons/internal/pkg/utils"
	"github.com/go-daemons/pkg/daemon"
)

// daemonSignaling is used by worker() and terminator() to send signals to each other for a clean shutdown.
type daemonSignaling struct {
	shutdown    chan bool
	shutdownAck chan bool
}

// terminator is the signal handler called when a SIGTERM is sent to the daemon. This function will attempt
// to perform a clean shutdown.
func terminator(signals *daemonSignaling) func(_ os.Signal) error {
	return func(_ os.Signal) error {
		log.Info("daemon terminator() called...")
		// Signal the daemon worker that it needs to shutdown.
		signals.shutdown <- true
		// Wait for the daemon work to signal that it's in a safe place to shut-down.
		<-signals.shutdownAck
		return godaemon.ErrStop
	}
}

// worker is the actual daemon infinite loop itself.
func worker(signals *daemonSignaling) func() {
	return func() {
		// Setup logging
		logger := log.New()
		file := logutil.SetupLogging(logger, helloworldconfigs.LogName)
		defer utils.Close(file, logger)

		ctx := &context.AppContext{Logger: logger}

		logger.Info("- - - - - - - - - - - - - - -")
		logger.Infof("Daemon %v started", helloworldconfigs.AppName)

		logger.WithFields(log.Fields{
			"app_name":    helloworldconfigs.AppName,
			"pid_file":    helloworldconfigs.PidFile,
			"log_file":    helloworldconfigs.LogName,
			"working_dir": helloworldconfigs.WorkingDir}).Info("Starting daemon")

		term := false
		timer := time.Now().UTC()
		for {
			select {
			case <-signals.shutdown:
				// Check to see if we need to shutdown the worker
				term = true
			default:
				// Make sure we don't call the orchestration too quickly.
				// This also provides a minor pause on daemon start-up.
				if utils.IsTimeUp(helloworldconfigs.OrchestrationWaitTime, timer) {
					timer = time.Now().UTC()
					_ = helloworld.Daemon(ctx)
				}
			}

			if term {
				break
			}

			time.Sleep(1 * time.Second)
		}

		// Signal the terminator that it's safe to proceed with a shutdown.
		logger.Info("daemon worker() graceful shutdown")
		signals.shutdownAck <- true
	}
}

func main() {
	// Entry point for both the "parent" command line processing, and the daemon processing. The underlying
	// daemon package handles knowing which invocation is which.

	signals := &daemonSignaling{make(chan bool), make(chan bool)}
	ctx := &daemon.Context{}
	logFile := filepath.Join(configs.LogPath, helloworldconfigs.LogName)
	_ = ctx.New(helloworldconfigs.PidFile, logFile, helloworldconfigs.WorkingDir, helloworldconfigs.AppName)
	fmt.Printf(helloworldconfigs.WorkingDir)
	ctx.SetWorkerHandler(worker(signals))
	ctx.SetTerminatorHandler(terminator(signals))
	daemon.ProcessCommandLine(ctx)
}
