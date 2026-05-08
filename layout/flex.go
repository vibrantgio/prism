package layout

import gio "gioui.org/layout"

// Row lays out widgets as Rigid children in a horizontal Flex row, left-to-right.
// For mixed Rigid/Flexed children use gioui.org/layout.Flex directly.
func Row(gtx gio.Context, children ...gio.Widget) gio.Dimensions {
	return gio.Flex{}.Layout(gtx, rigid(children)...)
}

// Col lays out widgets as Rigid children in a vertical Flex column, top-to-bottom.
// For mixed Rigid/Flexed children use gioui.org/layout.Flex directly.
func Col(gtx gio.Context, children ...gio.Widget) gio.Dimensions {
	return gio.Flex{Axis: gio.Vertical}.Layout(gtx, rigid(children)...)
}

func rigid(ws []gio.Widget) []gio.FlexChild {
	cs := make([]gio.FlexChild, len(ws))
	for i, w := range ws {
		cs[i] = gio.Rigid(w)
	}
	return cs
}
