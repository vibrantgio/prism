package input_test

import (
	"image"
	"testing"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/vibrantgio/prism/input"
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
