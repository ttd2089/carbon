package carbon

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ttd2089/proof"
)

func TestReadonlyMap(t *testing.T) {

	t.Run("Len()", func(t *testing.T) {

		t.Run("returns 0 for nil slices", func(t *testing.T) {
			underTest := NewReadonlySlice[string](nil)
			proof.AssertEq(t, "underTest.Len()", underTest.Len(), 0)
		})

		t.Run("returns 0 for empty slices", func(t *testing.T) {
			underTest := NewReadonlySlice(make([]string, 0))
			proof.AssertEq(t, "underTest.Len()", underTest.Len(), 0)
		})

		t.Run("returns length of non-empty slices", func(t *testing.T) {
			underTest := NewReadonlySlice([]string{"zero", "one", "two", "three"})
			proof.AssertEq(t, "underTest.Len()", underTest.Len(), 4)
		})
	})

	t.Run("Get()", func(t *testing.T) {

		t.Run("returns the value for the given index in the slice", func(t *testing.T) {
			underTest := NewReadonlySlice([]string{"zero", "one", "two", "three"})
			proof.AssertEq(t, "underTest.Get(2)", underTest.Get(2), "two")
		})
	})
}

func TestCopySlice(t *testing.T) {

	t.Run("returns a copy of the ReadonlySlice", func(t *testing.T) {
		original := NewReadonlySlice([]string{"zero", "one", "two"})
		returned := CopySlice(original)
		proof.AssertEq(t, "returned.Len()", returned.Len(), original.Len())
		for i := 0; i < original.Len(); i++ {
			proof.AssertEq(t, fmt.Sprintf("returned.Get(%d)", i), returned.Get(i), original.Get(i))
		}
	})

	t.Run("does not reflect changes in the original underlying slice", func(t *testing.T) {
		underlying := []string{"zero", "one", "two"}
		original := NewReadonlySlice(underlying)
		returned := CopySlice(original)
		notExpected := "four"
		underlying[0] = notExpected
		proof.ExpectNeq(t, "returned.Get(0)", returned.Get(0), notExpected)
	})
}

func TestCopySliceMut(t *testing.T) {

	for _, tc := range []struct {
		name string
		impl func([]string) ReadonlySlice[string]
	}{
		{
			name: "default implementation",
			impl: func(s []string) ReadonlySlice[string] {
				return NewReadonlySlice(s)
			},
		},
		{
			name: "non-default implementation",
			impl: func(s []string) ReadonlySlice[string] {
				return &readonlySliceWrapper[string]{NewReadonlySlice(s)}
			},
		},
	} {
		t.Run("returns a copy of the underlying slice", func(t *testing.T) {
			underlying := []string{"zero", "one", "two"}
			original := tc.impl(underlying)
			if actual := CopySliceMut(original); !reflect.DeepEqual(actual, underlying) {
				t.Fatalf("expected %+v; got %+v", underlying, actual)
			}
		})

		t.Run("does not reflect changes in riginal underlying map", func(t *testing.T) {
			underlying := []string{"zero", "one", "two"}
			original := tc.impl(underlying)
			returned := CopySliceMut(original)
			notExpected := "four"
			underlying[0] = notExpected
			proof.ExpectNeq(t, "returned[0]", returned[0], notExpected)
		})
	}
}

// The CopySliceMut function contains an optimization for the library provided implementation of
// ReadonlySlice. This wrapper exists so that we can test the function against a non-default
// implementation as well.
type readonlySliceWrapper[V any] struct {
	impl ReadonlySlice[V]
}

func (m *readonlySliceWrapper[V]) Len() int {
	return m.impl.Len()
}

func (m *readonlySliceWrapper[V]) Get(index int) V {
	return m.impl.Get(index)
}

func BenchmarkCopySliceMut(b *testing.B) {

	defaultImpl := NewReadonlySlice([]string{"zero", "one", "two"})
	wrapper := &readonlySliceWrapper[string]{defaultImpl}

	b.Run("default implementation", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = CopySliceMut(defaultImpl)
		}
	})

	b.Run("interface-based implementation", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = CopySliceMut(wrapper)
		}
	})
}
