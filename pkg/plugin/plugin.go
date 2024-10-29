package plugin

import (
	"context"
	"math"

	"github.com/go-logr/logr"
	"github.com/kasefuchs/lazygate/pkg/provider"
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
	options  *Options          // Plugin options.
	eventMgr event.Manager     // Plugin proxy's event manager.
	provider provider.Provider // Runner provider.
}

// NewPlugin creates new instance of plugin.
func NewPlugin(ctx context.Context, proxy *proxy.Proxy, options ...*Options) *Plugin {
	opts := DefaultOptions()
	if len(options) > 0 {
		opts = options[0]
	}

	return &Plugin{
		ctx:     ctx,
		proxy:   proxy,
		options: opts,
	}
}

// NewProxyPlugin creates new instance of Gate Proxy plugin.
func NewProxyPlugin(options ...*Options) proxy.Plugin {
	return proxy.Plugin{
		Name: Name,
		Init: func(ctx context.Context, proxy *proxy.Proxy) error {
			return NewPlugin(ctx, proxy, options...).Init()
		},
	}
}

// initProvider initializes server provider.
func (p *Plugin) initProvider() error {
	var err error
	p.provider, err = p.options.ProviderSelector()
	if err != nil {
		return err
	}

	opt := &provider.InitOptions{
		Log: p.log.WithName("provider"),
	}

	return p.provider.Init(opt)
}

// bindHandlers binds event handlers.
func (p *Plugin) bindHandlers() error {
	p.eventMgr = p.proxy.Event()

	event.Subscribe(p.eventMgr, math.MaxInt, p.onDisconnectEvent)
	event.Subscribe(p.eventMgr, math.MaxInt, p.onPlayerChooseInitialServerEvent)

	p.log.Info("subscribed handlers")
	return nil
}

// Init initializes plugin functionality.
func (p *Plugin) Init() error {
	p.log = logr.FromContextOrDiscard(p.ctx).WithName(Name)

	if err := p.initProvider(); err != nil {
		return err
	}
	if err := p.bindHandlers(); err != nil {
		return err
	}

	return nil
}
