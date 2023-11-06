package carbon

// A ReadonlySlice is the read-only interface of the primitive slice type.
type ReadonlySlice[V any] interface {

	// Len returns the number of elements in the slice.
	Len() int

	// Get returns the element at the given index of the slice. If the index is out of range, then
	// Get will panic.
	Get(index int) V
}

// NewReadonlySlice returns a ReadonlySlice that wraps the given slice. Note that NewReadonlySlice
// does not make a copy so changes to the slice will affect the the NewReadonlySlice's content.
func NewReadonlySlice[V any](s []V) ReadonlySlice[V] {
	return readonlySlice[V](s)
}

type readonlySlice[V any] []V

func (s readonlySlice[V]) Len() int {
	return len(s)
}

func (s readonlySlice[V]) Get(index int) V {
	return s[index]
}

// CopySlice makes a copy of the given ReadonlySlice that owns its own data so changes to the
// original underlying slice wont affect the new ReadonlySlice's content.
func CopySlice[V any](s ReadonlySlice[V]) ReadonlySlice[V] {
	return NewReadonlySlice(CopySliceMut(s))
}

// CopySliceMut copies the contents of the given ReadonlySlice into a new slice.
func CopySliceMut[V any](s ReadonlySlice[V]) []V {
	l := s.Len()
	c := make([]V, l)
	if s, ok := s.(readonlySlice[V]); ok {
		copy(c, s)
		return c
	}
	for i := 0; i < l; i++ {
		c[i] = s.Get(i)
	}
	return c
}
