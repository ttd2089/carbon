package carbon

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ttd2089/proof"
)

func TestReadonlySlice(t *testing.T) {

	t.Run("Len()", func(t *testing.T) {

		t.Run("returns 0 for nil maps", func(t *testing.T) {
			underTest := NewReadonlyMap[int, string](nil)
			proof.AssertEq(t, "underTest.Len()", underTest.Len(), 0)
		})

		t.Run("returns 0 for empty maps", func(t *testing.T) {
			underTest := NewReadonlyMap(make(map[int]string))
			proof.AssertEq(t, "underTest.Len()", underTest.Len(), 0)
		})

		t.Run("returns length of non-empty maps", func(t *testing.T) {
			underTest := NewReadonlyMap(map[int]string{
				0: "zero",
				1: "one",
				2: "two",
				3: "three",
			})
			proof.AssertEq(t, "underTest.Len()", underTest.Len(), 4)
		})

		t.Run("reflects changes in the underlying map", func(t *testing.T) {
			underlying := make(map[int]string)
			underTest := NewReadonlyMap(underlying)
			underlying[0] = "zero"
			underlying[1] = "one"
			proof.AssertEq(t, "underTest.Len()", underTest.Len(), 2)
		})
	})

	t.Run("Get()", func(t *testing.T) {

		t.Run("returns false for keys not in the map", func(t *testing.T) {
			underTest := NewReadonlyMap(map[int]string{
				0: "zero",
				1: "one",
			})
			_, ok := underTest.Get(2)
			proof.AssertEq(t, "_, ok := underTest.Get(2); ok", ok, false)
		})

		t.Run("returns true and the value for keys in the map", func(t *testing.T) {
			underTest := NewReadonlyMap(map[int]string{
				0: "zero",
				1: "one",
				2: "two",
				3: "three",
			})
			v, ok := underTest.Get(2)
			proof.AssertEq(t, "_, ok := underTest.Get(2); ok", ok, true)
			proof.AssertEq(t, "v, _ := underTest.Get(2); ok", v, "two")
		})

		t.Run("reflects changes in the underlying map", func(t *testing.T) {
			underlying := make(map[int]string)
			underTest := NewReadonlyMap(underlying)
			underlying[2] = "two"
			v, ok := underTest.Get(2)
			proof.AssertEq(t, "_, ok := underTest.Get(2); ok", ok, true)
			proof.AssertEq(t, "v, _ := underTest.Get(2); ok", v, "two")
		})
	})

	t.Run("Entries()", func(t *testing.T) {
		t.Run("returns each entry in the map", func(t *testing.T) {
			expected := []MapEntry[int, string]{
				{0, "zero"},
				{1, "one"},
				{2, "two"},
				{3, "three"},
			}
			underTest := NewReadonlyMap(map[int]string{
				0: "zero",
				1: "one",
				2: "two",
				3: "three",
			})
			actual := underTest.Entries()
			proof.AssertSetEq(t, "underTest.Entries()", actual, expected)
		})

		t.Run("reflects changes in the underlying map", func(t *testing.T) {
			underlying := map[int]string{
				0: "zero",
				1: "one",
			}
			underTest := NewReadonlyMap(underlying)
			underlying[2] = "two"
			underlying[3] = "three"
			expected := []MapEntry[int, string]{
				{0, "zero"},
				{1, "one"},
				{2, "two"},
				{3, "three"},
			}
			actual := underTest.Entries()
			proof.AssertSetEq(t, "underTest.Entries()", actual, expected)
		})
	})
}

func TestCopyMap(t *testing.T) {

	t.Run("returns a copy of the ReadonlyMap", func(t *testing.T) {
		original := NewReadonlyMap(map[int]string{
			0: "zero",
			1: "one",
			2: "two",
			3: "three",
		})
		returned := CopyMap(original)
		proof.AssertEq(t, "returned.Len()", returned.Len(), original.Len())
		for _, e := range original.Entries() {
			v, ok := returned.Get(e.Key)
			proof.AssertEq(t, fmt.Sprintf("_, ok := returned.Get(%d); ok", e.Key), ok, true)
			proof.AssertEq(t, fmt.Sprintf("v, _ := returned.Get(%d); v", e.Key), v, e.Value)
		}
	})

	t.Run("does not reflect changes in original underlying map", func(t *testing.T) {
		underlying := map[int]string{
			0: "zero",
			1: "one",
			2: "two",
			3: "three",
		}
		original := NewReadonlyMap(underlying)
		returned := CopyMap(original)
		notExpected := "four"
		underlying[0] = notExpected
		actual, _ := returned.Get(0)
		proof.ExpectNeq(t, "v, _ := returned.Get(0); v", actual, notExpected)
	})
}

func TestCopyMapMut(t *testing.T) {

	for _, tc := range []struct {
		name string
		impl func(map[int]string) ReadonlyMap[int, string]
	}{
		{
			name: "default implementation",
			impl: func(m map[int]string) ReadonlyMap[int, string] {
				return NewReadonlyMap(m)
			},
		},
		{
			name: "non-default implementation",
			impl: func(m map[int]string) ReadonlyMap[int, string] {
				return &readonlyMapWrapper[int, string]{NewReadonlyMap(m)}
			},
		},
	} {
		t.Run("returns a copy of the underlying map", func(t *testing.T) {
			underlying := map[int]string{
				0: "zero",
				1: "one",
				2: "two",
				3: "three",
			}
			original := tc.impl(underlying)
			if actual := CopyMapMut(original); !reflect.DeepEqual(actual, underlying) {
				t.Fatalf("expected %+v; got %+v", underlying, actual)
			}
		})

		t.Run("does not reflect changes in riginal underlying map", func(t *testing.T) {
			underlying := map[int]string{
				0: "zero",
				1: "one",
				2: "two",
				3: "three",
			}
			original := tc.impl(underlying)
			returned := CopyMapMut(original)
			notExpected := "four"
			underlying[0] = notExpected
			proof.ExpectNeq(t, "returned[0]", returned[0], notExpected)
		})
	}
}

// The CopyMapMut function contains an optimization for the library provided implementation of
// ReadonlyMap. This wrapper exists so that we can test the function against a non-default
// implementation as well.
type readonlyMapWrapper[K comparable, V any] struct {
	impl ReadonlyMap[K, V]
}

func (m *readonlyMapWrapper[K, V]) Len() int {
	return m.impl.Len()
}

func (m *readonlyMapWrapper[K, V]) Get(key K) (V, bool) {
	return m.impl.Get(key)
}

func (m *readonlyMapWrapper[K, V]) Entries() []MapEntry[K, V] {
	return m.impl.Entries()
}

func BenchmarkCopyMapMut(b *testing.B) {

	defaultImpl := NewReadonlyMap(map[int]string{
		0: "zero",
		1: "one",
		2: "two",
		3: "three",
	})
	wrapper := &readonlyMapWrapper[int, string]{defaultImpl}

	b.Run("default implementation", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = CopyMapMut(defaultImpl)
		}
	})

	b.Run("interface-based implementation", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = CopyMapMut(wrapper)
		}
	})
}
