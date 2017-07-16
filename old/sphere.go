package main

import "math"

type sphere struct {
	c  Vec
	r2 float64
}

func Sphere(center Vec, radius float64) *sphere {
	return &sphere{center, sqr(radius)}
}

func (s *sphere) Normal(r *Ray, t float64) Vec {
	return r.At(t).Sub(s.c).Normalized()
}

func (s *sphere) Hit(r *Ray) float64 {
	v := r.Start.Sub(s.c)
	d := r.Dir
	vd := v.Dot(d)
	D := sqr(vd) - (v.Len2() - s.r2)
	if D < 0 {
		return 0
	}
	if t1 := (-vd - math.Sqrt(D)); t1 > 0 {
		return t1
	}
	if t2 := (-vd + math.Sqrt(D)); t2 > 0 {
		return t2
	}

	return 0
}

func (s *sphere) Inters(r *Ray) Interval {
	v := r.Start.Sub(s.c)
	d := r.Dir
	vd := v.Dot(d)
	D := sqr(vd) - (v.Len2() - s.r2)
	if D < 0 {
		return empty
	}
	t1 := (-vd - math.Sqrt(D))
	t2 := (-vd + math.Sqrt(D))
	t1, t2 = Sort(t1, t2)
	return Interval{t1, t2}
}

func (s *sphere) Transl(dx, dy, dz float64) *sphere {
	return &sphere{s.c.Add(Vec{dx, dy, dz}), s.r2}
}