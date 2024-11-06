package registry

import (
	"github.com/kasefuchs/lazygate/pkg/provider"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Registry contains scheduler entries stored by server names.
type Registry struct {
	data map[string]*Entry // Internal registry data.

	proxy    *proxy.Proxy      // Plugin proxy.
	provider provider.Provider // Plugin provider.
}

// NewRegistry creates new instance of Registry.
func NewRegistry(proxy *proxy.Proxy, provider provider.Provider) *Registry {
	return &Registry{
		data: make(map[string]*Entry),

		proxy:    proxy,
		provider: provider,
	}
}

// Clear clears all entries from registry.
func (r *Registry) Clear() {
	r.data = make(map[string]*Entry)
}

// Refresh updates registry data with new info.
func (r *Registry) Refresh(namespace string) {
	r.Clear()

	for _, srv := range r.proxy.Servers() {
		alloc, err := r.provider.AllocationGet(srv)
		if err != nil {
			continue
		}

		cfg, err := alloc.Config()
		if err != nil {
			continue
		}

		if cfg.Namespace == namespace {
			ent := NewEntry(srv, alloc)
			r.EntryRegister(ent)
		}
	}
}

// EntryGet returns entry from registry
func (r *Registry) EntryGet(srv proxy.RegisteredServer) *Entry {
	name := srv.ServerInfo().Name()
	return r.data[name]
}

// EntryList returns list of entries in registry.
func (r *Registry) EntryList() []*Entry {
	res := make([]*Entry, 0, len(r.data))
	for _, ent := range r.data {
		res = append(res, ent)
	}

	return res
}

// EntryRegister registers new entry to registry.
func (r *Registry) EntryRegister(ent *Entry) *Entry {
	name := ent.Server.ServerInfo().Name()
	if _, ok := r.data[name]; !ok {
		r.data[name] = ent
	}

	return r.data[name]
}
