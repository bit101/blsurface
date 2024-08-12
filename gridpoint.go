// Package blsurface implements a 3d grid surface.
package blsurface

import "math"

// GridPoint is a 3d point on a grid.
type GridPoint struct {
	X, Y, Z float64
}

// RotateX rotates a grid point on the x axis.
func (g *GridPoint) RotateX(t float64) {
	cos := math.Cos(t)
	sin := math.Sin(t)
	y := g.Y*cos + g.Z*sin
	z := g.Z*cos - g.Y*sin
	g.Y = y
	g.Z = z
}

// RotateY rotates a grid point on the y axis.
func (g *GridPoint) RotateY(t float64) {
	cos := math.Cos(t)
	sin := math.Sin(t)
	x := g.X*cos + g.Z*sin
	z := g.Z*cos - g.X*sin
	g.X = x
	g.Z = z
}

// RotateZ rotates a grid point on the z axis.
func (g *GridPoint) RotateZ(t float64) {
	cos := math.Cos(t)
	sin := math.Sin(t)
	x := g.X*cos + g.Y*sin
	y := g.Y*cos - g.X*sin
	g.X = x
	g.Y = y
}
