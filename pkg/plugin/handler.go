package plugin

import (
	"math/rand"

	"go.minekube.com/gate/pkg/edition/java/proto/version"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/componentutil"
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
			cfg := alloc.Config()
			if err := alloc.Start(); err != nil {
				p.log.Error(err, "failed to resume server")
				if textComponent, err := componentutil.ParseTextComponent(version.MinimumVersion.Protocol, cfg.DisconnectReasons.ActionFailed); err == nil {
					plr.Disconnect(textComponent)
				}

				return
			}

			if textComponent, err := componentutil.ParseTextComponent(version.MinimumVersion.Protocol, cfg.DisconnectReasons.Starting); err == nil {
				plr.Disconnect(textComponent)
			}

			return
		}

		return
	}

	allocs, err := p.provider.AllocationList()
	if err != nil || len(allocs) == 0 {
		return
	}

	alloc := allocs[rand.Intn(len(allocs))]
	cfg := alloc.Config()
	if err := alloc.Start(); err != nil {
		p.log.Error(err, "failed to resume server")
		if textComponent, err := componentutil.ParseTextComponent(version.MinimumVersion.Protocol, cfg.DisconnectReasons.ActionFailed); err == nil {
			plr.Disconnect(textComponent)
		}

		return
	}

	if textComponent, err := componentutil.ParseTextComponent(version.MinimumVersion.Protocol, cfg.DisconnectReasons.Starting); err == nil {
		plr.Disconnect(textComponent)
	}
}
