package plugin

import (
	"time"

	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/queue"
	"github.com/kasefuchs/lazygate/pkg/utils"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func (p *Plugin) onDisconnectEvent(event *proxy.DisconnectEvent) {
	if conn := event.Player().CurrentServer(); conn != nil {
		srv := conn.Server()
		if ent := p.registry.EntryGet(srv); ent != nil {
			ent.UpdateLastActive()
		}
	}
}

func (p *Plugin) onServerPreConnectEvent(event *proxy.ServerPreConnectEvent) {
	plr := event.Player()
	srv := event.Server()
	ctx := plr.Context()

	ent := p.registry.EntryGet(srv)
	if ent == nil || ent.Ping(ctx, p.proxy.Config()) {
		return
	}

	cfg, err := provider.ParseAllocationConfig(ent.Allocation)
	if err != nil {
		return
	}

	name := srv.ServerInfo().Name()
	p.log.Info("starting server allocation", "server", name)
	if err := ent.Allocation.Start(); err != nil {
		p.log.Error(err, "failed to start server allocation", "server", name)

		return
	}

	ent.KeepOnlineFor(time.Duration(cfg.Time.MinimumOnline))

	for _, qn := range cfg.Queues {
		q := p.queues.Get(qn)
		if q == nil {
			continue
		}

		tcfg, err := ent.Allocation.ParseConfig(q.DefaultTicketConfig(), utils.ChildLabel("queue", qn))
		if err != nil {
			continue
		}

		ticket := &queue.Ticket{
			Entry:  ent,
			Config: tcfg,
			Player: plr,
		}

		if q.Enter(ticket) {
			return
		}
	}
}

func (p *Plugin) onServerRegistrationEvent(event *proxy.ServerRegisteredEvent) {
	srv := event.Server()
	p.registry.RegisterServer(srv, p.config.Namespace)
}
