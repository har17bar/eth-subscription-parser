package job

import (
	"context"
	"time"
	"txparser/internal/service"
)

const cronInterval = 15 * time.Second

type Eth struct {
	ethService service.Eth
}

func NewEth(ethService service.Eth) *Eth {
	return &Eth{
		ethService: ethService,
	}
}

func (e *Eth) Start(ctx context.Context) error {
	go func() {

		err := e.ethService.Process(ctx)
		if err != nil {
			return
		}
		ticker := time.NewTicker(cronInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := e.ethService.Process(ctx)
				if err != nil {
					return
				}
			}
		}
	}()

	return nil
}
