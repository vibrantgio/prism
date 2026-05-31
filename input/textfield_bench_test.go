package input_test

import (
	"image"
	"testing"

	"gioui.org/font/gofont"
	gioinput "gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/reactivego/rx"
	"github.com/vibrantgio/prism/input"
	"github.com/vibrantgio/prism/theme"
	"github.com/vibrantgio/prism/tokens"
)

// BenchmarkTextFieldRender exercises widget(gtx) for b.N synthetic frames,
// per DESIGN §"Performance — Profiling". b.ReportAllocs is enabled so CI
// can gate on per-frame allocation regressions (>5% threshold).
func BenchmarkTextFieldRender(b *testing.B) {
	shaper := text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	w := input.Render(
		shaper, "Placeholder",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		input.RenderState{},
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var ops op.Ops
		gtx := layout.Context{
			Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
			Constraints: layout.Exact(image.Pt(300, 60)),
			Ops:         &ops,
		}
		w(gtx)
	}
}

// BenchmarkTextFieldRenderFocused benchmarks the focused state which draws a
// wider border stroke.
func BenchmarkTextFieldRenderFocused(b *testing.B) {
	shaper := text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	w := input.Render(
		shaper, "Placeholder",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		input.RenderState{Focused: true},
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var ops op.Ops
		gtx := layout.Context{
			Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
			Constraints: layout.Exact(image.Pt(300, 60)),
			Ops:         &ops,
		}
		w(gtx)
	}
}

// BenchmarkTextFieldRenderSubmit exercises the live submit-enabled TextField
// path so per-frame allocations on chat-style inputs are measurable. The
// editor and event loop are taken through a real input.Router; submit
// configures the editor to translate Enter into widget.SubmitEvent, which
// changes the event-handling branch but not the visible rendering.
func BenchmarkTextFieldRenderSubmit(b *testing.B) {
	shaper := text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	props := input.TextFieldProps{
		Placeholder:   "Placeholder",
		Submit:        true,
		SubmitMessage: func(s string) any { return s },
		OnSubmit:      func(_ layout.Context, _ string) {},
		Shaper:        shaper,
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
	var ops op.Ops

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ops.Reset()
		gtx := layout.Context{
			Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
			Constraints: layout.Exact(image.Pt(300, 60)),
			Ops:         &ops,
			Source:      r.Source(),
		}
		w(gtx)
		r.Frame(&ops)
	}
}
