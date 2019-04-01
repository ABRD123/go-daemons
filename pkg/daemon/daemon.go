// Package daemon implements a my own version of a golang daemon.
//
// This package builds upon the functionality of https://github.com/sevlyar/go-daemon.
//
// The documentation for this package is at https://godoc.org/github.com/sevlyar/go-daemon.
//
// The "trick" with this daemon implementation, and the underlying go-daemon implementation, is that a single
// Go executable acts as both the parent command line processer, and as the daemon itself. Starting the process from
// the command line (or any other script) invokes the parent command line processing. Whereas passing either *start*
// or *restart* to the parent causes it to re-invoke the Go executable, but this time running in daemon mode.
package daemon

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	godaemon "github.com/sevlyar/go-daemon"
)

const usage = `%[1]v is a daemon

Usage:

	%[1]v <command>

The commands are:

	help       Print this usage information
	start      start as a daemon
	stop       find the running daemon, and shut it down
	restart    perform a stop + start operation
	reload     find the running daemon, and have it reload it's configuration
	status     find the running daemon, and print out it's PID
	debug      start, not as a daemon, but as a foreground process for debugging purposes`

// HandlerFunc is the function signature for daemon run-time functions that implement the signal handling for various
// events. Currently used by SetReloadHandler() and  SetTerminatorHandler().
type HandlerFunc func(sig os.Signal) (err error)

// WorkerFunc is the function signature for the daemon worker function. Used by SetWorkerHandler.
type WorkerFunc func()

// Context is the persistent data structure used to store internal daemon data in between function/method calls.
type Context struct {
	goctx      *godaemon.Context
	reloader   HandlerFunc
	terminator HandlerFunc
	worker     func()
}

// New allocates a Context structure using the supplied parameters. Once allocated this context maintains
// the persistant data storage for this package. It is up to the daemon implementor to decide where the PID file,
// log file and working directory are located. The appname parameter is the value that is displayed when someone runs
// 'ps -eaf' on the same system where the daemon is running.
func (ctx *Context) New(pidfile string, logfile string, workingdir string, appname string) error {
	ctx.goctx = &godaemon.Context{
		PidFileName: pidfile,
		PidFilePerm: 0644,
		LogFileName: logfile,
		LogFilePerm: 0640,
		WorkDir:     workingdir,
		Umask:       027,
		Args:        []string{appname},
	}

	return nil
}

// SetWorkerHandler is used to set the implementors worker function that is the basis of the daemon's run-time.
// This worker function can use any features of the go language (e.g. goroutines). Typically the worker function
// will enter an endless loop to perform it's processing.
func (ctx *Context) SetWorkerHandler(f WorkerFunc) {
	ctx.worker = f
}

// SetReloadHandler is an optional method used to set the function called when a "<daemon> reload" CLI operation is
// called. The implementation of this function is provided by the implementor of the daemon. If not provided
// "<daemon> reload" is a noop.
func (ctx *Context) SetReloadHandler(f HandlerFunc) {
	ctx.reloader = f
}

// SetTerminatorHandler is an optional method used to set the function called when a "<daemon> stop" CLI operation is
// called. The implementation of this function is provided by the implementor of the daemon. If not provided the
// default stop handler (from go-daemon) will forcefully terminate the running daemon.
func (ctx *Context) SetTerminatorHandler(f HandlerFunc) {
	ctx.terminator = f
}

// debugDaemon is the function used to run as a foreground process instead of a daemon process. Since there is no
// underlying daemon implementation running the application will behave like a normal application. Ctrl-C or the
// debugger stop command is required to terminate the process.
func debugDaemon(ctx *Context) {
	// Execute the daemon implementors code, but *DO NOT* execute as a goroutine. Executing as a goroutine will
	// allow the main thread of execution to continue on, and will result in the application exiting.
	ctx.worker()
}

// reloadDaemonConfig is the function used by the parent process to signal the running daemon to call the
// handler set by SetReloadHandler().
func reloadDaemonConfig(ctx *Context, proc *os.Process) error {
	if ctx.reloader != nil {
		fmt.Println("Sending reload signal to daemon with PID,", proc.Pid)
		_ = proc.Signal(syscall.SIGHUP)
	} else {
		fmt.Println("Daemon does not support reload functionality.")
	}
	return nil
}

// displayStatus is the function used by the parent process to display information about the running daemon.
func displayStatus(ctx *Context, proc *os.Process) {
	fmt.Printf("%v running as a daemon with PID %v\n", ctx.goctx.Args[0], proc.Pid)
}

// printUsage is the function used by the parent process to display the command line usage information.
func printUsage(ctx *Context) {
	fmt.Printf(usage, ctx.goctx.Args[0])
}

// runDaemon is the function used by the daemon process to register the reload and terminator handlers, and then
// execute the daemon implementors specific daemon processing code. The daemon implementors worker function is
// executed as a goroutine, and this function will block on the underlying ServeSignals() function; this is what
// allows the daemon to receive the SIGTERM and SIGHUP signals from the parent process.
func runDaemon(ctx *Context) {
	if ctx.terminator != nil {
		godaemon.SetSigHandler(godaemon.SignalHandlerFunc(ctx.terminator), syscall.SIGTERM)
	}
	if ctx.reloader != nil {
		godaemon.SetSigHandler(godaemon.SignalHandlerFunc(ctx.reloader), syscall.SIGHUP)
	}

	// Execute the daemon implementors code as a goroutine.

	go ctx.worker()

	// Block here and process any incoming signals from the parent process. Currently only SIGTERM and SIGHUP
	// are supported.

	err := godaemon.ServeSignals()
	if err != nil {
		log.Printf("Daemon %v failed with error, %v\n", ctx.goctx.Args[0], err)
		return
	}
	log.Printf("Daemon %v terminated normally", ctx.goctx.Args[0])
}

// startDaemon is the function used by the parent process to start up a daemon. Note, the parent and the daemon are
// the same exact go binary, they just execute a different path depending on which one they are.
func startDaemon(ctx *Context) error {
	_, err := ctx.goctx.Reborn()
	return err
}

// stopDaemon is the function used by the parent process to signal a running daemon process to perform a shutdown.
// This function will monitor the running daemon and if it doesn't shut down in a timely manner first try sending a
// SIGINT, and if that doesn't work send a SIGKILL.
func stopDaemon(_ *Context, proc *os.Process) error {
	fmt.Printf("Shutting down daemon with PID, %v", proc.Pid)
	_ = proc.Signal(syscall.SIGTERM)

	// Monitor the process for 10seconds, waiting for it to self-terminate.

	terminated := false
	monitorStart := time.Now().UTC()
	for time.Now().UTC().Before(monitorStart.Add(10 * time.Second)) {
		fmt.Printf(".")
		time.Sleep(1 * time.Second)
		err := proc.Signal(syscall.Signal(0))
		if err != nil {
			terminated = true
			break
		}
	}
	fmt.Println("")

	// If the process never self terminates, then we'll try a SIGINT.

	if !terminated {
		fmt.Printf("Sending a SIGINT to the daemon with PID, %v", proc.Pid)
		_ = proc.Signal(syscall.SIGINT)

		// Monitor the process for 10seconds, waiting for the SIGINT to take effect.

		terminated = false
		monitorStart = time.Now().UTC()
		for time.Now().UTC().Before(monitorStart.Add(10 * time.Second)) {
			fmt.Printf(".")
			time.Sleep(1 * time.Second)
			err := proc.Signal(syscall.Signal(0))
			if err != nil {
				terminated = true
				break
			}
		}
		fmt.Println("")

		// If the SIGINT doesn't work, send a SIGKILL.

		if !terminated {
			fmt.Println("Sending a SIGKILL to the daemon with PID,", proc.Pid)
			_ = proc.Signal(syscall.SIGKILL)
		}
	}

	return nil
}

// restartDaemon is the function used by the parent process to stop and (re)start a running daemon.
func restartDaemon(ctx *Context, proc *os.Process) error {
	_ = stopDaemon(ctx, proc)

	// Actively monitor the old daemon process to determine when it really goes away.
	for {
		time.Sleep(1 * time.Second)
		err := proc.Signal(syscall.Signal(0))
		if err != nil {
			break
		}
	}

	return startDaemon(ctx)
}

// ProcessCommandLine is a dual-purpose function. When the daemon is starting up calling ProcessCommandLine will
// setup the daemons run-time environment then call the implementors worker function to commence the actual daemon
// processing. The other purpose is to act as the parent process to manage daeomon operations like start, stop,
// restart, reload and status.
//noinspection GoUnusedExportedFunction
func ProcessCommandLine(ctx *Context) {
	if godaemon.WasReborn() {
		// Daemon processing. No exit codes are returned.

		proc, err := ctx.goctx.Reborn()
		if err != nil {
			log.Fatal("Fatal error during daemon start-up, aborting")
		}
		if proc == nil {
			// This is the daemon processing starting up. We need to run the code for normal daemon processing.
			//noinspection GoUnhandledErrorResult
			defer ctx.goctx.Release()
			runDaemon(ctx)
		}
	} else {
		// Command line application processing. This will *always* send an OS exit code back to the caller, for easier
		// bash shell script integration.

		retCode := 0
		defer func() { os.Exit(retCode) }()

		proc, err := ctx.goctx.Search()
		if err == nil && proc != nil {
			// In *nix need an additional test to see if the process actually exists...
			err = proc.Signal(syscall.Signal(0))
			if err != nil {
				proc = nil
			}
		}

		if len(os.Args) != 2 {
			printUsage(ctx)
			retCode = 1
			return
		}

		switch os.Args[1] {
		case "help":
			printUsage(ctx)
		case "start":
			if err != nil && proc == nil {
				if ctx.worker != nil {
					err = startDaemon(ctx)
					if err == nil {
						fmt.Println("Daemon successfully started")
						return
					}
					fmt.Printf("Cannot start the daemon %v, err = %v\n", os.Args[0], err)
				} else {
					fmt.Println("Cannot start daemon, no worker function/implementation supplied")
				}
			} else if proc != nil {
				fmt.Printf("There is already a daemon running for %v, with process ID %v\n", os.Args[0], proc.Pid)
			} else {
				fmt.Println("Should never make it here.")
			}
		case "stop":
			if err == nil && proc != nil {
				err = stopDaemon(ctx, proc)
				if err == nil {
					return
				}
				fmt.Printf("Cannot stop the daemon %v, err = %v\n", os.Args[0], err)
			} else {
				fmt.Println("Cannot find a daemon running for", os.Args[0])
			}
		case "restart":
			if err == nil && proc != nil {
				err = restartDaemon(ctx, proc)
				if err == nil {
					fmt.Println("Daemon successfully restarted")
					return
				}
				fmt.Printf("Cannot restart the daemon %v, err = %v\n", os.Args[0], err)
			} else {
				// Daemon is not currently running, just process as a start.
				err = startDaemon(ctx)
				if err == nil {
					fmt.Println("Daemon successfully started")
					return
				}
				fmt.Printf("Cannot start the daemon %v, err = %v\n", os.Args[0], err)
			}
		case "reload":
			if err == nil && proc != nil {
				err = reloadDaemonConfig(ctx, proc)
				if err == nil {
					return
				}
				fmt.Printf("Cannot reload configs for daemon %v, err = %v\n", os.Args[0], err)
			} else {
				fmt.Println("Cannot find a daemon running for", os.Args[0])
			}
		case "status":
			if err != nil || proc == nil {
				fmt.Println("Cannot find a daemon running for", os.Args[0])
			} else {
				displayStatus(ctx, proc)
				return
			}
		case "debug":
			if err == nil && proc != nil {
				fmt.Printf("Cannot run in DEBUG mode, the daemon %v is running with PID %v\n", os.Args[0], proc.Pid)
			} else {
				debugDaemon(ctx)
			}
		default:
			printUsage(ctx)
		}

		// Getting here indicates that a problem occurred above, so make sure the return code is set properly.
		retCode = 1
	}

	return
}
