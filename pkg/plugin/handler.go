package plugin

import (
	"time"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func (p *Plugin) onDisconnectEvent(event *proxy.DisconnectEvent) {
	if conn := event.Player().CurrentServer(); conn != nil {
		srv := conn.Server()
		if ent := p.monitor.Registry().EntryGet(srv); ent != nil {
			ent.UpdateLastActive()
		}
	}
}

func (p *Plugin) onServerPreConnectEvent(event *proxy.ServerPreConnectEvent) {
	plr := event.Player()
	srv := event.Server()

	ent := p.monitor.Registry().EntryGet(srv)
	if ent == nil || ent.Ping(plr.Context(), p.proxy.Config()) {
		return
	}

	cfg, err := ent.Allocation.Config()
	if err != nil {
		return
	}

	name := srv.ServerInfo().Name()
	p.log.Info("Starting server allocation", "server", name)

	ent.KeepOnlineFor(time.Duration(cfg.Time.MinimumOnline))

	if err := ent.Allocation.Start(); err != nil {
		p.log.Error(err, "Failed to start server allocation", "server", name)
		plr.Disconnect(cfg.DisconnectReasons.StartFailed.TextComponent())

		return
	}

	plr.Disconnect(cfg.DisconnectReasons.Starting.TextComponent())
}

func (p *Plugin) onServerPostConnectEvent(event *proxy.ServerPostConnectEvent) {
	srv := event.Player().CurrentServer().Server()
	if ent := p.monitor.Registry().EntryGet(srv); ent != nil {
		ent.UpdateLastActive()
	}
}
