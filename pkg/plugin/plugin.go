package plugin

import (
	"context"
	"math"

	"github.com/go-logr/logr"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/registry"
	"github.com/kasefuchs/lazygate/pkg/scheduler"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

const (
	Name    = "lazygate"        // Name represents plugin name.
	logName = "lazygate.plugin" // Logger name to log with.
)

// Plugin is the LazyGate Gate plugin.
type Plugin struct {
	ctx      context.Context      // Plugin context.
	log      logr.Logger          // Plugin logger.
	proxy    *proxy.Proxy         // Gate proxy instance.
	monitor  *scheduler.Scheduler // Plugin server monitor.
	options  *Options             // Plugin options.
	registry *registry.Registry   // Plugin registry.
	provider provider.Provider    // Allocation provider.
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
		Ctx: p.ctx,
	}

	return p.provider.Init(opt)
}

// initRegistry initializes new registry.
func (p *Plugin) initRegistry() error {
	p.registry = registry.NewRegistry(p.proxy, p.provider)
	p.registry.Refresh()

	return nil
}

// initMonitor initializes server monitor.
func (p *Plugin) initMonitor() error {
	p.monitor = scheduler.NewScheduler(p.ctx, p.registry, p.provider)

	return p.monitor.Init()
}

// initHandlers subscribes event handlers.
func (p *Plugin) initHandlers() error {
	eventMgr := p.proxy.Event()

	event.Subscribe(eventMgr, math.MaxInt, p.onDisconnectEvent)
	event.Subscribe(eventMgr, math.MaxInt, p.onServerPreConnectEvent)

	return nil
}

// Init initializes plugin functionality.
func (p *Plugin) Init() error {
	p.log = logr.FromContextOrDiscard(p.ctx).WithName(logName)

	if err := p.initProvider(); err != nil {
		return err
	}
	if err := p.initRegistry(); err != nil {
		return err
	}
	if err := p.initMonitor(); err != nil {
		return err
	}
	if err := p.initHandlers(); err != nil {
		return err
	}

	return nil
}
