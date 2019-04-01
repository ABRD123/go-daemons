#golang daemon

Basic implementation of my own version of a golang daemon.

Uses low-level daemon implementation from [go-daemon](https://github.com/sevlyar/go-daemon).

Currently supports OSX and Linuxes. Tested on Mac OSX 10.11 and Ubuntu 16.04.

For example see daemon_example_test.go

## How the CLI works

*daemon* help
* Prints a usage message and exits immediately.

*daemon* start
* Attemtps to start the daemon, if it's not already running.

*daemon* stop
* Searches for a running daemon matching the supplied context info, and signals the daemon with SIGTERM.
* If SetTerminatorHandler() has supplied a valid function, this function will be invoked in the running daemon.
* If SetTerminatorHandler() hasn't been called, the default handler will immediately terminate the running daemon.

*daemon* restart
* Searches for a running daemon matching the supplied context info, and signals the daemon with SIGTERM, then...
* Attempts to start the daemon.

*daemon* reload
* Searches for a running daemon matching the supplied context info, and signals the daemen with SIGHUP.
* If a SetReloadHandler() has supplied a valid function, it will be invoked in the running daemon.
* If SetReloadHandler() has not been called, this is a NOOP.

*daemon* status 
* Searches for a running daemon matching the supplied context info, and displays the PID info.

*daemon* debug
* Runs the daemon in debug mode, as a foreground application. Bypasses all go-daemon functionality.

## How the daemon works

When the daemon CLI processer detects a 'start' or 'restart' (after all validations have passed) the parent process
invokes the current application as a *daemon* process. When the newly running process identifies itself as the *daemon*
invocation, rather than parsing the command line it performs a goroutine invocation of the *worker* function.

For more information on the low-level *daemon* invocation see [go-daemon](https://github.com/sevlyar/go-daemon).
