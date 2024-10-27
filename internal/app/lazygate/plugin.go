package lazygate

import (
	"context"
	"math"

	"github.com/go-logr/logr"
	"github.com/kasefuchs/lazygate/internal/pkg/provider"
	"github.com/kasefuchs/lazygate/internal/pkg/provider/nomad"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Plugin is the LazyGate Gate plugin.
type Plugin struct {
	ctx      context.Context   // Plugin context.
	log      logr.Logger       // Plugin logger.
	proxy    *proxy.Proxy      // Gate proxy instance.
	eventMgr event.Manager     // Plugin proxy's event manager.
	provider provider.Provider // Runner provider.
}

// NewPlugin creates new instance of LazyGate plugin.
func NewPlugin(ctx context.Context, proxy *proxy.Proxy) *Plugin {
	return &Plugin{
		ctx:   ctx,
		proxy: proxy,
	}
}

// initProvider initializes server provider.
func (p *Plugin) initProvider() error {
	p.provider = &nomad.Provider{}
	if err := p.provider.Init(); err != nil {
		return err
	}

	p.log.Info("initialized provider", "name", p.provider.Name())
	return nil
}

// initHandlers binds events handlers
func (p *Plugin) initHandlers() error {
	p.eventMgr = p.proxy.Event()

	event.Subscribe(p.eventMgr, math.MaxInt, p.onDisconnectEvent)
	event.Subscribe(p.eventMgr, math.MaxInt, p.onConnectionErrorEvent)

	p.log.Info("subscribed plugin handlers")
	return nil
}

// Init initializes plugin functionality.
func (p *Plugin) Init() error {
	p.log = logr.FromContextOrDiscard(p.ctx)

	if err := p.initProvider(); err != nil {
		return err
	}
	if err := p.initHandlers(); err != nil {
		return err
	}

	return nil
}
