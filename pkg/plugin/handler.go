package plugin

import (
	"time"

	"github.com/kasefuchs/lazygate/pkg/queue"
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

	cfg, err := ent.Allocation.Config()
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

	ticket := &queue.Ticket{
		Entry:  ent,
		Config: cfg,
		Player: plr,
	}

	for _, qn := range cfg.Queue.Try {
		q := p.queues.Get(qn)

		if q == nil || q.Enter(ticket) {
			return
		}
	}
}
