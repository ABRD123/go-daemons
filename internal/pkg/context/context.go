// Package context is contains different contexts used in daemons and processors
package context

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Context interface exports methods for daemons to implement
type Context interface {
	// GetHeartBeat gets the heartbeat value stored in context.
	GetHeartBeat() time.Time
	// GetLogger returns an instance of logger
	GetLogger() *log.Logger
	// SetHeartBeat sets the heartbeat value stored in the context.
	SetHeartBeat(t time.Time)
}
