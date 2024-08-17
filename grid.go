// Package blsurface implements a 3d grid surface.
package blsurface

import (
	"log"
	"math"

	"github.com/bit101/bitlib/blcolor"
	"github.com/bit101/bitlib/blmath"
	cairo "github.com/bit101/blcairo"
)

//////////////////////////////
// Types
//////////////////////////////

// YFunction is the definition for a function that returns a y value for given x, z values.
type YFunction func(x, z float64) float64

// ColorFunc is the definition of a function that returns a color based on a 3d point.
type ColorFunc func(x, y, z float64) blcolor.Color

// Grid represents a grid of grid points.
type Grid struct {
	cells                  []*GridPoint
	w, d                   int
	originX, originY       float64
	rotation               float64
	tilt                   float64
	yFunc                  YFunction
	colorFunc              ColorFunc
	xMin, zMin, xMax, zMax float64
	yScale                 float64
	width                  float64
	axes                   []*GridPoint
}

//////////////////////////////
// Constructor
//////////////////////////////

// NewGrid creates a new grid.
func NewGrid() *Grid {

	return &Grid{
		w:         20,
		d:         20,
		originX:   0,
		originY:   0,
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
		axes: []*GridPoint{
			NewGridPoint(0, 0, 0),
			NewGridPoint(1, 0, 0),
			NewGridPoint(0, -1, 0),
			NewGridPoint(0, 0, 1),
		},
	}
}

//////////////////////////////
// Drawing
//////////////////////////////

// DrawPoints draws each point in the grid.
func (g *Grid) DrawPoints(context *cairo.Context, radius float64) {
	for _, p := range g.cells {
		context.FillCircle(p.X, p.Y, radius)
	}
}

// DrawCells draws the whole grid.
func (g *Grid) DrawCells(context *cairo.Context) {
	g.makeGrid()
	g.applyFunc()
	g.transform()
	context.Save()
	context.Translate(g.originX, g.originY)
	for i := 0; i < g.w; i++ {
		for j := 0; j < g.d; j++ {
			x := g.getXIndex(i)
			z := g.getZIndex(j)
			g.drawCell(context, x, z)
		}
	}
	context.Restore()
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

func (g *Grid) DrawOrigin(context *cairo.Context, x, y, size float64) {
	context.Save()
	m := *cairo.NewMatrix()
	m.InitIdentity()
	context.SetMatrix(m)
	context.Translate(x, y)
	context.MoveTo(g.axes[0].X*size, g.axes[0].Y*size)
	context.LineTo(g.axes[1].X*size, g.axes[1].Y*size)
	context.MoveTo(g.axes[0].X*size, g.axes[0].Y*size)
	context.LineTo(g.axes[2].X*size, g.axes[2].Y*size)
	context.MoveTo(g.axes[0].X*size, g.axes[0].Y*size)
	context.LineTo(g.axes[3].X*size, g.axes[3].Y*size)
	context.Stroke()

	context.FillText("x", g.axes[1].X*size+5, g.axes[1].Y*size)
	context.FillText("y", g.axes[2].X*size+5, g.axes[2].Y*size)
	context.FillText("z", g.axes[3].X*size+5, g.axes[3].Y*size)

	context.Restore()
}

// ////////////////////////////
// Helpers
// ////////////////////////////
func (g *Grid) applyFunc() {
	for _, p := range g.cells {
		p.Y = g.yFunc(p.X, p.Z) * g.yScale
		p.origY = p.Y
	}
}

func (g *Grid) getCell(x, z int) *GridPoint {
	index := z*(g.w+1) + x
	return g.cells[index]
}

func (g *Grid) getXIndex(i int) int {
	if g.rotation <= math.Pi {
		return g.w - 1 - i
	}
	return i
}

func (g *Grid) getZIndex(j int) int {
	if g.rotation > math.Pi*0.5 && g.rotation <= math.Pi*1.5 {
		return g.d - 1 - j
	}
	return j
}

func (g *Grid) makeGrid() {
	wf := float64(g.w)
	df := float64(g.d)
	grid := []*GridPoint{}
	for z := range g.d + 1 {
		zf := blmath.Map(float64(z), 0, df, g.zMax, g.zMin)
		for x := range g.w + 1 {
			xf := blmath.Map(float64(x), 0, wf, g.xMin, g.xMax)
			yf := 0.0
			p := NewGridPoint(xf, yf, zf)
			grid = append(grid, p)
		}
	}
	g.cells = grid
}

func (g *Grid) project(p *GridPoint) (float64, float64) {
	// someday I'll get perspective working again.
	perspective := false
	scale := g.width / (g.xMax - g.xMin)
	x := p.X * scale
	y := p.Y * scale
	if perspective {
		fl := 300.0
		z := p.Z * scale
		scale = fl / (fl + z + 100)
		return x * scale, y * scale
	}
	return x, y
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
	for _, p := range g.axes {
		p.RotateY(g.rotation)
	}
	for _, p := range g.axes {
		p.RotateX(g.tilt)
	}
}

//////////////////////////////
// Setters
//////////////////////////////

// SetColorFunc sets the function that computes the color for a given x, y, z.
func (g *Grid) SetColorFunc(colorFunc ColorFunc) {
	g.colorFunc = colorFunc
}

// SetGridSize sets how many cells to draw across the x-axis.
// The number of cells for the z-axis will be computed based on the x and z ranges.
func (g *Grid) SetGridSize(gridSize int) {
	xRange := g.xMax - g.xMin
	zRange := g.zMax - g.zMin
	g.w = gridSize
	g.d = int(float64(g.w) * zRange / xRange)
}

func (g *Grid) SetOrigin(x, y float64) {
	g.originX = x
	g.originY = y
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

// SetWidth sets the width of the graph on the x axis.
func (g *Grid) SetWidth(w float64) {
	g.width = w
}

// SetXRange sets the min and max x values.
func (g *Grid) SetXRange(xMin, xMax float64) {
	g.xMin = xMin
	g.xMax = xMax
	xRange := g.xMax - g.xMin
	zRange := g.zMax - g.zMin
	g.d = int(float64(g.w) * zRange / xRange)
}

// SetYFunc sets the function that computes the y value for a given x and z.
func (g *Grid) SetYFunc(yFunc YFunction) {
	g.yFunc = yFunc
}

// SetYScale sets how much to scale the y axis.
func (g *Grid) SetYScale(yScale float64) {
	g.yScale = yScale
}

// SetZRange sets the min and max z values.
func (g *Grid) SetZRange(zMin, zMax float64) {
	g.zMin = zMin
	g.zMax = zMax
	xRange := g.xMax - g.xMin
	zRange := g.zMax - g.zMin
	g.d = int(float64(g.w) * zRange / xRange)
}
