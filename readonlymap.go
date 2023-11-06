package carbon

// A MapEntry represents a key and its associated value in a ReadonlyMap.
type MapEntry[K comparable, V any] struct {

	// Key is the key for the entry within the ReadonlyMap.
	Key K

	// Value is the value for the entry within the ReadonlyMap.
	Value V
}

// A ReadonlyMap is the read-only interface of the primitive map type.
type ReadonlyMap[K comparable, V any] interface {

	// Len returns the number of elements in the map.
	Len() int

	// Get returns the value in the map for the given key if present, and a boolean indicating
	// whether the value was found.
	Get(key K) (V, bool)

	// Entries returns the entries in the map.
	Entries() []MapEntry[K, V]
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

func (m readonlyMap[K, V]) Entries() []MapEntry[K, V] {
	entries := make([]MapEntry[K, V], 0, len(m))
	for k, v := range m {
		entries = append(entries, MapEntry[K, V]{k, v})
	}
	return entries
}

// CopyMap makes a copy of the given ReadonlyMap that owns its own data so changes to the original
// underlying map wont affect the new ReadonlyMap's content.
func CopyMap[K comparable, V any](m ReadonlyMap[K, V]) ReadonlyMap[K, V] {
	return NewReadonlyMap(CopyMapMut(m))
}

// CopyMapMut copies the contents of the given ReadonlyMap into a new map.
func CopyMapMut[K comparable, V any](m ReadonlyMap[K, V]) map[K]V {
	c := make(map[K]V)
	if m, ok := m.(readonlyMap[K, V]); ok {
		for k, v := range m {
			c[k] = v
		}
		return c
	}
	for _, entry := range m.Entries() {
		c[entry.Key] = entry.Value
	}
	return c
}
