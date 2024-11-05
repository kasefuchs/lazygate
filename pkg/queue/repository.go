package queue

// Repository contains queues.
type Repository struct {
	data map[string]Queue
}

// NewRepository creates new instance of repository.
func NewRepository() *Repository {
	return &Repository{
		data: make(map[string]Queue),
	}
}

// Get receives queue by name.
func (r *Repository) Get(name string) Queue {
	return r.data[name]
}

// Push add queue to this repository
func (r *Repository) Push(queue Queue) {
	r.data[queue.Name()] = queue
}
