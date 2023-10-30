package carbon

import (
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
