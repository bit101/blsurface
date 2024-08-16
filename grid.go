// Package blsurface implements a 3d grid surface.
package blsurface

import (
	"log"
	"math"

	"github.com/bit101/bitlib/blcolor"
	"github.com/bit101/bitlib/blmath"
	cairo "github.com/bit101/blcairo"
)

// YFunction is the definition for a function that returns a y value for given x, z values.
type YFunction func(x, z float64) float64

// ColorFunc is the definition of a function that returns a color based on a 3d point.
type ColorFunc func(x, y, z float64) blcolor.Color

// Grid represents a grid of grid points.
type Grid struct {
	cells                  []*GridPoint
	w, d                   int
	rotation               float64
	tilt                   float64
	yFunc                  YFunction
	colorFunc              ColorFunc
	xMin, zMin, xMax, zMax float64
	yScale                 float64
	width                  float64
}

// NewGrid creates a new grid.
func NewGrid() *Grid {

	return &Grid{
		w:         20,
		d:         20,
		rotation:  math.Pi / 6,
		tilt:      math.Pi / 6,
		yFunc:     func(x, z float64) float64 { return 0.0 },
		colorFunc: func(x, y, z float64) blcolor.Color { return blcolor.White },
		yScale:    1.0,
		xMin:      -1,
		zMin:      -1,
		xMax:      1,
		zMax:      1,
		width:     400,
	}
}

func (g *Grid) SetXRange(xMin, xMax float64) {
	g.xMin = xMin
	g.xMax = xMax
	xRange := g.xMax - g.xMin
	zRange := g.zMax - g.zMin
	g.d = int(float64(g.w) * zRange / xRange)
}

func (g *Grid) SetZRange(zMin, zMax float64) {
	g.zMin = zMin
	g.zMax = zMax
	xRange := g.xMax - g.xMin
	zRange := g.zMax - g.zMin
	g.d = int(float64(g.w) * zRange / xRange)
}

func (g *Grid) SetYScale(yScale float64) {
	g.yScale = yScale
}

func (g *Grid) SetGridSize(gridSize int) {
	xRange := g.xMax - g.xMin
	zRange := g.zMax - g.zMin
	g.w = gridSize
	g.d = int(float64(g.w) * zRange / xRange)
}

// SetYFunc sets the function that computes the y value for a given x and z.
func (g *Grid) SetYFunc(yFunc YFunction) {
	g.yFunc = yFunc
}

// SetColorFunc sets the function that computes the color for a given x, y, z.
func (g *Grid) SetColorFunc(colorFunc ColorFunc) {
	g.colorFunc = colorFunc
}

// SetTilt rotates all the points on the x axis.
func (g *Grid) SetTilt(t float64) {
	if t < -math.Pi/2 || t > math.Pi/2 {
		log.Fatal("tilt must be between -90 and +90 degrees (-PI/2 to PI/2 radians)")
		// this can be removed once I figure out how to fix the hidden surface removal
		// when tilt is outside this range.
	}
	g.tilt = t
}

// SetTiltDegrees rotates all the points on the x axis.
func (g *Grid) SetTiltDegrees(t float64) {
	g.SetTilt(t / 180.0 * math.Pi)
}

// SetRotation rotates all the points on the y axis.
func (g *Grid) SetRotation(t float64) {
	for t < 0 {
		t += blmath.Tau
	}
	for t > blmath.Tau {
		t -= blmath.Tau
	}
	g.rotation = t
}

// SetRotationDegrees rotates all the points on the y axis.
func (g *Grid) SetRotationDegrees(t float64) {
	g.SetRotation(t / 180.0 * math.Pi)
}

// SetWidth sets the width of the graph on the x axis.
func (g *Grid) SetWidth(w float64) {
	g.width = w
}

// DrawPoints draws each point in the grid.
func (g *Grid) DrawPoints(context *cairo.Context, radius float64) {
	for _, p := range g.cells {
		context.FillCircle(p.X, p.Y, radius)
	}
}

func (g *Grid) makeGrid() {
	wf := float64(g.w)
	df := float64(g.d)
	grid := []*GridPoint{}
	for z := range g.d + 1 {
		zf := blmath.Map(float64(z), 0, df, g.zMin, g.zMax)
		for x := range g.w + 1 {
			xf := blmath.Map(float64(x), 0, wf, g.xMin, g.xMax)
			yf := 0.0
			p := NewGridPoint(xf, yf, zf)
			grid = append(grid, p)
		}
	}
	g.cells = grid
}

// DrawCells draws the whole grid.
func (g *Grid) DrawCells(context *cairo.Context) {
	g.makeGrid()
	g.applyFunc()
	for i := 0; i < g.w; i++ {
		x := i
		if g.rotation <= math.Pi {
			x = g.w - 1 - i
		}
		for j := 0; j < g.d; j++ {
			z := j
			if g.rotation > math.Pi*0.5 && g.rotation <= math.Pi*1.5 {
				z = g.d - 1 - j
			}
			g.drawCell(context, x, z)
		}
	}
}

func (g *Grid) getCell(x, z int) *GridPoint {
	index := z*(g.w+1) + x
	return g.cells[index]
}

func (g *Grid) applyFunc() {
	for _, p := range g.cells {
		p.Y = g.yFunc(p.X, p.Z) * g.yScale
		p.origY = p.Y
	}
	g.transform()
}

func (g *Grid) transform() {
	xRange := g.xMax - g.xMin
	zRange := g.zMax - g.zMin
	for _, p := range g.cells {
		p.X = blmath.Map(p.X, g.xMin, g.xMax, -xRange/2, xRange/2)
		p.Z = blmath.Map(p.Z, g.zMin, g.zMax, -zRange/2, zRange/2)
	}
	for _, p := range g.cells {
		p.RotateY(g.rotation)
	}
	for _, p := range g.cells {
		p.RotateX(g.tilt)
	}
}

func (g *Grid) drawCell(context *cairo.Context, x, z int) {
	p0 := g.getCell(x, z)
	p1 := g.getCell(x+1, z)
	p2 := g.getCell(x+1, z+1)
	p3 := g.getCell(x, z+1)
	avg := &GridPoint{
		X: (p0.origX + p1.origX + p2.origX + p3.origX) / 4,
		Y: (p0.origY + p1.origY + p2.origY + p3.origY) / 4,
		Z: (p0.origZ + p1.origZ + p2.origZ + p3.origZ) / 4,
	}
	context.Save()

	context.MoveTo(g.project(p0))
	context.LineTo(g.project(p1))
	context.LineTo(g.project(p2))
	context.LineTo(g.project(p3))
	context.ClosePath()

	context.SetSourceColor(g.colorFunc(avg.X, avg.Y, avg.Z))
	context.FillPreserve()

	context.Save()
	context.SetLineWidth(1)
	context.StrokePreserve()
	context.Restore()
	context.SetSourceBlack()
	context.Stroke()
	context.Restore()
}

func (g *Grid) project(p *GridPoint) (float64, float64) {
	// someday I'll get perspective working again.
	// fl := 500.0
	// scale := fl / (fl - p.Z + 300)
	scale := g.width / (g.xMax - g.xMin)
	// fmt.Println(scale)
	return p.X * scale, p.Y * scale
}
