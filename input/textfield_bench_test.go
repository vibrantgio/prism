package input_test

import (
	"image"
	"testing"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	gioinput "gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/prism/bench"
	"github.com/vibrantgio/prism/input"
	"github.com/vibrantgio/prism/theme"
	"github.com/vibrantgio/prism/tokens"
)

// BenchmarkTextFieldRender measures the static, unfocused render path
// (input.Render → drawTextFieldStatic) via the shared bench.BenchFrame harness.
// It is the idle baseline that the live caret-blink frame below is contrasted
// against; b.ReportAllocs is enabled by the harness.
func BenchmarkTextFieldRender(b *testing.B) {
	shaper := text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	w := input.Render(
		shaper, "Placeholder",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		input.RenderState{},
	)
	bench.BenchFrame(b, w)
}

// BenchmarkTextFieldCaretBlink measures the cursor-blinking frame: a live
// widget.Editor that is focused and holds text, so drawTextFieldLive renders
// the caret on every frame (the caret is drawn only while gtx.Focused(editor)
// is true — textfield.go). This is the typing hot path, so it is NOT
// apples-to-apples with the static render above: it includes editor layout plus
// the input.Router frame cost.
//
// Focus is established before the timed loop with the proven click-to-focus
// dance from TestTextFieldSubmitFiresCallbacksAndClears, then a key.EditEvent
// types text. The OnChange guard fails the benchmark loudly if focus or input
// silently did not take — otherwise we would benchmark the unfocused
// placeholder frame under a "caret-blink" label.
func BenchmarkTextFieldCaretBlink(b *testing.B) {
	shaper := text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	var changed string
	props := input.TextFieldProps{
		Placeholder: "Placeholder",
		Shaper:      shaper,
		OnChange:    func(_ layout.Context, s string) { changed = s },
	}

	var w layout.Widget
	if err := input.TextField(rx.Of(theme.Default()), props).Subscribe(func(next layout.Widget, _ error, done bool) {
		if !done && next != nil {
			w = next
		}
	}, rx.NewScheduler()).Wait(); err != nil {
		b.Fatalf("TextField subscribe: %v", err)
	}
	if w == nil {
		b.Fatal("TextField did not emit a widget")
	}

	r := new(gioinput.Router)
	ops := new(op.Ops)
	size := image.Pt(300, 60)

	// Frame 1 registers the editor's pointer/keyboard regions; we need its
	// dimensions to click at the centre.
	dims := driveTextFieldFrame(w, ops, r, size)
	centre := f32.Pt(float32(dims.Size.X)/2, float32(dims.Size.Y)/2)
	r.Queue(
		pointer.Event{Kind: pointer.Press, Position: centre, Buttons: pointer.ButtonPrimary, Source: pointer.Mouse},
		pointer.Event{Kind: pointer.Release, Position: centre, Buttons: pointer.ButtonPrimary, Source: pointer.Mouse},
	)
	// Frame 2 lets the editor process the press and request focus.
	driveTextFieldFrame(w, ops, r, size)

	// Type into the now-focused editor so the caret sits inside real content.
	r.Queue(key.EditEvent{Range: key.Range{Start: 0, End: 0}, Text: "caret"})
	driveTextFieldFrame(w, ops, r, size)

	if changed != "caret" {
		b.Fatalf("editor did not focus/ingest input (OnChange=%q): caret path is dead, benchmark would measure the wrong frame", changed)
	}

	bench.BenchFrame(b, w, bench.WithRouter(r), bench.WithSize(size))
}
