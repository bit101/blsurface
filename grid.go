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
	cells     []*GridPoint
	w, d      int
	yrotation float64
	tilt      float64
	yFunc     YFunction
}

// NewGrid creates a new grid.
func NewGrid(w, d int, spacing float64) *Grid {
	wf := float64(w)
	df := float64(d)
	halfW := wf * spacing / 2
	halfD := df * spacing / 2
	grid := []*GridPoint{}
	for z := range d {
		zf := blmath.Map(float64(z), 0, df, -halfD, halfD)
		for x := range w {
			xf := blmath.Map(float64(x), 0, wf, -halfW, halfW)
			yf := 0.0
			p := GridPoint{xf, yf, zf}
			grid = append(grid, &p)
		}
	}
	return &Grid{grid, w, d, 0, 0, func(x, z float64) float64 { return 0.0 }}
}

// SetYFunc sets the function that computes the y value for a given x and z.
func (g *Grid) SetYFunc(yFunc YFunction) {
	g.yFunc = yFunc
}

// SetTilt rotates all the points on the x axis.
func (g *Grid) SetTilt(t float64) {
	g.tilt = t
}

// SetRotation rotates all the points on the y axis.
func (g *Grid) SetRotation(t float64) {
	for t < 0 {
		t += blmath.Tau
	}
	g.yrotation = t
}

// // RotateZ rotates all the points on the z axis.
// func (g *Grid) RotateZ(t float64) {
// 	for _, p := range g.Grid {
// 		p.RotateZ(t)
// 	}
// }

// DrawPoints draws each point in the grid.
func (g *Grid) DrawPoints(context *cairo.Context, radius float64) {
	for _, p := range g.cells {
		context.FillCircle(p.X, p.Y, radius)
	}
}

// DrawCells draws the whole grid.
func (g *Grid) DrawCells(context *cairo.Context) {
	g.applyFunc()
	for i := 0; i < g.w-1; i++ {
		xIndex := i
		if g.yrotation > 0 && g.yrotation <= math.Pi {
			xIndex = g.w - 2 - i
		}
		for j := 0; j < g.d-1; j++ {
			zIndex := g.w * j
			if g.yrotation > math.Pi*0.5 && g.yrotation <= math.Pi*1.5 {
				zIndex = (g.d - 2 - j) * g.w
			}
			g.drawCell(context, xIndex+zIndex)
		}
	}
}

func (g *Grid) applyFunc() {
	for _, p := range g.cells {
		p.Y = g.yFunc(p.X, p.Z)
	}
	g.transform()
}

func (g *Grid) transform() {
	for _, p := range g.cells {
		p.RotateY(g.yrotation)
	}
	for _, p := range g.cells {
		p.RotateX(g.tilt)
	}
}

func (g *Grid) drawCell(context *cairo.Context, index int) {
	p0 := g.cells[index]
	p1 := g.cells[index+g.w]
	p2 := g.cells[index+g.w+1]
	p3 := g.cells[index+1]

	context.MoveTo(project(p0))
	context.LineTo(project(p1))
	context.LineTo(project(p2))
	context.LineTo(project(p3))
	context.ClosePath()

	context.SetSourceWhite()
	context.FillPreserve()

	context.Save()
	context.SetLineWidth(1)
	context.SetSourceWhite()
	context.StrokePreserve()
	context.Restore()
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
