package wait

import (
	"context"
	"time"

	"github.com/kasefuchs/lazygate/pkg/queue"
	"github.com/kasefuchs/lazygate/pkg/utils"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var _ queue.Queue = (*Queue)(nil)

const name = "wait"

type Queue struct {
	proxy *proxy.Proxy
}

func (q *Queue) Name() string {
	return name
}

func (q *Queue) Init(opts *queue.InitOptions) error {
	q.proxy = opts.Proxy

	return nil
}

func (q *Queue) Enter(ticket *queue.Ticket) bool {
	pcfg := q.proxy.Config()
	pctx := ticket.Player.Context()

	ctx, cancel := context.WithTimeout(pctx, time.Duration(ticket.Config.Queue.Wait.Timeout))
	utils.Tick(ctx, time.Duration(ticket.Config.Queue.Wait.PingInterval), func() {
		if ticket.Entry.Ping(ctx, pcfg) {
			cancel()
		}
	})

	return ticket.Entry.Ping(pctx, pcfg)
}
