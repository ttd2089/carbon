package carbon

// A ReadonlyMap is the read-only interface of the primitive map type.
type ReadonlyMap[K comparable, V any] interface {

	// Len returns the number of elements in the map.
	Len() int

	// Get returns the value in the map for the given key if present, and a boolean indicating
	// whether the value was found.
	Get(K) (V, bool)
}

// NewReadonlyMap returns a ReadonlyMap that wraps the given map. Note that NewReadonlyMap does not
// make a copy so changes to the map will affect the the ReadonlyMap's content.
func NewReadonlyMap[K comparable, V any](m map[K]V) ReadonlyMap[K, V] {
	return readonlyMap[K, V](m)
}

type readonlyMap[K comparable, V any] map[K]V

func (m readonlyMap[K, V]) Len() int {
	return len(m)
}

func (m readonlyMap[K, V]) Get(key K) (V, bool) {
	v, ok := m[key]
	return v, ok
}
