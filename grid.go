// Package blsurface implements a 3d grid surface.
package blsurface

import (
	"math"
	"slices"

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
	cells                     []*GridPoint
	faces                     []*Face
	w, d                      int
	originX, originY, originZ float64
	rotation                  float64
	tilt                      float64
	yFunc                     YFunction
	colorFunc                 ColorFunc
	xMin, zMin, xMax, zMax    float64
	yScale                    float64
	width                     float64
	axes                      []*GridPoint
	perspective               bool
	fl                        float64
}

//////////////////////////////
// Constructor
//////////////////////////////

// NewGrid creates a new grid.
func NewGrid() *Grid {

	return &Grid{
		w:           20,
		d:           20,
		originX:     0,
		originY:     0,
		originZ:     200,
		rotation:    math.Pi / 6,
		tilt:        math.Pi / 6,
		yFunc:       func(x, z float64) float64 { return 0.0 },
		colorFunc:   func(x, y, z float64) blcolor.Color { return blcolor.White },
		yScale:      1.0,
		xMin:        -1,
		zMin:        -1,
		xMax:        1,
		zMax:        1,
		width:       400,
		perspective: false,
		fl:          350.0,
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

// DrawAxes draws a representation of the x, y, z axis lines.
func (g *Grid) DrawAxes(context *cairo.Context, x, y, size float64) {
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

// DrawCells draws the whole grid.
func (g *Grid) DrawCells(context *cairo.Context) {
	g.makeGrid()
	g.applyFunc()
	g.transform()
	context.Save()
	context.Translate(g.originX, g.originY)
	scale := g.width / (g.xMax - g.xMin)
	slices.SortFunc(g.faces, sortFaces)
	for _, face := range g.faces {
		face.Draw(context, g.perspective, g.originZ, g.fl, scale, g.colorFunc)
	}
	context.Restore()
}

func sortFaces(a, b *Face) int {
	if a.Zpos() > b.Zpos() {
		return -1
	}
	return 1
}

// DrawPoints draws each point in the grid.
func (g *Grid) DrawPoints(context *cairo.Context, radius float64) {
	for _, p := range g.cells {
		context.FillCircle(p.X, p.Y, radius)
	}
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
	for i := 0; i < g.w; i++ {
		for j := 0; j < g.d; j++ {
			x := i
			z := j
			p0 := g.getCell(x, z)
			p1 := g.getCell(x+1, z)
			p2 := g.getCell(x+1, z+1)
			p3 := g.getCell(x, z+1)

			g.faces = append(g.faces, NewFace(p0, p1, p2, p3))
		}
	}
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

// SetFocalLength sets the value used to create perspective.
func (g *Grid) SetFocalLength(fl float64) {
	g.fl = fl
}

// SetGridSize sets how many cells to draw across the x-axis.
// The number of cells for the z-axis will be computed based on the x and z ranges.
func (g *Grid) SetGridSize(gridSize int) {
	xRange := g.xMax - g.xMin
	zRange := g.zMax - g.zMin
	g.w = gridSize
	g.d = int(float64(g.w) * zRange / xRange)
}

// SetOrigin sets the center x, y, z point from which the surface will be drawn.
func (g *Grid) SetOrigin(x, y, z float64) {
	g.originX = x
	g.originY = y
	g.originZ = z
}

// SetPerspective sets whether the surface will be drawn with perspective (true) or orthographically (false).
func (g *Grid) SetPerspective(b bool) {
	g.perspective = b
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
