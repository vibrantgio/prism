package tokens

// ElevationScale holds shadow-depth stops in device-independent pixels,
// following Material Design 3 elevation levels 0–5.
type ElevationScale struct {
	Level0 float32 // 0 dp
	Level1 float32 // 1 dp
	Level2 float32 // 3 dp
	Level3 float32 // 6 dp
	Level4 float32 // 8 dp
	Level5 float32 // 12 dp
}

// Elevation is the default scale instance.
var Elevation = ElevationScale{
	Level0: 0,
	Level1: 1,
	Level2: 3,
	Level3: 6,
	Level4: 8,
	Level5: 12,
}
