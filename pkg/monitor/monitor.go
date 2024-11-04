package monitor

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/utils"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Monitor logger name.
const logName = "lazygate.monitor"

// Monitor checks servers.
type Monitor struct {
	ctx      context.Context   // Plugin context.
	log      logr.Logger       // Monitor logger.
	proxy    *proxy.Proxy      // Plugin proxy.
	provider provider.Provider // Plugin provider.
	registry *Registry         // Monitor entry registry.
}

// NewMonitor creates new instance of Monitor.
func NewMonitor(ctx context.Context, proxy *proxy.Proxy, provider provider.Provider) *Monitor {
	return &Monitor{
		ctx:      ctx,
		proxy:    proxy,
		provider: provider,
	}
}

// initRegistry initializes new server entry registry.
func (m *Monitor) initRegistry() error {
	m.registry = NewRegistry()

	return m.refreshRegistry()
}

// initTicker starts tickers.
func (m *Monitor) initTicker() error {
	go utils.Tick(m.ctx, 10*time.Second, m.allocationStopTicker)

	return nil
}

// Refreshes registry by adding new entries.
func (m *Monitor) refreshRegistry() error {
	m.registry.Clear()

	for _, srv := range m.proxy.Servers() {
		alloc, err := m.provider.AllocationGet(srv)
		if err != nil {
			continue
		}

		ent := NewEntry(srv, alloc)
		m.registry.EntryRegister(ent)
	}

	return nil
}

// Init initializes monitor.
func (m *Monitor) Init() error {
	m.log = logr.FromContextOrDiscard(m.ctx).WithName(logName)

	if err := m.initRegistry(); err != nil {
		return err
	}
	if err := m.initTicker(); err != nil {
		return err
	}

	m.log.Info("initialized")
	return nil
}

// Registry returns registry associated with this Monitor.
func (m *Monitor) Registry() *Registry {
	return m.registry
}
