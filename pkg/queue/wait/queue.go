package wait

import (
	"context"
	"time"

	"github.com/traefik/paerser/types"

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

func (q *Queue) DefaultTicketConfig() interface{} {
	return &TicketConfig{
		Timeout:      types.Duration(25 * time.Second),
		PingInterval: types.Duration(3 * time.Second),
	}
}

func (q *Queue) Enter(ticket *queue.Ticket) bool {
	pcfg := q.proxy.Config()
	pctx := ticket.Player.Context()
	tcfg := ticket.Config.(*TicketConfig)

	ctx, cancel := context.WithTimeout(pctx, time.Duration(tcfg.Timeout))
	utils.Tick(ctx, time.Duration(tcfg.PingInterval), func() {
		if ticket.Entry.Ping(ctx, pcfg) {
			cancel()
		}
	})

	return ticket.Entry.Ping(pctx, pcfg)
}
