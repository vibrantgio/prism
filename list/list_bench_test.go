package list_test

import (
	"fmt"
	"image"
	"testing"

	"gioui.org/layout"

	"github.com/vibrantgio/prism/bench"
	"github.com/vibrantgio/prism/list"
)

// BenchmarkListLayout demonstrates O(visible) layout cost: ns/op stays roughly
// constant as total item count grows from 10 to 10000, because only the
// viewport-visible rows (~5 at viewH=150px, rowPx=30px) are ever laid out.
// Each sub-case plugs into the shared bench.BenchFrame harness; the BASELINE.md
// "Phase 1 component baseline" row is taken from the N=1000 sub-case.
func BenchmarkListLayout(b *testing.B) {
	for _, n := range []int{10, 100, 1000, 10000} {
		items := makeItems(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			state := list.NewState()
			w := func(gtx layout.Context) layout.Dimensions {
				return list.Layout(gtx, state, items, colorRowFn)
			}
			bench.BenchFrame(b, w, bench.WithSize(image.Pt(viewW, viewH)))
		})
	}
}
