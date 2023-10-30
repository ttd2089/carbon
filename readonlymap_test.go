package carbon

import (
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
	})
}
