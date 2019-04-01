package context

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// AppContext is the collection of data items that is needed by the daemon, and implementation of Context interface.
type AppContext struct {
	// Logger is the reference to the daemon log file. This logger should be used for writing out daemon level events.
	// 	NOTE: For operational events associated with FlexUp and FlexDown *DO NOT* write them to the daemon level logger.
	Logger *log.Logger

	// Heartbeat is a timer used to intermittently log a message to the daemon log indicating that the daemon is
	// actively processing.
	Heartbeat time.Time
}

// GetHeartBeat gets the value of the heartbeat timer.
func (ctx *AppContext) GetHeartBeat() time.Time {
	return ctx.Heartbeat
}

// GetLogger gets an instance of log.Logger
func (ctx *AppContext) GetLogger() *log.Logger {
	return ctx.Logger
}

// SetHeartBeat sets the value of the heartbeat timer based on the value passed in.
func (ctx *AppContext) SetHeartBeat(t time.Time) {
	ctx.Heartbeat = t
}
