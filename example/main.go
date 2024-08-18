// Package main renders an image or video
package main

import (
	"math"

	"github.com/bit101/bitlib/blcolor"
	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/noise"
	"github.com/bit101/bitlib/random"
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/blcairo/render"
	"github.com/bit101/blcairo/target"
	"github.com/bit101/blsurface"
)

//revive:disable:unused-parameter
const (
	tau = blmath.Tau
	pi  = math.Pi
)

func main() {
	random.RandSeed()
	renderTarget := target.Video
	fileName := "blsurface_smooth"

	if renderTarget == target.Image {
		render.CreateAndViewImage(600, 600, "out/"+fileName+".png", scene1, 0.0)
	} else if renderTarget == target.Video {
		program := render.NewProgram(600, 600, 30)
		program.AddSceneWithFrames(scene1, 120)
		program.RenderAndPlayVideo("out/frames", "out/"+fileName+".mp4")
	}
}

func scene1(context *cairo.Context, width, height, percent float64) {
	random.Seed(0)
	context.BlackOnWhite()
	// context.SetLineWidth(0.33)

	//////////////////////////////
	// make surface
	//////////////////////////////
	grid := blsurface.NewGrid()

	grid.SetOrigin(width/2, height/2, 100)
	grid.SetPerspective(true)
	// grid.SetFocalLength(100)

	grid.SetGridSize(40)

	grid.SetWidth(500)

	// grid.SetXRange(-2, 2)
	// grid.SetZRange(-2, 2)

	// grid.SetYScale(0.25)

	grid.SetYFunc(stepped)

	// grid.SetColorFunc(animColor(percent))

	// grid.SetRotation(0)
	// grid.SetRotationDegrees(140)
	grid.SetRotation(tau * percent)

	// grid.SetTilt(0)
	// grid.SetTiltDegrees(-40)
	// grid.SetTiltDegrees(360 * percent)
	// grid.SetTiltDegrees(blmath.LoopSin(percent, -45, 45))

	grid.DrawCells(context)

	// grid.DrawAxes(context, 300, 150, 100)
}

func concentricWave(x, z float64) float64 {
	return math.Sin(math.Hypot(x, z) * tau * 2)
}

func waves(x, z float64) float64 {
	return math.Sin(x*pi)*0.25 + math.Cos(z*pi*2)*0.2
}

func flat(x, z float64) float64 {
	return 0.0
}

func rando(x, z float64) float64 {
	return random.FloatRange(-0.07, 0.07)
}

func noisy(percent float64) blsurface.YFunction {
	s := 0.013
	return func(x, z float64) float64 {
		return noise.Simplex2(x*s, z*s) * blmath.LoopSin(percent, -30, 30)
	}
}

func staticNoise(x, z float64) float64 {
	s := 1.0
	n := noise.Simplex2(x*s, z*s) * 0.2
	return n
}

func stepped(x, z float64) float64 {
	s := 2.0
	n := noise.Simplex2(x*s, z*s) * 0.25
	// n = blmath.RoundToNearest(n, 0.125)
	r := globe(x, z)
	// r = blmath.RoundToNearest(r, 0.0625)
	if math.Hypot(x, z) < 0.1 {
		return 0.5
	}
	if math.Hypot(x, z) < 0.2 {
		return 0.35
	}

	return math.Max(n, r)
}

func globe(x, z float64) float64 {
	size := 0.5
	r := math.Hypot(x, z)
	if r < size {
		return math.Sqrt(size*size-r*r) * 1
	}
	return 0
}

//////////////////////////////
// Color funcs
//////////////////////////////

func concentricColor(x, y, z float64) blcolor.Color {
	dist := math.Hypot(x, z)
	return blcolor.HSV(dist*360, 0.5, 1)
	// return blcolor.HSV(p.Y*360+180, 0.5, 1)
}

func heightcolor(x, y, z float64) blcolor.Color {
	return blcolor.HSV(y*360+180, 0.5, 1)
}

func animColor(percent float64) blsurface.ColorFunc {
	return func(x, y, z float64) blcolor.Color {
		// return blcolor.HSV(y*360+percent*360, 0.5, 1)
		n := noise.Simplex3(x, y, z)
		n1 := noise.Simplex3(x*2+1, y*2+1, z*2+1) / 2
		return blcolor.HSV((n+n1)*180+percent*360, 0.5, 1)
	}
}
