// Package blsurface implements a 3d grid surface.
package blsurface

import (
	"math"

	"github.com/bit101/bitlib/blmath"
	cairo "github.com/bit101/blcairo"
)

// YFunction is the definition for a function that returns a y value for given x, z values.
type YFunction func(x, z float64) float64

// Grid represents a grid of grid points.
type Grid struct {
	Grid      []*GridPoint
	W, D      int
	yrotation float64
}

// NewGrid creates a new grid.
func NewGrid(w, d int, spacing float64, yfunc YFunction) *Grid {
	wf := float64(w)
	df := float64(d)
	halfW := wf * spacing / 2
	halfD := df * spacing / 2
	grid := []*GridPoint{}
	for z := range d {
		zf := blmath.Map(float64(z), 0, df, -halfD, halfD)
		for x := range w {
			xf := blmath.Map(float64(x), 0, wf, -halfW, halfW)
			yf := yfunc(xf, zf)
			p := GridPoint{xf, yf, zf}
			grid = append(grid, &p)
		}
	}
	return &Grid{grid, w, d, 0}
}

// Tilt rotates all the points on the x axis.
func (g *Grid) Tilt(t float64) {
	for _, p := range g.Grid {
		p.RotateX(t)
	}
}

// Rotate rotates all the points on the y axis.
func (g *Grid) Rotate(t float64) {
	for t < 0 {
		t += blmath.Tau
	}
	g.yrotation = t
	for _, p := range g.Grid {
		p.RotateY(t)
	}
}

// // RotateZ rotates all the points on the z axis.
// func (g *Grid) RotateZ(t float64) {
// 	for _, p := range g.Grid {
// 		p.RotateZ(t)
// 	}
// }

// DrawPoints draws each point in the grid.
func (g *Grid) DrawPoints(context *cairo.Context, radius float64) {
	for _, p := range g.Grid {
		context.FillCircle(p.X, p.Y, radius)
	}
}

// DrawCells draws the whole grid.
func (g *Grid) DrawCells(context *cairo.Context) {
	for i := 0; i < g.W-1; i++ {
		xIndex := i
		if g.yrotation > 0 && g.yrotation <= math.Pi {
			xIndex = g.W - 2 - i
		}
		for j := 0; j < g.D-1; j++ {
			zIndex := g.W * j
			if g.yrotation > math.Pi*0.5 && g.yrotation <= math.Pi*1.5 {
				zIndex = (g.D - 2 - j) * g.W
			}
			g.drawCell(context, xIndex+zIndex)
		}
	}
}

func (g *Grid) drawCell(context *cairo.Context, index int) {
	p0 := g.Grid[index]
	p1 := g.Grid[index+g.W]
	p2 := g.Grid[index+g.W+1]
	p3 := g.Grid[index+1]

	context.MoveTo(project(p0))
	context.LineTo(project(p1))
	context.LineTo(project(p2))
	context.LineTo(project(p3))
	context.ClosePath()

	context.SetSourceWhite()
	context.FillPreserve()

	context.SetSourceBlack()
	context.Stroke()
}

func project(p *GridPoint) (float64, float64) {
	// fl := 500.0
	// scale := fl / (fl - p.Z + 300)
	scale := 1.0
	// fmt.Println(scale)
	return p.X * scale, p.Y * scale
}
