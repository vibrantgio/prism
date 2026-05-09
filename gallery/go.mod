module github.com/vibrantgio/prism/gallery

go 1.25.1

require (
	gioui.org v0.9.0
	github.com/reactivego/rx v0.2.2
	github.com/vibrantgio/ivg/raster/gio v0.0.0
	github.com/vibrantgio/prism/a11y v0.0.0
	github.com/vibrantgio/prism/button v0.0.0
	github.com/vibrantgio/prism/coordination v0.0.0
	github.com/vibrantgio/prism/initial v0.0.0
	github.com/vibrantgio/prism/input v0.0.0
	github.com/vibrantgio/prism/layout v0.0.0
	github.com/vibrantgio/prism/list v0.0.0
	github.com/vibrantgio/prism/theme v0.0.0
	github.com/vibrantgio/prism/tokens v0.0.0
)

require (
	gioui.org/shader v1.0.8 // indirect
	github.com/go-text/typesetting v0.3.0 // indirect
	github.com/reactivego/scheduler v0.1.2 // indirect
	github.com/vibrantgio/ivg v0.1.3 // indirect
	github.com/vibrantgio/mvu v0.0.0 // indirect
	golang.org/x/exp/shiny v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/image v0.26.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace (
	github.com/vibrantgio/ivg/raster/gio => ../../ivg/raster/gio
	github.com/vibrantgio/mvu => ../../mvu
	github.com/vibrantgio/prism/a11y => ../a11y
	github.com/vibrantgio/prism/button => ../button
	github.com/vibrantgio/prism/coordination => ../coordination
	github.com/vibrantgio/prism/icon => ../icon
	github.com/vibrantgio/prism/initial => ../initial
	github.com/vibrantgio/prism/input => ../input
	github.com/vibrantgio/prism/layout => ../layout
	github.com/vibrantgio/prism/list => ../list
	github.com/vibrantgio/prism/theme => ../theme
	github.com/vibrantgio/prism/tokens => ../tokens
	github.com/vibrantgio/svg/driver/gio => ../../svg/driver/gio
)
