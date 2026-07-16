package scrollbar

import (
	"image"
	"math"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Layout draws the scrollbar along axis and registers its gesture areas.
// viewportStart and viewportEnd describe the visible fraction of the content
// in the range [0,1] (see FromListPosition).
//
// The bar renders whenever the viewport shows less than the full content and
// renders nothing (zero dimensions) when everything fits. It occupies the
// full major axis of the incoming constraints and Width() along the minor
// axis.
func (s Style) Layout(gtx layout.Context, state *State, axis layout.Axis, viewportStart, viewportEnd float32) layout.Dimensions {
	if viewportStart <= 0 && viewportEnd >= 1 {
		// Everything fits: no scrollbar.
		return layout.Dimensions{}
	}

	// Pin the constraints in an axis-independent way, then convert to the
	// correct representation for the current axis.
	convert := axis.Convert
	maxMajorAxis := convert(gtx.Constraints.Max).X
	gtx.Constraints.Min.X = maxMajorAxis
	gtx.Constraints.Min.Y = gtx.Dp(s.Width())
	gtx.Constraints.Min = convert(gtx.Constraints.Min)
	gtx.Constraints.Max = gtx.Constraints.Min

	// Process events against last frame's areas before reading hover state.
	state.Update(gtx, axis, viewportStart, viewportEnd)

	thumbColor := s.ThumbColor
	if state.IndicatorHovered() || state.Dragging() {
		thumbColor = s.ThumbHoverColor
	}

	inset := layout.Inset{
		Top:    s.TrackPadding,
		Bottom: s.TrackPadding,
		Left:   s.TrackPadding,
		Right:  s.TrackPadding,
	}

	return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			// Lay out the draggable track underneath the thumb.
			area := image.Rectangle{Max: gtx.Constraints.Min}
			pointerArea := clip.Rect(area)
			defer pointerArea.Push(gtx.Ops).Pop()
			state.AddDrag(gtx.Ops)

			// Stack a normal clickable area on top of the draggable area
			// to capture non-dragging clicks.
			defer pointer.PassOp{}.Push(gtx.Ops).Pop()
			defer pointerArea.Push(gtx.Ops).Pop()
			state.AddTrack(gtx.Ops)

			paint.FillShape(gtx.Ops, s.TrackColor, clip.Rect(area).Op())
			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
		func(gtx layout.Context) layout.Dimensions {
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Work in axis-independent space.
				gtx.Constraints.Min = convert(gtx.Constraints.Min)
				gtx.Constraints.Max = convert(gtx.Constraints.Max)

				// Compute the pixel size and position of the thumb
				// within the track.
				trackLen := gtx.Constraints.Min.X
				viewStart := int(math.Round(float64(viewportStart) * float64(trackLen)))
				viewEnd := int(math.Round(float64(viewportEnd) * float64(trackLen)))
				thumbLen := max(viewEnd-viewStart, gtx.Dp(s.ThumbMinLen))
				if viewStart+thumbLen > trackLen {
					viewStart = trackLen - thumbLen
				}
				thumbDims := convert(image.Point{
					X: thumbLen,
					Y: gtx.Dp(s.ThumbMinorWidth),
				})
				radius := gtx.Dp(s.ThumbCornerRadius)

				// Draw the thumb.
				offset := convert(image.Pt(viewStart, 0))
				defer op.Offset(offset).Push(gtx.Ops).Pop()
				paint.FillShape(gtx.Ops, thumbColor, clip.RRect{
					Rect: image.Rectangle{Max: thumbDims},
					SW:   radius,
					NW:   radius,
					NE:   radius,
					SE:   radius,
				}.Op(gtx.Ops))

				// Register the thumb's pointer hit area.
				area := clip.Rect(image.Rectangle{Max: thumbDims})
				defer pointer.PassOp{}.Push(gtx.Ops).Pop()
				defer area.Push(gtx.Ops).Pop()
				state.AddIndicator(gtx.Ops)

				return layout.Dimensions{Size: convert(gtx.Constraints.Min)}
			})
		},
	)
}
