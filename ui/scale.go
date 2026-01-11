package ui

import "math"

// Scale converts logical units to physical pixels.
type Scale struct {
	Device float64 // device scale factor (e.g. 2.0 on HiDPI)
	UI     float64 // additional UI scale (user preference), default 1.0
}

func (s Scale) Factor() float64 {
	f := s.Device * s.UI
	if f <= 0 {
		return 1
	}
	return f
}

func (s Scale) PxF(v float64) float64 { return v * s.Factor() }
func (s Scale) PxI(v int) int         { return int(math.Round(float64(v) * s.Factor())) }

// Snap rounds a physical-pixel coordinate to an integer pixel.
func (s Scale) Snap(v float64) float64 { return math.Round(v) }
