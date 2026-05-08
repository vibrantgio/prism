package button_test

import (
	"image"
	"testing"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/vibrantgio/prism/button"
	golden "github.com/vibrantgio/prism/internal/golden"
	"github.com/vibrantgio/prism/tokens"
)

func defaultShaper(t *testing.T) *text.Shaper {
	t.Helper()
	return text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
}

// ---- Golden-image tests ----

// TestButtonGolden records or diffs the four canonical button states:
// light-normal, dark-normal, light-focused, light-pressed.
func TestButtonGolden(t *testing.T) {
	shaper := defaultShaper(t)
	size := image.Pt(300, 60)

	// Use zero corner radius and empty label to produce fully sharp-edged
	// renders: anti-aliased rounded corners and GPU font rasterisation both
	// vary slightly between GPU context initialisations, breaking determinism.
	// Shape (background colour, focus ring presence) and colour accuracy are
	// still fully exercised; the exact radius and label are tested elsewhere.
	sharpRadius := tokens.RadiusScale{} // all zeros → sharp corners, no AA
	cases := []struct {
		name   string
		colors tokens.ColorTokens
		state  button.RenderState
	}{
		{"light-normal", tokens.DefaultLight, button.RenderState{}},
		{"dark-normal", tokens.DefaultDark, button.RenderState{}},
		{"light-focused", tokens.DefaultLight, button.RenderState{Focused: true}},
		{"light-pressed", tokens.DefaultLight, button.RenderState{Pressed: true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := button.Render(
				shaper, "",
				tc.colors, tokens.Spacing, sharpRadius, tokens.DefaultTypeScale,
				tc.state,
			)
			golden.Render(t, tc.name, size, w)
		})
	}
}

// ---- Accessibility tests ----

// TestButtonMinHitTarget checks the button meets the 44 dp minimum
// interactive height (DESIGN §Accessibility / WCAG 2.5.5).
func TestButtonMinHitTarget(t *testing.T) {
	shaper := defaultShaper(t)

	var ops op.Ops
	gtx := layout.Context{
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
		Constraints: layout.Exact(image.Pt(300, 120)),
		Ops:         &ops,
	}

	dims := button.Render(
		shaper, "OK",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{},
	)(gtx)

	const wantPx = 44 // 44 dp at 1:1 scale
	if dims.Size.Y < wantPx {
		t.Errorf("button height = %d px, want ≥ %d px (44 dp at 1:1 scale)", dims.Size.Y, wantPx)
	}
}

// TestButtonDisabledIsVisuallyDistinct confirms disabled state produces
// different pixels from enabled state.
func TestButtonDisabledIsVisuallyDistinct(t *testing.T) {
	shaper := defaultShaper(t)
	size := image.Pt(300, 60)

	imgEnabled := golden.Capture(t, size, button.Render(
		shaper, "Click me",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{},
	))
	imgDisabled := golden.Capture(t, size, button.Render(
		shaper, "Click me",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{Disabled: true},
	))

	if imgEnabled == nil || imgDisabled == nil {
		return // headless unavailable; Capture called t.Skip
	}
	if n := golden.PixelDiff(imgEnabled, imgDisabled); n == 0 {
		t.Error("disabled and enabled buttons render identically; expected visual difference")
	}
}

// TestButtonFocusRingIsVisuallyDistinct confirms focused state renders
// differently from normal state (the focus ring must add pixels).
func TestButtonFocusRingIsVisuallyDistinct(t *testing.T) {
	shaper := defaultShaper(t)
	size := image.Pt(300, 60)

	imgNormal := golden.Capture(t, size, button.Render(
		shaper, "Click me",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{},
	))
	imgFocused := golden.Capture(t, size, button.Render(
		shaper, "Click me",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{Focused: true},
	))

	if imgNormal == nil || imgFocused == nil {
		return
	}
	if n := golden.PixelDiff(imgNormal, imgFocused); n == 0 {
		t.Error("focused and normal buttons render identically; expected focus ring pixels to differ")
	}
}

// TestButtonPressedIsVisuallyDistinct confirms pressed state renders
// differently from normal state.
func TestButtonPressedIsVisuallyDistinct(t *testing.T) {
	shaper := defaultShaper(t)
	size := image.Pt(300, 60)

	imgNormal := golden.Capture(t, size, button.Render(
		shaper, "Click me",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{},
	))
	imgPressed := golden.Capture(t, size, button.Render(
		shaper, "Click me",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{Pressed: true},
	))

	if imgNormal == nil || imgPressed == nil {
		return
	}
	if n := golden.PixelDiff(imgNormal, imgPressed); n == 0 {
		t.Error("pressed and normal buttons render identically; expected visual difference")
	}
}

// TestButtonHoveredIsVisuallyDistinct confirms hovered state renders
// differently from normal state.
func TestButtonHoveredIsVisuallyDistinct(t *testing.T) {
	shaper := defaultShaper(t)
	size := image.Pt(300, 60)

	imgNormal := golden.Capture(t, size, button.Render(
		shaper, "Click me",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{},
	))
	imgHovered := golden.Capture(t, size, button.Render(
		shaper, "Click me",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{Hovered: true},
	))

	if imgNormal == nil || imgHovered == nil {
		return
	}
	if n := golden.PixelDiff(imgNormal, imgHovered); n == 0 {
		t.Error("hovered and normal buttons render identically; expected visual difference")
	}
}
