package plugin

import (
	"context"
	"fmt"
	"math"

	"github.com/kasefuchs/lazygate/pkg/provider/docker"

	"github.com/go-logr/logr"
	"github.com/kasefuchs/lazygate/pkg/config/static"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/provider/nomad"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Name represents plugin name.
const Name = "lazygate"

// Plugin is the LazyGate Gate plugin.
type Plugin struct {
	ctx      context.Context   // Plugin context.
	log      logr.Logger       // Plugin logger.
	proxy    *proxy.Proxy      // Gate proxy instance.
	config   *static.Config    // Static config of plugin.
	eventMgr event.Manager     // Plugin proxy's event manager.
	provider provider.Provider // Runner provider.
}

// NewPlugin creates new instance of plugin.
func NewPlugin(ctx context.Context, proxy *proxy.Proxy) *Plugin {
	return &Plugin{
		ctx:   ctx,
		proxy: proxy,
	}
}

// NewProxyPlugin creates new instance of Gate Proxy plugin.
func NewProxyPlugin() proxy.Plugin {
	return proxy.Plugin{
		Name: Name,
		Init: func(ctx context.Context, proxy *proxy.Proxy) error {
			return NewPlugin(ctx, proxy).Init()
		},
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
	switch p.config.Provider {
	case "nomad":
		p.provider = &nomad.Provider{}
	case "docker":
		p.provider = &docker.Provider{}
	case "":
		return fmt.Errorf("empty provider")
	default:
		return fmt.Errorf("unknown provider: %s", p.config.Provider)
	}

	opt := provider.InitOptions{}
	if err := p.provider.Init(opt); err != nil {
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
	p.log = logr.FromContextOrDiscard(p.ctx).WithName(Name)

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
