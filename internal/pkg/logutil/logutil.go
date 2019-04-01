// Package logutil implements a library for interacting with log files.
package logutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"

	"github.com/go-daemons/configs"
)

// SetupLogging will setup a log file using logName and return a handle to the file. The logName will have the
// project log path prepended to it so you only need to send just the log name. ex. directord.log
//
// NOTE: The caller is required to close the returned file when finished with it.
//  Param logger: in/out, Pointer to a Logger object.
//  Param logName: in, log name to be used when opening the log file. ex. "name.log"
//  Returns pointer to a File.
func SetupLogging(logger *log.Logger, logName string) *os.File {
	fullLogName := filepath.Join(configs.LogPath, logName)
	file, err := os.OpenFile(fullLogName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		// Cannot open log file.
		logger.Out = os.Stdout
		logger.Panic(err)
	} else {
		logger.Out = file
		filenameHook := filename.NewHook()
		filenameHook.Field = "source"
		filenameHook.Formatter = func(file, function string, line int) string {
			index := strings.LastIndex(function, "/")
			if index > -1 {
				function = function[index+1:]
			}
			return fmt.Sprintf("%s::%s::%d", file, function, line)
		}
		logger.AddHook(filenameHook)
		if !configs.Live || configs.ProdLogDebug {
			// Debug level if not live environment or bypassed.
			logger.SetLevel(log.DebugLevel)
		}
	}
	return file
}
