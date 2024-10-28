package plugin

import (
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var (
	resumeInProgress = &component.Text{
		Content: "Server is currently starting up!",
	}
	resumeFailed = &component.Text{
		Content: "Server is unable to resume, please try again later.",
	}
)

func (p *Plugin) onDisconnectEvent(event *proxy.DisconnectEvent) {
	if conn := event.Player().CurrentServer(); conn != nil {
		srv := conn.Server()

		if srv.Players().Len() == 0 {
			p.log.Info("pausing server", "name", srv.ServerInfo().Name())
			if err := p.provider.Pause(&srv); err != nil {
				p.log.Error(err, "failed to pause server")
			}
		}
	}
}

func (p *Plugin) onPlayerChooseInitialServerEvent(event *proxy.PlayerChooseInitialServerEvent) {
	srv := event.InitialServer()
	plr := event.Player()

	if srv != nil {
		req := plr.CreateConnectionRequest(srv)
		if _, err := req.Connect(plr.Context()); err == nil {
			return
		}
	}

	p.log.Info("resuming server")
	if err := p.provider.Resume(&srv); err != nil {
		plr.Disconnect(resumeFailed)
		p.log.Error(err, "failed to resume server")

		return
	}
	plr.Disconnect(resumeInProgress)
}
