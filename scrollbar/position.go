package scrollbar

import "gioui.org/layout"

// FromListPosition converts a layout.Position into two floats representing
// the location of the viewport on the underlying content. It needs to know
// the number of elements in the list and the major-axis size of the list in
// order to do this. The returned values are in the range [0,1] with
// start <= end. If there is nothing to scroll (no elements or a zero-length
// position), it returns (0, 1).
//
// This is a port of material's unexported fromListPosition
// (gioui.org/widget/material/list.go).
func FromListPosition(lp layout.Position, elements int, majorAxisSize int) (start, end float32) {
	if elements == 0 || lp.Length == 0 {
		return 0, 1
	}

	// Approximate the size of the scrollable content.
	lengthEstPx := float32(lp.Length)
	elementLenEstPx := lengthEstPx / float32(elements)

	// Determine how much of the content is visible.
	listOffsetF := float32(lp.Offset)
	listOffsetL := float32(lp.OffsetLast)

	// Compute the location of the beginning and end of the viewport using
	// the estimated element size and known pixel offsets.
	viewportStart := clamp1((float32(lp.First)*elementLenEstPx + listOffsetF) / lengthEstPx)
	viewportEnd := clamp1((float32(lp.First+lp.Count)*elementLenEstPx + listOffsetL) / lengthEstPx)
	viewportFraction := viewportEnd - viewportStart

	// Compute the expected visible proportion of the list content based
	// solely on the ratio of the visible size and the estimated total size.
	visiblePx := float32(majorAxisSize)
	visibleFraction := visiblePx / lengthEstPx

	// Compute the error between the two methods of determining the viewport
	// and diffuse the error on either end of the viewport based on how close
	// we are to each end.
	err := visibleFraction - viewportFraction
	adjStart := viewportStart
	adjEnd := viewportEnd
	if viewportFraction < 1 {
		startShare := viewportStart / (1 - viewportFraction)
		endShare := (1 - viewportEnd) / (1 - viewportFraction)
		startErr := startShare * err
		endErr := endShare * err

		adjStart -= startErr
		adjEnd += endErr
	}
	return adjStart, adjEnd
}

// clamp1 limits v to the range [0,1].
func clamp1(v float32) float32 {
	if v >= 1 {
		return 1
	} else if v <= 0 {
		return 0
	}
	return v
}
