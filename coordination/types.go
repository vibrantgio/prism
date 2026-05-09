// Package coordination provides the Subject primitive for cross-widget
// coordination in Gio applications.
//
// The package codifies four invariants established by Experiments C1–C2
// (see EXPERIMENT-C.md):
//
//  1. One-frame lag. Subject delivery is asynchronous. Cross-widget state
//     changes are visible on the frame AFTER the emitting frame. This is
//     correct for drag, modal, and tooltip concerns and is imperceptible at
//     ≥30 fps.
//
//  2. Mutable per-widget state must be hoisted outside FRP closures. Gesture
//     accumulators (gesture.Drag, gesture.Hover, widget.Clickable) must live
//     in the owning struct, not inside rx.Map closures that are regenerated on
//     every Subject emission.
//
//  3. Buffer capacity must exceed maximum burst. For pointer-event emitters,
//     BufCapPointer prevents frame-goroutine blocking. For infrequent signals
//     (modals, tooltips) BufCapSignal suffices.
//
//  4. Intermediate emissions are silently dropped under burst. The
//     mvu.Window atomic snapshot retains only the most recent widget closure
//     before the next frame fires. Signals where every value is load-bearing
//     (undo stacks, event logs) require a different mechanism.
package coordination

import "github.com/reactivego/rx"

// BufCapPointer is the recommended producer-side buffer depth for Subjects
// that emit on every pointer event (drag, hover). Sized to ~2×60 fps so the
// frame goroutine is never blocked under burst pointer events.
const BufCapPointer = 128

// BufCapSignal is the recommended producer-side buffer depth for infrequent
// coordination signals (modal depth, focus owner, tooltip arbitration).
const BufCapSignal = 8

// Subject creates a typed broadcast channel for cross-widget coordination.
// The Observer side is held by one producer; the Observable side may be
// subscribed by up to eight concurrent consumers.
//
// bufCap is the producer-side buffer depth. Use BufCapPointer for signals
// emitted on every pointer event; use BufCapSignal for infrequent signals.
func Subject[T any](bufCap int) (rx.Observer[T], rx.Observable[T]) {
	return rx.Subject[T](0, 0, bufCap, 8)
}
