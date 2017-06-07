package main

import (
	"flag"
	"fmt"
	"math"
	"time"
)

var (
	width        = flag.Int("w", 1024, "canvas width")
	height       = flag.Int("h", 768, "canvas height")
	focalLen     = flag.Float64("f", 1, "focal length")
	progressive  = flag.Int("p", 16, "progressively increase resolution")
	maxRecursion = flag.Int("r", 10, "maximum number of recursive rays")
)

const (
	Horiz = 20.0
	fine  = 0.02
	tol   = 1e-9
)

const deg = math.Pi / 180

var (
	Focal = Vec{0, 0, -1}
	scene *Scene
)

type Scene struct {
	light Vec
	amb   float64
	objs  []Obj
}

type Obj struct {
	Shape
	Shader ShaderFunc
}

func main() {
	Init()
	start := time.Now()

	const H = 2
	sp := coolSphere().RotX(-30*deg).Transl(0, 0, 6)
	sl := Slab(8, 0.1, 8).Transl(0, -H, 0).RotX(-30*deg).Transl(0, 0, 6)
	scene = &Scene{
		light: Vec{9, 3, -5},
		amb:   0.2,
		objs:  []Obj{{sp, ShadeDiffuse()}, {sl, ShadeReflect()}},
	}

	img := MakeImage(*width, *height)

	Render(scene, img)
	fmt.Println("done,", time.Since(start))
}

var nShade int

func Render(s *Scene, img [][]float64) {
	nShade = 0
	for sub := *progressive; sub > 0; sub /= 2 {
		refine(s, img, sub, sub == *progressive)
		Encode(img, "out.jpg")
	}
}

func refine(sc *Scene, img [][]float64, sub int, first bool) {
	W := *width
	H := *height
	for i := 0; i < H; i += sub {
		fmt.Printf("%.1f%%\n\u001B[F", float64(100*nShade)/float64((W+1)*(H+1)))
		for j := 0; j < W; j += sub {
			if i%(2*sub) == 0 && j%(2*sub) == 0 && !first {
				continue
			}
			nShade++
			y0 := (-float64(i) + float64(H)/2 + 0.5) / float64(H)
			x0 := (float64(j) - float64(W)/2 + 0.5) / float64(H)
			start := Vec{x0, y0, 0}
			r := Ray{start, start.Sub(Focal).Normalized()}

			v := PixelShade(sc, r, *maxRecursion)
			v = clip(v, 0, 1)

			for I := i; I < i+sub && I < H; I++ {
				for J := j; J < j+sub && J < W; J++ {
					img[I][J] = v
				}
			}
		}
	}
}

func PixelShade(sc *Scene, r Ray, N int) float64 {
	if N == 0 {
		return scene.amb
	}

	i, _ := Nearest(sc.objs, r)
	if i == -1 {
		return 0
	}
	obj := sc.objs[i]
	shape := obj.Shape

	pos, norm, ok := Normal(r, shape)
	if !ok {
		return 0
	}

	v := obj.Shader(pos, norm, r, N)

	//v = clip(v, 0, 1)
	return v
}

func Nearest(s []Obj, r Ray) (int, float64) {
	nearest := -1
	nearestZ := math.Inf(1)

	for i, s := range s {
		z, ok := Inters(r, s.Shape)
		if ok && z < nearestZ {
			nearestZ = z
			nearest = i
		}
	}
	return nearest, nearestZ
}

func inters(r Ray, s Shape) bool {
	_, ok := Inters(r, s)
	return ok
}

func intersAny(r Ray, s []Obj) bool {
	for _, s := range s {
		if inters(r, s.Shape) {
			return true
		}
	}
	return false
}

func Inters(r Ray, s Shape) (float64, bool) {
	for t := 0.0; t < Horiz; t += fine {
		if s(r.At(t)) {
			return t, true
		}
	}
	return 0, false
}

func Bisect(r Ray, s Shape) (Vec, bool) {
	in, ok := Inters(r, s)
	if !ok {
		return Vec{}, false
	}

	out := in - fine

	if s(r.At(out)) || !s(r.At(in)) {
		return Vec{}, false
	}

	for math.Abs(in-out)/(in+out) > tol {
		mid := (in + out) / 2
		if s(r.At(mid)) {
			in = mid
		} else {
			out = mid
		}
	}
	return r.At(in), true
}

func Normal(r Ray, s Shape) (Vec, Vec, bool) {
	c, ok := Bisect(r, s)
	if !ok {
		return Vec{}, Vec{}, false
	}

	ra := r
	ra.Dir = ra.Dir.Add(Vec{1e-5, 0, 0})
	a, okA := Bisect(ra, s)

	rb := r
	rb.Dir = rb.Dir.Add(Vec{0, 1e-5, 0})
	b, okB := Bisect(rb, s)

	if !okA || !okB {
		return Vec{}, Vec{}, false
	}

	a = a.Sub(c)
	b = b.Sub(c)

	n := b.Cross(a).Normalized()
	return c, n, true

}

func MakeImage(W, H int) [][]float64 {
	img := make([][]float64, H)
	for i := range img {
		img[i] = make([]float64, W)
	}
	return img
}

func clip(v, min, max float64) float64 {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	return v
}

func assert(t bool) {
	if !t {
		panic("assertion failed")
	}
}

func Init() {
	flag.Parse()
	Focal = Vec{0, 0, -*focalLen}
}
