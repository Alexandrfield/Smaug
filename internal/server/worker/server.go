package worker

import (
	"context"

	"github.com/Alexandrfield/Smaug/internal/common"
)

func ServerLoop(ctx context.Context, done chan struct{}, logger common.Logger) {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("Stop ServerLoop")
			close(done)
			return
		default:
			continue
		}
	}
}
