// Package layout provides spacing helpers, a FocusGroup, and flex/grid wrappers
// for Gio applications. It forms part of the Prism component foundation.
package layout

import (
	gio "gioui.org/layout"
	"gioui.org/unit"
)

// Inset returns a uniform inset with dp device-independent pixels on all sides.
func Inset(dp float32) gio.Inset {
	return gio.UniformInset(unit.Dp(dp))
}

// InsetXY returns an inset with h horizontal (left+right) and v vertical
// (top+bottom) device-independent padding.
func InsetXY(h, v float32) gio.Inset {
	return gio.Inset{Left: unit.Dp(h), Right: unit.Dp(h), Top: unit.Dp(v), Bottom: unit.Dp(v)}
}

// HSpacer returns a Widget that occupies dp device-independent pixels horizontally
// and no vertical space. Use as a gap between children inside a Row.
func HSpacer(dp float32) gio.Widget {
	s := gio.Spacer{Width: unit.Dp(dp)}
	return func(gtx gio.Context) gio.Dimensions {
		return s.Layout(gtx)
	}
}

// VSpacer returns a Widget that occupies dp device-independent pixels vertically
// and no horizontal space. Use as a gap between children inside a Col.
func VSpacer(dp float32) gio.Widget {
	s := gio.Spacer{Height: unit.Dp(dp)}
	return func(gtx gio.Context) gio.Dimensions {
		return s.Layout(gtx)
	}
}
