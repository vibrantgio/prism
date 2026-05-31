// Package button provides a Prism Button component for Gio applications.
//
// The component integrates with both FRP (via rx.Observable) and MVU (via
// mvu.MessageOp) application patterns, per the component contract in
// DESIGN §"Bridging FRP and MVU".
package button

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/mvu"
	"github.com/vibrantgio/prism/theme"
	"github.com/vibrantgio/prism/tokens"
)

// minHeight is the minimum interactive target height (WCAG 2.5.5, DESIGN §Accessibility).
const minHeight = unit.Dp(44)

// RenderState holds explicit visual interaction state for static rendering.
// All fields default to false (normal/idle state).
// Intended for golden-image testing; production code obtains state from the
// Gio event system via Button.
type RenderState struct {
	Hovered  bool
	Focused  bool
	Pressed  bool
	Disabled bool
}

// Props configures a Button instance.
type Props struct {
	// Label is the text rendered inside the button.
	Label string

	// Description is the screen-reader label. Falls back to Label when empty.
	Description string

	// Disabled, if non-nil, disables the button when it emits true.
	// A nil Disabled means always enabled.
	Disabled rx.Observable[bool]

	// OnClick is called when the button is activated by click or Space/Enter.
	// This is the FRP callback path. The gtx argument is the layout.Context
	// active on the frame when the click is processed, allowing consumers to
	// emit mvu.MessageOp{Message: ...}.Add(gtx.Ops) inside the callback.
	OnClick func(gtx layout.Context)

	// Message, if non-nil, causes the button to emit mvu.MessageOp{Message}
	// into gtx.Ops on activation. This is the MVU integration path.
	Message any

	// Shaper, if nil, defaults to a shaper backed by Go fonts.
	// The default shaper is created once per subscription inside the rx.Defer
	// scope, so it is not re-allocated on every theme change.
	Shaper *text.Shaper
}

// resolvedTokens is the concrete per-emission snapshot consumed by the widget closure.
type resolvedTokens struct {
	color   tokens.ColorTokens
	typ     tokens.TypeScale
	spacing tokens.SpacingScale
	radius  tokens.RadiusScale
}

// Button returns an rx.Observable[layout.Widget] that emits a new widget
// whenever the theme or disabled state changes. Widget state (clickable, hover,
// focus, press) lives in the rx.Defer scope and persists across emissions.
//
// Both integration paths are supported:
//   - FRP: set Props.OnClick; FRP consumers wrap with rx.NewSubject if needed.
//   - MVU: set Props.Message; the component emits mvu.MessageOp on activation.
func Button(th rx.Observable[theme.Theme], props Props) rx.Observable[layout.Widget] {
	disabled := props.Disabled
	if disabled == nil {
		disabled = rx.Of(false)
	}

	// Flatten the nested theme observables into a concrete snapshot.
	resolved := rx.SwitchMap(th, func(t theme.Theme) rx.Observable[resolvedTokens] {
		return rx.Map(
			rx.CombineLatest4(t.Color, t.Type, t.Spacing, t.Radius),
			func(n rx.Tuple4[tokens.ColorTokens, tokens.TypeScale, tokens.SpacingScale, tokens.RadiusScale]) resolvedTokens {
				return resolvedTokens{n.First, n.Second, n.Third, n.Fourth}
			},
		)
	})

	inputs := rx.CombineLatest2(resolved, disabled)

	return rx.Defer(func() rx.Observable[layout.Widget] {
		// Allocated once per subscription — survives all theme and disabled
		// emissions for the lifetime of this button instance.
		var click widget.Clickable
		shaper := props.Shaper
		if shaper == nil {
			shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
		}

		return rx.Map(inputs, func(next rx.Tuple2[resolvedTokens, bool]) layout.Widget {
			tok, dis := next.First, next.Second

			return func(gtx layout.Context) layout.Dimensions {
				if dis {
					gtx = gtx.Disabled()
				}

				// Process events; Clicked also handles Space/Enter via widget.Clickable.
				if click.Clicked(gtx) {
					if props.OnClick != nil {
						props.OnClick(gtx)
					}
					if props.Message != nil {
						mvu.MessageOp{Message: props.Message}.Add(gtx.Ops)
					}
				}

				hov := click.Hovered()
				prs := click.Pressed()
				foc := !dis && gtx.Focused(&click)

				desc := props.Description
				if desc == "" {
					desc = props.Label
				}

				return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					semantic.ClassOp(semantic.Button).Add(gtx.Ops)
					semantic.LabelOp(props.Label).Add(gtx.Ops)
					semantic.DescriptionOp(desc).Add(gtx.Ops)
					semantic.EnabledOp(!dis).Add(gtx.Ops)
					return drawButton(gtx, shaper, props.Label, tok, RenderState{
						Hovered:  hov,
						Focused:  foc,
						Pressed:  prs,
						Disabled: dis,
					})
				})
			}
		})
	})
}

// Render produces a layout.Widget for a button in an explicit visual state,
// without any event processing or rx machinery. Intended for golden-image
// testing and static demonstrations; production code should use Button.
func Render(
	shaper *text.Shaper,
	label string,
	colors tokens.ColorTokens,
	sp tokens.SpacingScale,
	rad tokens.RadiusScale,
	ts tokens.TypeScale,
	s RenderState,
) layout.Widget {
	tok := resolvedTokens{color: colors, spacing: sp, radius: rad, typ: ts}
	return func(gtx layout.Context) layout.Dimensions {
		return drawButton(gtx, shaper, label, tok, s)
	}
}

// drawButton renders the button visual into gtx. All visual state comes from s;
// no event queries are performed here.
func drawButton(gtx layout.Context, shaper *text.Shaper, label string, tok resolvedTokens, s RenderState) layout.Dimensions {
	padH := gtx.Dp(unit.Dp(tok.spacing.S4)) // 16 dp horizontal padding
	padV := gtx.Dp(unit.Dp(tok.spacing.S2)) // 8 dp vertical padding
	minH := gtx.Dp(minHeight)               // 44 dp minimum height
	rad := gtx.Dp(unit.Dp(tok.radius.Md))   // 6 dp corner radius

	bg, fg := buttonColors(tok.color, s)

	// Record the text material (fg color op) — replayed inside the label layout.
	mColor := op.Record(gtx.Ops)
	paint.ColorOp{Color: fg}.Add(gtx.Ops)
	textMaterial := mColor.Stop()

	// Record the label render to obtain its size before drawing the background.
	labelGtx := gtx
	labelGtx.Constraints.Min = image.Pt(0, 0)
	maxLabelW := gtx.Constraints.Max.X - 2*padH
	if maxLabelW > 0 {
		labelGtx.Constraints.Max.X = maxLabelW
	}
	mLabel := op.Record(gtx.Ops)
	wl := widget.Label{MaxLines: 1}
	labelDims := wl.Layout(labelGtx, shaper, font.Font{}, unit.Sp(tok.typ.LabelLarge), label, textMaterial)
	labelCall := mLabel.Stop()

	// Button dimensions: fill available width, enforce 44 dp minimum height.
	btnW := gtx.Constraints.Max.X
	if btnW < labelDims.Size.X+2*padH {
		btnW = labelDims.Size.X + 2*padH
	}
	btnH := labelDims.Size.Y + 2*padV
	if btnH < minH {
		btnH = minH
	}
	btnSize := image.Pt(btnW, btnH)

	// Background fill.
	rrect := clip.RRect{Rect: image.Rectangle{Max: btnSize}, SE: rad, SW: rad, NE: rad, NW: rad}
	paint.FillShape(gtx.Ops, bg, rrect.Op(gtx.Ops))

	// Focus ring (2 dp stroke on the background boundary).
	if s.Focused {
		paint.FillShape(gtx.Ops, tok.color.Outline, clip.Stroke{
			Path:  rrect.Path(gtx.Ops),
			Width: float32(gtx.Dp(2)),
		}.Op())
	}

	// Replay the label centered within the button.
	offX := (btnW - labelDims.Size.X) / 2
	offY := (btnH - labelDims.Size.Y) / 2
	st := op.Offset(image.Pt(offX, offY)).Push(gtx.Ops)
	labelCall.Add(gtx.Ops)
	st.Pop()

	if !s.Disabled {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	return layout.Dimensions{Size: btnSize}
}

// buttonColors returns the background and foreground colors for the given state.
func buttonColors(c tokens.ColorTokens, s RenderState) (bg, fg color.NRGBA) {
	bg, fg = c.Primary, c.OnPrimary
	switch {
	case s.Disabled:
		// 38% opacity — WCAG disabled state convention.
		bg = withAlpha(bg, 0x61)
		fg = withAlpha(fg, 0x61)
	case s.Pressed:
		// ~15% black overlay for press feedback.
		bg = blend(bg, color.NRGBA{A: 0x26})
	case s.Focused, s.Hovered:
		// ~10% white overlay for hover/focus feedback.
		bg = blend(bg, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x1a})
	}
	return
}

// withAlpha returns c with its alpha scaled by factor a (0–255).
func withAlpha(c color.NRGBA, a uint8) color.NRGBA {
	c.A = uint8(uint16(c.A) * uint16(a) / 255)
	return c
}

// blend alpha-composites overlay on top of base using straight alpha.
func blend(base, overlay color.NRGBA) color.NRGBA {
	a := float32(overlay.A) / 255
	return color.NRGBA{
		R: uint8(float32(base.R)*(1-a) + float32(overlay.R)*a + 0.5),
		G: uint8(float32(base.G)*(1-a) + float32(overlay.G)*a + 0.5),
		B: uint8(float32(base.B)*(1-a) + float32(overlay.B)*a + 0.5),
		A: base.A,
	}
}
