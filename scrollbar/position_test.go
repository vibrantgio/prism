package scrollbar

import (
	"testing"

	"gioui.org/layout"
)

func TestFromListPosition(t *testing.T) {
	// A uniform list of 100 elements, 20px each (Length 2000), viewed
	// through a 400px viewport that shows 20 elements at a time.
	const (
		elements = 100
		elemPx   = 20
		length   = elements * elemPx
		viewport = 400
	)

	tests := []struct {
		name          string
		lp            layout.Position
		elements      int
		majorAxisSize int
		wantStart     float32
		wantEnd       float32
		tol           float32
	}{
		{
			name:          "top of list",
			lp:            layout.Position{First: 0, Offset: 0, Count: 20, OffsetLast: 0, Length: length},
			elements:      elements,
			majorAxisSize: viewport,
			wantStart:     0,
			wantEnd:       0.2,
			tol:           1e-3,
		},
		{
			name:          "bottom of list",
			lp:            layout.Position{First: 80, Offset: 0, Count: 20, OffsetLast: 0, Length: length},
			elements:      elements,
			majorAxisSize: viewport,
			wantStart:     0.8,
			wantEnd:       1,
			tol:           1e-3,
		},
		{
			name:          "middle of uniform list",
			lp:            layout.Position{First: 40, Offset: 0, Count: 20, OffsetLast: 0, Length: length},
			elements:      elements,
			majorAxisSize: viewport,
			wantStart:     0.4,
			wantEnd:       0.6,
			tol:           1e-3,
		},
		{
			name:          "zero elements",
			lp:            layout.Position{Count: 0, Length: length},
			elements:      0,
			majorAxisSize: viewport,
			wantStart:     0,
			wantEnd:       1,
			tol:           0,
		},
		{
			name:          "zero length",
			lp:            layout.Position{First: 0, Count: 0, Length: 0},
			elements:      elements,
			majorAxisSize: viewport,
			wantStart:     0,
			wantEnd:       1,
			tol:           0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			start, end := FromListPosition(tc.lp, tc.elements, tc.majorAxisSize)
			if diff := start - tc.wantStart; diff < -tc.tol || diff > tc.tol {
				t.Errorf("start = %v, want %v (tol %v)", start, tc.wantStart, tc.tol)
			}
			if diff := end - tc.wantEnd; diff < -tc.tol || diff > tc.tol {
				t.Errorf("end = %v, want %v (tol %v)", end, tc.wantEnd, tc.tol)
			}
		})
	}

	t.Run("middle viewport fraction matches visible/total", func(t *testing.T) {
		lp := layout.Position{First: 40, Offset: 0, Count: 20, OffsetLast: 0, Length: length}
		start, end := FromListPosition(lp, elements, viewport)
		frac := end - start
		want := float32(viewport) / float32(length)
		if diff := frac - want; diff < -1e-3 || diff > 1e-3 {
			t.Errorf("end-start = %v, want %v within 1e-3", frac, want)
		}
	})

	t.Run("variable offsets keep invariant", func(t *testing.T) {
		lp := layout.Position{First: 3, Offset: 7, Count: 21, OffsetLast: -13, Length: length}
		start, end := FromListPosition(lp, elements, viewport)
		checkInvariant(t, lp, start, end)
	})
}

// checkInvariant asserts 0 <= start <= end <= 1.
func checkInvariant(t *testing.T, lp layout.Position, start, end float32) {
	t.Helper()
	if !(0 <= start && start <= end && end <= 1) {
		t.Errorf("invariant 0 <= start <= end <= 1 violated for %+v: start=%v end=%v", lp, start, end)
	}
}

// TestFromListPositionSweep exercises a few hundred synthetic positions and
// asserts the [0,1] ordering invariant everywhere, plus that start is
// non-decreasing as First grows with all other parameters fixed.
func TestFromListPositionSweep(t *testing.T) {
	const (
		elements = 100
		elemPx   = 20
		length   = elements * elemPx
		viewport = 400
	)

	for count := 18; count <= 22; count++ {
		for offset := 0; offset < elemPx; offset += 4 {
			for offsetLast := 0; offsetLast >= -elemPx; offsetLast -= 4 {
				prevStart := float32(-1)
				for first := 0; first+count <= elements; first += 2 {
					lp := layout.Position{
						First:      first,
						Offset:     offset,
						Count:      count,
						OffsetLast: offsetLast,
						Length:     length,
					}
					start, end := FromListPosition(lp, elements, viewport)
					checkInvariant(t, lp, start, end)
					if start < prevStart {
						t.Errorf("start not non-decreasing at %+v: start=%v < prev=%v", lp, start, prevStart)
					}
					prevStart = start
				}
			}
		}
	}
}
