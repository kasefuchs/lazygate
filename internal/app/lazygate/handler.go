package lazygate

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
			if err := p.provider.Pause(srv); err != nil {
				p.log.Error(err, "failed to pause server")
			}
		}
	}
}

func (p *Plugin) onConnectionErrorEvent(event *proxy.ConnectionErrorEvent) {
	srv := event.Server()
	plr := event.Player()

	p.log.Info("resuming server", "name", srv.ServerInfo().Name())
	if err := p.provider.Resume(srv); err != nil {
		plr.Disconnect(resumeFailed)
		p.log.Error(err, "failed to resume server")

		return
	}

	plr.Disconnect(resumeInProgress)
}

func (p *Plugin) onPlayerChooseInitialServerEvent(event *proxy.PlayerChooseInitialServerEvent) {
	plr := event.Player()

	if event.InitialServer() == nil {
		p.log.Info("resuming first server")
		if err := p.provider.ResumeAny(); err != nil {
			plr.Disconnect(resumeFailed)
			p.log.Error(err, "failed to resume server")

			return
		}

		plr.Disconnect(resumeInProgress)
	}
}
