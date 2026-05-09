package list_test

import (
	"fmt"
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/vibrantgio/prism/list"
)

// BenchmarkListLayout demonstrates O(visible) layout cost: ns/op stays roughly
// constant as total item count grows from 10 to 10000, because only the
// viewport-visible rows (~5 at viewH=150px, rowPx=30px) are ever laid out.
func BenchmarkListLayout(b *testing.B) {
	for _, n := range []int{10, 100, 1000, 10000} {
		items := makeItems(n)
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			state := list.NewState()
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var ops op.Ops
				gtx := layout.Context{
					Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
					Constraints: layout.Exact(image.Pt(viewW, viewH)),
					Ops:         &ops,
				}
				list.Layout(gtx, state, items, colorRowFn)
			}
		})
	}
}
