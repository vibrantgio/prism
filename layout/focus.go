package layout

import (
	"gioui.org/io/event"
	"gioui.org/io/key"
	gio "gioui.org/layout"
)

// FocusGroup tracks keyboard focus among a fixed set of interactive items.
// Allocate once per logical group; call Update every frame before registering
// items via their tags.
//
// Usage pattern:
//
//	var g layout.FocusGroup
//	g.Grow(3)
//
//	// each frame:
//	g.Update(gtx)
//	for i := 0; i < g.Len(); i++ {
//	    // inside a clip area: event.Op(gtx.Ops, g.Tag(i))
//	}
//	focused := g.Focused() // -1 if none
type FocusGroup struct {
	items     []focusItem
	focused   int // valid only when haveFocus is true
	haveFocus bool
}

// focusItem is a non-zero-size type so that each element of a FocusGroup's
// items slice has a distinct address, which is required for unique event.Tag
// identity. A zero-size type (struct{}) would cause all items to share the
// same address in a slice, breaking tag-based event routing.
type focusItem [1]byte

// Grow ensures the group has at least n items. Items already present are
// left untouched.
func (g *FocusGroup) Grow(n int) {
	for len(g.items) < n {
		g.items = append(g.items, focusItem{})
	}
}

// Len returns the number of items in the group.
func (g *FocusGroup) Len() int { return len(g.items) }

// Tag returns the event.Tag for item i. Pass it to event.Op inside a clip
// area during layout so the item participates in focus traversal.
// i must be in [0, Len()).
func (g *FocusGroup) Tag(i int) event.Tag { return &g.items[i] }

// Focused returns the index of the currently focused item, or -1 if no
// item in the group currently has focus.
func (g *FocusGroup) Focused() int {
	if !g.haveFocus {
		return -1
	}
	return g.focused
}

// Update refreshes focus state from the current router state. Call once per
// frame before the layout pass that registers item tags via event.Op.
//
// Each item's FocusFilter is registered so the router retains focus when
// assigned to an item. gtx.Focused is used to query focus rather than
// draining FocusEvents, because focus commands applied between frames take
// effect immediately in the router's state.
func (g *FocusGroup) Update(gtx gio.Context) {
	g.haveFocus = false
	for i := range g.items {
		// Register the item as focusable so the router keeps focus when set.
		// Drain any synthetic reset events the router delivers on first use.
		for {
			if _, ok := gtx.Event(key.FocusFilter{Target: &g.items[i]}); !ok {
				break
			}
		}
		if gtx.Focused(&g.items[i]) {
			g.focused = i
			g.haveFocus = true
		}
	}
}
