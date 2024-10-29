package plugin

import (
	"math/rand"

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
			if alloc, err := p.provider.AllocationGet(srv); err == nil {
				p.log.Info("stopping server", "name", srv.ServerInfo().Name())
				if err := alloc.Stop(); err != nil {
					p.log.Error(err, "failed to stop server")
				}
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

		if alloc, err := p.provider.AllocationGet(srv); err == nil {
			if err := alloc.Start(); err != nil {
				plr.Disconnect(resumeFailed)
				p.log.Error(err, "failed to resume server")

				return
			}

			plr.Disconnect(resumeInProgress)
		}

		return
	}

	allocs, err := p.provider.AllocationList()
	if err != nil || len(allocs) == 0 {
		return
	}

	alloc := allocs[rand.Intn(len(allocs))]
	if err := alloc.Start(); err != nil {
		plr.Disconnect(resumeFailed)
		p.log.Error(err, "failed to resume server")

		return
	}

	plr.Disconnect(resumeInProgress)
}
