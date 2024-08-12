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
	fileName := "blsurface"

	if renderTarget == target.Image {
		render.CreateAndViewImage(700, 500, "out/"+fileName+".png", scene1, 0.0)
	} else if renderTarget == target.Video {
		program := render.NewProgram(560, 460, 30)
		program.AddSceneWithFrames(scene1, 360)
		program.RenderAndPlayVideo("out/frames", "out/"+fileName+".mp4")
	}
}

func scene1(context *cairo.Context, width, height, percent float64) {
	random.Seed(0)
	context.BlackOnWhite()
	context.SetLineWidth(0.25)
	context.Save()
	context.TranslateCenter()

	grid := blsurface.NewGrid(100, 100, 10, globe)
	grid.Rotate(-percent * tau)
	grid.Tilt(pi * 0.2)
	// grid.Tilt(blmath.LoopSin(percent, -pi*0.15, pi*0.35))
	// grid.DrawPoints(context)
	grid.DrawCells(context)
	context.Restore()
}

func concentricWave(x, z float64) float64 {
	return math.Sin(math.Hypot(x, z)*0.05) * 40
}

func waves(x, z float64) float64 {
	return math.Sin(x*0.07)*10 + math.Cos(z*0.05)*20
}

func flat(x, z float64) float64 {
	return 0.0
}

func rando(x, z float64) float64 {
	return random.FloatRange(-20, 20)
}

func noisy(percent float64) blsurface.YFunction {
	s := 0.013
	return func(x, z float64) float64 {
		return noise.Simplex2(x*s, z*s) * blmath.LoopSin(percent, -30, 30)
	}
}

func stepped(x, z float64) float64 {
	s := 0.005
	n := noise.Simplex2(x*s, z*s)*40 + 40
	n = blmath.RoundToNearest(n, 70)
	return n
}

func globe(x, z float64) float64 {
	size := 181.0
	r := math.Hypot(x, z)
	if r < size {
		s := -math.Sqrt(size*size - r*r)
		size *= 0.75
		if r < size {
			s += math.Sqrt(size*size - r*r)
		}
		size *= 0.25
		if r < size {
			return -180
		}
		return s
	}
	return 0
}
