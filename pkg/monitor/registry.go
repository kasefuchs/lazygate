package monitor

import "go.minekube.com/gate/pkg/edition/java/proxy"

// Registry contains monitor entries stored by server names.
type Registry struct {
	v map[string]*Entry
}

// NewRegistry creates new instance of Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Clear clears all entries from registry.
func (r *Registry) Clear() {
	r.v = make(map[string]*Entry)
}

// EntryGet returns entry from registry
func (r *Registry) EntryGet(srv proxy.RegisteredServer) *Entry {
	name := srv.ServerInfo().Name()
	return r.v[name]
}

// EntryList returns list of entries in registry.
func (r *Registry) EntryList() []*Entry {
	res := make([]*Entry, 0, len(r.v))
	for _, ent := range r.v {
		res = append(res, ent)
	}

	return res
}

// EntryRegister registers new entry to registry.
func (r *Registry) EntryRegister(ent *Entry) *Entry {
	name := ent.Server.ServerInfo().Name()
	if _, ok := r.v[name]; !ok {
		r.v[name] = ent
	}

	return r.v[name]
}
