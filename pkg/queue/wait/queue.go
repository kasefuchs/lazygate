package wait

import (
	"context"
	"time"

	"codeberg.org/kasefuchs/lazygate/pkg/queue"
	"codeberg.org/kasefuchs/lazygate/pkg/utils"
	"github.com/traefik/paerser/types"
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
	defer cancel()

	return utils.WaitUntil(ctx, time.Duration(tcfg.PingInterval), func(ctx context.Context) bool {
		return ticket.Entry.Ping(ctx, pcfg)
	})
}
