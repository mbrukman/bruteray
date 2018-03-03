package br_test

import (
	"testing"

	. "github.com/barnex/bruteray/br"
	. "github.com/barnex/bruteray/mat"
	"github.com/barnex/bruteray/raster"
	. "github.com/barnex/bruteray/shape"
)

func BenchmarkSphere(b *testing.B) {
	e := NewEnv()
	e.Add(Sphere(Vec{0, 0, 1}, 0.25, Flat(WHITE)))
	c := raster.Camera(0)

	benchmark(b, e, c)
}

func Benchmark9Spheres(b *testing.B) {
	e := NewEnv()
	r := 0.5

	//nz := ShadeNormal(Ez)
	e.Add(Sphere(Vec{0, 0, 0}, r, Flat(WHITE)))
	e.Add(Sphere(Vec{0, 0, 2}, r, Flat(WHITE)))
	e.Add(Sphere(Vec{0, 0, 4}, r, Flat(WHITE)))

	e.Add(Sphere(Vec{2, 0, 0}, r, Flat(WHITE)))
	e.Add(Sphere(Vec{2, 0, 2}, r, Flat(WHITE)))
	e.Add(Sphere(Vec{2, 0, 4}, r, Flat(WHITE)))

	e.Add(Sphere(Vec{-2, 0, 0}, r, Flat(WHITE)))
	e.Add(Sphere(Vec{-2, 0, 2}, r, Flat(WHITE)))
	e.Add(Sphere(Vec{-2, 0, 4}, r, Flat(WHITE)))

	c := raster.Camera(1).Transl(0, 4, -4).Transf(RotX4(Pi / 5))
	benchmark(b, e, c)
}

func benchmark(b *testing.B, e *Env, c *raster.Cam) {
	b.SetBytes((testW + 1) * (testH + 1))
	img := raster.MakeImage(testW, testH)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raster.SinglePass(c, e, img)
	}
}