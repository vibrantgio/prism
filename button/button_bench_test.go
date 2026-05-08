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
	"github.com/vibrantgio/prism/tokens"
)

// BenchmarkButtonRender exercises widget(gtx) for b.N synthetic frames,
// per DESIGN §"Performance — Profiling". b.ReportAllocs is enabled so CI
// can gate on per-frame allocation regressions (>5% threshold).
func BenchmarkButtonRender(b *testing.B) {
	shaper := text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	w := button.Render(
		shaper, "Benchmark",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{},
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

// BenchmarkButtonRenderFocused benchmarks the focused state which draws an
// extra clip.Stroke path for the focus ring.
func BenchmarkButtonRenderFocused(b *testing.B) {
	shaper := text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	w := button.Render(
		shaper, "Benchmark",
		tokens.DefaultLight, tokens.Spacing, tokens.Radius, tokens.DefaultTypeScale,
		button.RenderState{Focused: true},
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
