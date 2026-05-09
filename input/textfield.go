// Package input provides Prism input components for Gio applications.
package input

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
// Gio event system via TextField.
type RenderState struct {
	Focused  bool
	Disabled bool
}

// TextFieldProps configures a TextField instance.
type TextFieldProps struct {
	// Placeholder is shown when the field is empty and unfocused.
	Placeholder string

	// Description is the screen-reader label. Falls back to Placeholder when empty.
	Description string

	// Disabled, if non-nil, disables the field when it emits true.
	Disabled rx.Observable[bool]

	// OnChange is called with the new value on every text change.
	// This is the FRP callback path.
	OnChange func(string)

	// Message, if non-nil, causes the field to emit mvu.MessageOp{Message}
	// on every text change. This is the MVU integration path.
	Message any

	// Shaper, if nil, defaults to a shaper backed by Go fonts.
	Shaper *text.Shaper
}

// resolvedTokens is the concrete per-emission snapshot consumed by the widget closure.
type resolvedTokens struct {
	color   tokens.ColorTokens
	typ     tokens.TypeScale
	spacing tokens.SpacingScale
	radius  tokens.RadiusScale
}

// TextField returns an rx.Observable[layout.Widget] that emits a new widget
// whenever the theme or disabled state changes. Widget state (editor content,
// focus) lives in the rx.Defer scope and persists across emissions.
//
// Both integration paths are supported:
//   - FRP: set TextFieldProps.OnChange; FRP consumers wrap with rx.NewSubject if needed.
//   - MVU: set TextFieldProps.Message; the component emits mvu.MessageOp on text change.
func TextField(th rx.Observable[theme.Theme], props TextFieldProps) rx.Observable[layout.Widget] {
	disabled := props.Disabled
	if disabled == nil {
		disabled = rx.Of(false)
	}

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
		// emissions for the lifetime of this TextField instance.
		editor := &widget.Editor{SingleLine: true}
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

				// Drain editor events; fire callbacks on text change.
				for {
					ev, ok := editor.Update(gtx)
					if !ok {
						break
					}
					if _, ok := ev.(widget.ChangeEvent); ok {
						val := editor.Text()
						if props.OnChange != nil {
							props.OnChange(val)
						}
						if props.Message != nil {
							mvu.MessageOp{Message: props.Message}.Add(gtx.Ops)
						}
					}
				}

				desc := props.Description
				if desc == "" {
					desc = props.Placeholder
				}

				foc := !dis && gtx.Focused(editor)
				showPh := !foc && editor.Len() == 0

				return drawTextFieldLive(gtx, shaper, editor, props.Placeholder, desc, tok, RenderState{
					Focused:  foc,
					Disabled: dis,
				}, showPh)
			}
		})
	})
}

// Render produces a layout.Widget for a text field in an explicit visual state,
// without any event processing or rx machinery. Intended for golden-image
// testing and static demonstrations; production code should use TextField.
func Render(
	shaper *text.Shaper,
	placeholder string,
	colors tokens.ColorTokens,
	sp tokens.SpacingScale,
	rad tokens.RadiusScale,
	ts tokens.TypeScale,
	s RenderState,
) layout.Widget {
	tok := resolvedTokens{color: colors, spacing: sp, radius: rad, typ: ts}
	return func(gtx layout.Context) layout.Dimensions {
		return drawTextFieldStatic(gtx, shaper, placeholder, tok, s)
	}
}

// drawTextFieldLive renders a live text field containing a widget.Editor.
func drawTextFieldLive(gtx layout.Context, shaper *text.Shaper, editor *widget.Editor, placeholder, desc string, tok resolvedTokens, s RenderState, showPlaceholder bool) layout.Dimensions {
	padH := gtx.Dp(unit.Dp(tok.spacing.S3))
	padV := gtx.Dp(unit.Dp(tok.spacing.S2))
	minH := gtx.Dp(minHeight)
	rad := gtx.Dp(unit.Dp(tok.radius.Md))
	textSize := unit.Sp(tok.typ.BodyLarge)

	bg, textColor, borderColor, phColor := textFieldColors(tok.color, s)

	fieldW := gtx.Constraints.Max.X
	innerW := fieldW - 2*padH
	if innerW < 1 {
		innerW = 1
	}

	innerGtx := gtx
	innerGtx.Constraints = layout.Constraints{
		Min: image.Pt(0, 0),
		Max: image.Pt(innerW, gtx.Constraints.Max.Y),
	}

	// Semantic accessibility ops.
	semantic.ClassOp(semantic.Editor).Add(gtx.Ops)
	semantic.DescriptionOp(desc).Add(gtx.Ops)
	semantic.EnabledOp(!s.Disabled).Add(gtx.Ops)

	// Measure content height via recorded label layout so we can size the
	// field before drawing the background. Replay if placeholder is needed.
	mPhColor := op.Record(gtx.Ops)
	paint.ColorOp{Color: phColor}.Add(gtx.Ops)
	phMat := mPhColor.Stop()

	mPh := op.Record(gtx.Ops)
	wl := widget.Label{MaxLines: 1}
	contentDims := wl.Layout(innerGtx, shaper, font.Font{}, textSize, placeholder, phMat)
	phCall := mPh.Stop()

	fieldH := contentDims.Size.Y + 2*padV
	if fieldH < minH {
		fieldH = minH
	}
	fieldSize := image.Pt(fieldW, fieldH)

	// Border as nested fills: outer rect in border color, inner rect in
	// background color. Avoids clip.Stroke anti-aliasing variance in tests.
	borderPx := gtx.Dp(1)
	if s.Focused {
		borderPx = gtx.Dp(2)
	}
	innerRad := rad - borderPx
	if innerRad < 0 {
		innerRad = 0
	}
	rrectOuter := clip.RRect{Rect: image.Rectangle{Max: fieldSize}, SE: rad, SW: rad, NE: rad, NW: rad}
	paint.FillShape(gtx.Ops, borderColor, rrectOuter.Op(gtx.Ops))
	rrectInner := clip.RRect{
		Rect: image.Rectangle{
			Min: image.Pt(borderPx, borderPx),
			Max: image.Pt(fieldSize.X-borderPx, fieldSize.Y-borderPx),
		},
		SE: innerRad, SW: innerRad, NE: innerRad, NW: innerRad,
	}
	paint.FillShape(gtx.Ops, bg, rrectInner.Op(gtx.Ops))

	offY := (fieldH - contentDims.Size.Y) / 2

	// Placeholder overlay (only when empty and unfocused).
	if showPlaceholder {
		st := op.Offset(image.Pt(padH, offY)).Push(gtx.Ops)
		phCall.Add(gtx.Ops)
		st.Pop()
	}

	// Editor materials.
	mText := op.Record(gtx.Ops)
	paint.ColorOp{Color: textColor}.Add(gtx.Ops)
	textMat := mText.Stop()

	mSel := op.Record(gtx.Ops)
	paint.ColorOp{Color: withAlpha(tok.color.Primary, 0x40)}.Add(gtx.Ops)
	selMat := mSel.Stop()

	// Editor — always laid out so it receives pointer/keyboard events.
	// Min.X = innerW ensures the event clip covers the full field width even
	// when the editor is empty (otherwise the clip shrinks to the caret width).
	editorGtx := innerGtx
	editorGtx.Constraints = layout.Constraints{
		Min: image.Pt(innerW, 0),
		Max: image.Pt(innerW, contentDims.Size.Y),
	}
	st := op.Offset(image.Pt(padH, offY)).Push(gtx.Ops)
	editor.Layout(editorGtx, shaper, font.Font{}, textSize, textMat, selMat)
	st.Pop()

	if !s.Disabled {
		pointer.CursorText.Add(gtx.Ops)
	}

	return layout.Dimensions{Size: fieldSize}
}

// drawTextFieldStatic renders a static text field for golden-image testing.
// It always shows the placeholder text; there is no live editor.
func drawTextFieldStatic(gtx layout.Context, shaper *text.Shaper, placeholder string, tok resolvedTokens, s RenderState) layout.Dimensions {
	padH := gtx.Dp(unit.Dp(tok.spacing.S3))
	padV := gtx.Dp(unit.Dp(tok.spacing.S2))
	minH := gtx.Dp(minHeight)
	rad := gtx.Dp(unit.Dp(tok.radius.Md))
	textSize := unit.Sp(tok.typ.BodyLarge)

	bg, _, borderColor, phColor := textFieldColors(tok.color, s)

	fieldW := gtx.Constraints.Max.X
	innerW := fieldW - 2*padH
	if innerW < 1 {
		innerW = 1
	}

	innerGtx := gtx
	innerGtx.Constraints = layout.Constraints{
		Min: image.Pt(0, 0),
		Max: image.Pt(innerW, gtx.Constraints.Max.Y),
	}

	// Record placeholder label for measurement and deferred rendering.
	mPhColor := op.Record(gtx.Ops)
	paint.ColorOp{Color: phColor}.Add(gtx.Ops)
	phMat := mPhColor.Stop()

	mLabel := op.Record(gtx.Ops)
	wl := widget.Label{MaxLines: 1}
	labelDims := wl.Layout(innerGtx, shaper, font.Font{}, textSize, placeholder, phMat)
	labelCall := mLabel.Stop()

	fieldH := labelDims.Size.Y + 2*padV
	if fieldH < minH {
		fieldH = minH
	}
	fieldSize := image.Pt(fieldW, fieldH)

	// Border as nested fills: outer rect in border color, inner rect in
	// background color. Avoids clip.Stroke anti-aliasing variance in tests.
	borderPx := gtx.Dp(1)
	if s.Focused {
		borderPx = gtx.Dp(2)
	}
	innerRad := rad - borderPx
	if innerRad < 0 {
		innerRad = 0
	}
	rrectOuter := clip.RRect{Rect: image.Rectangle{Max: fieldSize}, SE: rad, SW: rad, NE: rad, NW: rad}
	paint.FillShape(gtx.Ops, borderColor, rrectOuter.Op(gtx.Ops))
	rrectInner := clip.RRect{
		Rect: image.Rectangle{
			Min: image.Pt(borderPx, borderPx),
			Max: image.Pt(fieldSize.X-borderPx, fieldSize.Y-borderPx),
		},
		SE: innerRad, SW: innerRad, NE: innerRad, NW: innerRad,
	}
	paint.FillShape(gtx.Ops, bg, rrectInner.Op(gtx.Ops))

	// Placeholder label centered vertically.
	offY := (fieldH - labelDims.Size.Y) / 2
	st := op.Offset(image.Pt(padH, offY)).Push(gtx.Ops)
	labelCall.Add(gtx.Ops)
	st.Pop()

	return layout.Dimensions{Size: fieldSize}
}

// textFieldColors returns (bg, text, border, placeholder) colors for the given state.
func textFieldColors(c tokens.ColorTokens, s RenderState) (bg, text, border, placeholder color.NRGBA) {
	bg = c.Surface
	text = c.OnSurface
	border = c.Outline
	placeholder = withAlpha(c.OnSurfaceVariant, 0x99)
	switch {
	case s.Disabled:
		bg = withAlpha(c.SurfaceVariant, 0x80)
		text = withAlpha(c.OnSurface, 0x61)
		border = withAlpha(c.Outline, 0x61)
		placeholder = withAlpha(c.OnSurfaceVariant, 0x40)
	case s.Focused:
		border = c.Primary
	}
	return
}

// withAlpha returns c with its alpha scaled by factor a (0–255).
func withAlpha(c color.NRGBA, a uint8) color.NRGBA {
	c.A = uint8(uint16(c.A) * uint16(a) / 255)
	return c
}
