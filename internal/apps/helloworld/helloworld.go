package helloworld

import (
	"time"

	"github.com/go-daemons/configs"
	"github.com/go-daemons/internal/pkg/context"
	"github.com/go-daemons/internal/pkg/utils"
)

// Daemon is the function that performs the orchestration for the HelloWorld daemon. Call this function from the
// HelloWorld daemon framework.
func Daemon(ctx context.Context) error {
	if utils.IsTimeUp(configs.HeartBeatTime, ctx.GetHeartBeat()) {
		ctx.SetHeartBeat(time.Now().UTC())
		ctx.GetLogger().Info("Daemon orchestration heartbeat")
	}
	ctx.GetLogger().Info("****Hello world******")
	return nil
}
