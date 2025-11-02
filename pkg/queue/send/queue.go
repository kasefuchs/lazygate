package send

import (
	"context"
	"time"

	"github.com/kasefuchs/lazygate/pkg/queue"
	"github.com/kasefuchs/lazygate/pkg/utils"
	"github.com/traefik/paerser/types"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var _ queue.Queue = (*Queue)(nil)

const name = "send"

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
		To:           "limbo",
		Timeout:      types.Duration(2 * time.Minute),
		PingInterval: types.Duration(5 * time.Second),
	}
}

func (q *Queue) Enter(ticket *queue.Ticket) bool {
	pcfg := q.proxy.Config()
	pctx := ticket.Player.Context()
	tcfg := ticket.Config.(*TicketConfig)

	ctx, cancel := context.WithTimeout(pctx, time.Duration(tcfg.Timeout))
	defer cancel()

	if to := q.proxy.Server(tcfg.To); to != nil {
		go ticket.Player.CreateConnectionRequest(to).ConnectWithIndication(ctx)
	} else {
		return false
	}

	return utils.WaitUntil(ctx, time.Duration(tcfg.PingInterval), func(ctx context.Context) bool {
		return ticket.Entry.Ping(ctx, pcfg)
	})
}
