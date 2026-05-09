// Package list provides a virtual-scrolling list component for Gio.
//
// Only items in the current viewport are laid out — O(visible), not O(total).
//
// For lists with reorderable rows that contain interactive Gio widgets (editors,
// checkboxes, etc.), pair with keyed.Defer from prism/keyed to keep per-row
// widget state stable across reorders.
package list

import "gioui.org/layout"

// State holds the scroll position across frames.
// Allocate once per list instance and reuse on every frame.
type State struct {
	l layout.List
}

// NewState returns a State for a vertical list starting at the top.
func NewState() *State {
	return &State{l: layout.List{Axis: layout.Vertical}}
}

// NewStateAt returns a State whose initial first-visible item index is first.
// Intended for golden-image testing; production code uses NewState and lets
// pointer events drive scrolling.
func NewStateAt(first int) *State {
	return &State{l: layout.List{
		Axis:     layout.Vertical,
		Position: layout.Position{First: first},
	}}
}

// Layout lays out items in a virtual scrolling list. rowFn is called only for
// items in the current viewport: O(visible), not O(len(items)).
//
// rowFn must not retain gtx past its call; the closure is invoked once per
// visible item inside the layout.List callback.
func Layout[T any](
	gtx layout.Context,
	state *State,
	items []T,
	rowFn func(gtx layout.Context, item T) layout.Dimensions,
) layout.Dimensions {
	return state.l.Layout(gtx, len(items), func(gtx layout.Context, i int) layout.Dimensions {
		return rowFn(gtx, items[i])
	})
}
