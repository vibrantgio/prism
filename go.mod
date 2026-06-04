module github.com/vibrantgio/prism

go 1.25.1

require (
	gioui.org v0.9.0
	github.com/reactivego/rx v0.2.2
	github.com/vibrantgio/ivg/raster/gio v0.0.0
	github.com/vibrantgio/mvu v0.0.0
	github.com/vibrantgio/pulse v0.0.0
	github.com/vibrantgio/svg v0.0.3
	github.com/vibrantgio/svg/driver/gio v0.0.0
)

require (
	gioui.org/shader v1.0.8 // indirect
	github.com/go-text/typesetting v0.3.0 // indirect
	github.com/reactivego/scheduler v0.1.2 // indirect
	github.com/vibrantgio/ivg v0.1.3 // indirect
	github.com/vibrantgio/traer v0.0.4 // indirect
	golang.org/x/exp/shiny v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/image v0.26.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace (
	github.com/vibrantgio/ivg/raster/gio => ../ivg/raster/gio
	github.com/vibrantgio/mvu => ../mvu
	github.com/vibrantgio/pulse => ../pulse
	// prism/icon's local svg/driver/gio pairs with the local svg (whose driver
	// interface lacks the SetupDrawers method that published v0.0.3 added);
	// replace svg too so the standalone (GOWORK=off) icon build matches the
	// workspace, which already resolves svg via its ./svg use directive.
	github.com/vibrantgio/svg => ../svg
	github.com/vibrantgio/svg/driver/gio => ../svg/driver/gio
	// Local traer is ahead of published v0.0.4 (adds NewParticleSystem, used by
	// pulse, which prism/gallery composes); replace to the sibling repo so the
	// standalone (GOWORK=off) gallery build resolves it.
	github.com/vibrantgio/traer => ../traer
)
