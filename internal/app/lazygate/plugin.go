package lazygate

import (
	"context"
	"math"

	"github.com/go-logr/logr"
	"github.com/kasefuchs/lazygate/internal/pkg/config/static"
	"github.com/kasefuchs/lazygate/internal/pkg/provider"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Plugin is the LazyGate Gate plugin.
type Plugin struct {
	ctx      context.Context   // Plugin context.
	log      logr.Logger       // Plugin logger.
	proxy    *proxy.Proxy      // Gate proxy instance.
	config   *static.Config    // Static config of plugin.
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

// loadConfig loads static config.
func (p *Plugin) loadConfig() error {
	var err error
	p.config, err = static.ParseEnv()
	if err != nil {
		return err
	}

	p.log.Info("loaded plugin configuration")
	return nil
}

// initProvider initializes server provider.
func (p *Plugin) initProvider() error {
	var err error
	p.provider, err = provider.NewProvider(p.config.Provider)
	if err != nil {
		return err
	}

	if err := p.provider.Init(); err != nil {
		return err
	}

	p.log.Info("initialized provider", "name", p.config.Provider)
	return nil
}

// initHandlers binds events handlers
func (p *Plugin) initHandlers() error {
	p.eventMgr = p.proxy.Event()

	event.Subscribe(p.eventMgr, math.MaxInt, p.onDisconnectEvent)
	event.Subscribe(p.eventMgr, math.MaxInt, p.onPlayerChooseInitialServerEvent)

	p.log.Info("subscribed plugin handlers")
	return nil
}

// Init initializes plugin functionality.
func (p *Plugin) Init() error {
	p.log = logr.FromContextOrDiscard(p.ctx)

	if err := p.loadConfig(); err != nil {
		return err
	}
	if err := p.initProvider(); err != nil {
		return err
	}
	if err := p.initHandlers(); err != nil {
		return err
	}

	return nil
}
