// Package blsurface implements a 3d grid surface.
package blsurface

import "math"

// GridPoint is a 3d point on a grid.
// Each point contains X, Y, Z values which are tranformed versions of the original values,
// and origX, origY, origZ, which are not transformed.
// The color algorithm needs the untransformed values.
// Possible to make this immutable and have transform functions return copies.
// But this will add other complications. Look into it.
type GridPoint struct {
	origX, origY, origZ float64
	X, Y, Z             float64
}

// NewGridPoint creates a new GridPoint.
func NewGridPoint(x, y, z float64) *GridPoint {
	return &GridPoint{x, y, z, x, y, z}
}

// RotateX rotates a grid point on the x axis.
func (g *GridPoint) RotateX(t float64) {
	cos := math.Cos(t)
	sin := math.Sin(t)
	y := g.Y*cos - g.Z*sin
	z := g.Z*cos + g.Y*sin
	g.Y = y
	g.Z = z
}

// RotateY rotates a grid point on the y axis.
func (g *GridPoint) RotateY(t float64) {
	cos := math.Cos(t)
	sin := math.Sin(t)
	x := g.X*cos - g.Z*sin
	z := g.Z*cos + g.X*sin
	g.X = x
	g.Z = z
}

// RotateZ rotates a grid point on the z axis.
func (g *GridPoint) RotateZ(t float64) {
	cos := math.Cos(t)
	sin := math.Sin(t)
	x := g.X*cos - g.Y*sin
	y := g.Y*cos + g.X*sin
	g.X = x
	g.Y = y
}
