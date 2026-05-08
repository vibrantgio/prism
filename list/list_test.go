package list_test

import (
	"image"
	"image/color"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	golden "github.com/vibrantgio/prism/internal/golden"
	"github.com/vibrantgio/prism/list"
)

const (
	rowPx = 30  // fixed row height in pixels
	viewW = 200 // viewport width in pixels
	viewH = 150 // viewport height in pixels; fits exactly 5 rows
)

// rowColor maps item index to a unique shade of gray in [15, 243].
// Items 0..19 all produce distinct values, so short/long/scrolled-mid
// goldens are visually distinguishable.
func rowColor(item int) color.NRGBA {
	v := uint8(15 + (item%19)*12)
	return color.NRGBA{R: v, G: v, B: v, A: 0xff}
}

// colorRowFn draws a solid-gray row of fixed pixel height.
func colorRowFn(gtx layout.Context, item int) layout.Dimensions {
	size := image.Pt(gtx.Constraints.Max.X, rowPx)
	paint.FillShape(gtx.Ops, rowColor(item), clip.Rect{Max: size}.Op())
	return layout.Dimensions{Size: size}
}

func makeItems(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

// TestListGolden records or diffs three canonical list configurations:
// short (all items visible), long (many items at top), scrolled-mid (scrolled
// into the middle of a list).
func TestListGolden(t *testing.T) {
	size := image.Pt(viewW, viewH)
	cases := []struct {
		name  string
		items []int
		state *list.State
	}{
		{"short", makeItems(3), list.NewState()},
		{"long", makeItems(20), list.NewState()},
		{"scrolled-mid", makeItems(20), list.NewStateAt(8)},
	}
	for _, tc := range cases {
		items := tc.items
		state := tc.state
		t.Run(tc.name, func(t *testing.T) {
			golden.Render(t, tc.name, size, func(gtx layout.Context) layout.Dimensions {
				return list.Layout(gtx, state, items, colorRowFn)
			})
		})
	}
}

// TestLayoutCallsRowFnOnlyForVisibleItems is the counter-based proof that
// Layout calls rowFn O(visible) times rather than O(len(items)).
//
// With viewH=150px and rowPx=30px, exactly 5 rows fit in the viewport.
// Gio may fetch one extra item for look-ahead, so ≤10 calls is the safe bound.
// This is verified against a 1000-item list where O(total) would be 1000 calls.
func TestLayoutCallsRowFnOnlyForVisibleItems(t *testing.T) {
	items := makeItems(1000)
	state := list.NewState()

	var calls int
	var ops op.Ops
	gtx := layout.Context{
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
		Constraints: layout.Exact(image.Pt(viewW, viewH)),
		Ops:         &ops,
	}

	list.Layout(gtx, state, items, func(gtx layout.Context, item int) layout.Dimensions {
		calls++
		return colorRowFn(gtx, item)
	})

	const maxVisible = 10 // 5 rows + generous look-ahead buffer; far less than 1000
	if calls > maxVisible {
		t.Errorf("rowFn called %d times for 1000-item list (viewport %dpx, row %dpx); want ≤ %d (O(visible))",
			calls, viewH, rowPx, maxVisible)
	}
	if calls == 0 {
		t.Error("rowFn never called; list should render at least 1 row")
	}
}
