package scheduler

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/registry"
	"github.com/kasefuchs/lazygate/pkg/utils"
)

// Scheduler logger name.
const logName = "lazygate.scheduler"

// Scheduler performs periodic operations.
type Scheduler struct {
	ctx      context.Context    // Plugin context.
	log      logr.Logger        // Scheduler logger.
	provider provider.Provider  // Plugin provider.
	registry *registry.Registry // Scheduler entry registry.
}

// NewScheduler creates new instance of Scheduler.
func NewScheduler(ctx context.Context, registry *registry.Registry, provider provider.Provider) *Scheduler {
	return &Scheduler{
		ctx:      ctx,
		registry: registry,
		provider: provider,
	}
}

// startTasks starts tasks.
func (m *Scheduler) startTasks() error {
	go utils.Tick(m.ctx, 10*time.Second, m.stopStaleAllocations)

	return nil
}

// Init initializes scheduler.
func (m *Scheduler) Init() error {
	m.log = logr.FromContextOrDiscard(m.ctx).WithName(logName)

	if err := m.startTasks(); err != nil {
		return err
	}

	m.log.Info("initialized")
	return nil
}
