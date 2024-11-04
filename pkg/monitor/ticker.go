package monitor

// Checks if allocation should be stopped and stops it.
func (m *Monitor) allocationStopTicker() {
	for _, ent := range m.registry.EntryList() {
		if !ent.ShouldStop() {
			continue
		}

		name := ent.Server.ServerInfo().Name()
		m.log.Info("stopping server allocation", "server", name)
		if err := ent.Allocation.Stop(); err != nil {
			m.log.Error(err, "failed to stop server allocation", "server", name)
		}
	}
}
