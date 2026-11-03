package scheduler

import (
	"context"
	"time"

	"codeberg.org/kasefuchs/lazygate/pkg/provider"
	"codeberg.org/kasefuchs/lazygate/pkg/registry"
	"codeberg.org/kasefuchs/lazygate/pkg/utils"
	"github.com/go-logr/logr"
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
