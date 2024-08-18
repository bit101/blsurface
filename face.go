// Package blsurface implements a 3d grid surface.
package blsurface

import (
	cairo "github.com/bit101/blcairo"
)

// Face holds a single face of the surface (four points).
type Face struct {
	p0, p1, p2, p3 *GridPoint
}

// NewFace creates a new face from 4 points.
func NewFace(p0, p1, p2, p3 *GridPoint) *Face {
	return &Face{p0, p1, p2, p3}
}

// Zpos is the average z position of this face.
func (f *Face) Zpos() float64 {
	return (f.p0.Z + f.p1.Z + f.p2.Z + f.p3.Z) / 4
}

// Draw draws a single face
func (f *Face) Draw(context *cairo.Context, scale float64, colorFunc ColorFunc) {
	avg := &GridPoint{
		X: (f.p0.origX + f.p1.origX + f.p2.origX + f.p3.origX) / 4,
		Y: (f.p0.origY + f.p1.origY + f.p2.origY + f.p3.origY) / 4,
		Z: (f.p0.origZ + f.p1.origZ + f.p2.origZ + f.p3.origZ) / 4,
	}

	context.Save()
	context.MoveTo(f.project(f.p0, scale))
	context.LineTo(f.project(f.p1, scale))
	context.LineTo(f.project(f.p2, scale))
	context.LineTo(f.project(f.p3, scale))
	context.ClosePath()

	context.SetSourceColor(colorFunc(avg.X, avg.Y, avg.Z))
	context.FillPreserve()

	context.Save()
	context.SetLineWidth(1)
	context.StrokePreserve()
	context.Restore()
	context.SetSourceBlack()
	context.Stroke()
	context.Restore()
}

func (f *Face) project(p *GridPoint, scale float64) (float64, float64) {
	fl := 300.0
	x := p.X * scale
	y := p.Y * scale
	z := p.Z * scale
	scale = fl / (fl + z + 200)
	x *= scale
	y *= scale
	return x, y
}
