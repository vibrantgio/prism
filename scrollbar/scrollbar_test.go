package scrollbar

import (
	"image/color"
	"testing"

	"gioui.org/unit"

	"github.com/vibrantgio/prism/tokens"
)

func TestFromTokens(t *testing.T) {
	cases := []struct {
		name string
		c    tokens.ColorTokens
	}{
		{"DefaultLight", tokens.DefaultLight},
		{"DefaultDark", tokens.DefaultDark},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := FromTokens(tc.c)

			wantThumb := tc.c.OnSurfaceVariant
			wantThumb.A = 100
			if s.ThumbColor != wantThumb {
				t.Errorf("ThumbColor = %v, want %v", s.ThumbColor, wantThumb)
			}

			wantHover := tc.c.OnSurfaceVariant
			wantHover.A = 170
			if s.ThumbHoverColor != wantHover {
				t.Errorf("ThumbHoverColor = %v, want %v", s.ThumbHoverColor, wantHover)
			}

			if s.TrackColor != (color.NRGBA{}) {
				t.Errorf("TrackColor = %v, want transparent zero value", s.TrackColor)
			}

			metrics := []struct {
				name string
				got  unit.Dp
				want unit.Dp
			}{
				{"ThumbMinorWidth", s.ThumbMinorWidth, 6},
				{"TrackPadding", s.TrackPadding, 2},
				{"ThumbCornerRadius", s.ThumbCornerRadius, 3},
				{"ThumbMinLen", s.ThumbMinLen, 16},
				{"Width()", s.Width(), 10},
			}
			for _, m := range metrics {
				if m.got != m.want {
					t.Errorf("%s = %v, want %v", m.name, m.got, m.want)
				}
			}
		})
	}
}
