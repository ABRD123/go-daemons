// Package utils contains common utility methods for the project.
package utils

import (
	"io"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// IsTimeUp implements a duration check, and returns true or false based on the results of the check. If a zero
// value is passed in for time value the assumption is that time is up.
//
// Note: The following conditions result in a return of true...
// t + d >= time.Now().UTC()
// t == 0
func IsTimeUp(d time.Duration, t time.Time) bool {
	r := true

	if !t.IsZero() {
		if t.Add(d).After(time.Now().UTC()) {
			r = false
		}
	}

	return r
}

// Close is a common method to be used when using defer to close resources.
func Close(c io.Closer, logger *log.Logger) {
	err := c.Close()

	if err != nil {
		// SSH pseudo-terminal connections can get closed early, so gobble up this error. Ugly way to check for
		// this error, but it's only defined internally in the Golang source.
		if strings.Contains(err.Error(), "use of closed network connection") {
			err = nil
		}
	}

	if logger != nil && err != nil {
		logger.Error(err.Error())
	}
}
