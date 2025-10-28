package plugin

import (
	"context"
	"fmt"
	"math"

	"github.com/go-logr/logr"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/queue"
	"github.com/kasefuchs/lazygate/pkg/registry"
	"github.com/kasefuchs/lazygate/pkg/scheduler"
	"github.com/kasefuchs/lazygate/pkg/utils"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

const (
	Name    = "lazygate"        // Name represents plugin name.
	logName = "lazygate.plugin" // Logger name to log with.
)

// Plugin is the LazyGate Gate plugin.
type Plugin struct {
	ctx       context.Context      // Plugin context.
	log       logr.Logger          // Plugin logger.
	proxy     *proxy.Proxy         // Gate proxy instance.
	queues    *queue.Repository    // Plugin queues repository.
	config    *Config              // Plugin configuration.
	options   *Options             // Plugin options.
	registry  *registry.Registry   // Plugin registry.
	provider  provider.Provider    // Allocation provider.
	scheduler *scheduler.Scheduler // Plugin server scheduler.
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

// initConfig loads plugin config.
func (p *Plugin) initConfig() error {
	cfg, err := utils.ParseEnv(DefaultConfig(), utils.ChildEnvPrefix("plugin"))
	if err != nil {
		return err
	}

	p.config = cfg.(*Config)
	return nil
}

// initProvider initializes server provider.
func (p *Plugin) initProvider() error {
	for _, prv := range p.options.Providers {
		if prv.Name() != p.config.Provider {
			continue
		}

		cfg, err := utils.ParseEnv(prv.DefaultConfig(), utils.ChildEnvPrefix("provider", prv.Name()))
		if err != nil {
			return err
		}

		p.provider = prv
		return p.provider.Init(&provider.InitOptions{
			Ctx:    p.ctx,
			Config: cfg,
		})
	}

	return fmt.Errorf("no such provider %s", p.config.Provider)
}

// initRegistry initializes new registry.
func (p *Plugin) initRegistry() error {
	p.registry = registry.NewRegistry(p.proxy, p.provider)
	p.registry.Refresh(p.config.Namespace)

	return nil
}

// initScheduler initializes server scheduler.
func (p *Plugin) initScheduler() error {
	p.scheduler = scheduler.NewScheduler(p.ctx, p.registry, p.provider)

	return p.scheduler.Init()
}

// initQueues initializes player queues.
func (p *Plugin) initQueues() error {
	opts := &queue.InitOptions{
		Proxy: p.proxy,
	}

	p.queues = queue.NewRepository()
	for _, q := range p.options.Queues {
		if err := q.Init(opts); err != nil {
			return err
		}

		p.queues.Push(q)
	}

	return nil
}

// initHandlers subscribes event handlers.
func (p *Plugin) initHandlers() error {
	eventMgr := p.proxy.Event()

	event.Subscribe(eventMgr, math.MaxInt, p.onDisconnectEvent)
	event.Subscribe(eventMgr, math.MaxInt, p.onServerPreConnectEvent)
	event.Subscribe(eventMgr, math.MaxInt, p.onServerRegistrationEvent)

	return nil
}

// Init initializes plugin functionality.
func (p *Plugin) Init() error {
	p.log = logr.FromContextOrDiscard(p.ctx).WithName(logName)

	if err := p.initConfig(); err != nil {
		return err
	}
	if err := p.initProvider(); err != nil {
		return err
	}
	if err := p.initRegistry(); err != nil {
		return err
	}
	if err := p.initScheduler(); err != nil {
		return err
	}
	if err := p.initQueues(); err != nil {
		return err
	}
	if err := p.initHandlers(); err != nil {
		return err
	}

	return nil
}
