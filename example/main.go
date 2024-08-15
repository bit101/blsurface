// Package main renders an image or video
package main

import (
	"math"

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
	renderTarget := target.Video
	fileName := "blsurface_smooth"

	if renderTarget == target.Image {
		render.CreateAndViewImage(700, 500, "out/"+fileName+".png", scene1, 0.0)
	} else if renderTarget == target.Video {
		program := render.NewProgram(600, 600, 30)
		program.AddSceneWithFrames(scene1, 360)
		program.RenderAndPlayVideo("out/frames", "out/"+fileName+".mp4")
	}
}

func scene1(context *cairo.Context, width, height, percent float64) {
	random.Seed(0)
	context.BlackOnWhite()
	context.SetLineWidth(0.33)
	context.Save()
	context.TranslateCenter()

	//////////////////////////////
	// make surface
	//////////////////////////////
	grid := blsurface.NewGrid(-1, -1, 1, 1, 60)
	grid.SetYFunc(stepped)
	grid.SetWidth(500)
	// grid.SetWidth(blmath.LoopSin(percent, 100, 600))

	// grid.SetTiltDegrees(360 * percent)
	// grid.SetTiltDegrees(210)
	grid.SetRotationDegrees(15)

	// grid.SetTiltDegrees(blmath.LoopSin(percent, -90, 90))
	grid.SetRotation(tau * percent)

	grid.DrawCells(context)
	context.Restore()
}

func concentricWave(x, z float64) float64 {
	return math.Sin(math.Hypot(x, z)*tau*2) * 0.25
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
	s := 1.0
	n := noise.Simplex2(x*s, z*s) * 0.5
	// n = blmath.RoundToNearest(n, 0.125)
	r := globe(x, z)

	return math.Min(n, r)
}

func globe(x, z float64) float64 {
	size := 0.5
	r := math.Hypot(x, z)
	if r < size {
		return -math.Sqrt(size*size-r*r) * 1
	}
	return 0
}
