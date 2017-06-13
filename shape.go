package main

import "math"

type Shape interface {
	//Inters(r Ray) (float64, float64)
	Normal(r Ray) (float64, Vec, bool)
}

//type ShapeFunc func(r Ray) (t float64, normal Vec, intersect bool)
//
//func (f ShapeFunc) Normal(r Ray) (t float64, normal Vec, intersect bool) {
//	return f(r)
//}

//func And(a, b Shape) Shape {
//	return ShapeFunc(func(r Ray) (float64, Vec, bool) {
//		ta, na, oka := a.Normal(r)
//		if !oka || ta < 0 {
//			return 0, Vec{}, false
//		}
//		tb, nb, okb := b.Normal(r)
//		if !okb || tb < 0 {
//			return 0, Vec{}, false
//		}
//
//		if ta < tb {
//			return ta, na, oka
//		} else {
//			return tb, nb, okb
//		}
//	})
//}

func Sphere(c Vec, r float64) Shape {
	return sphere{c, r}
}

type sphere struct {
	c Vec
	r float64
}

func (s sphere) Normal(ray Ray) (float64, Vec, bool) {
	v := ray.Start.Sub(s.c)
	d := ray.Dir
	D := sqr(v.Dot(d)) - (v.Dot(v) - s.r*s.r)
	if D < 0 {
		return 0, Vec{}, false
	}
	t := (-v.Dot(d) - math.Sqrt(D)) / d.Len2()
	if t < 0 {
		return 0, Vec{}, false
	}
	n := ray.At(t).Sub(s.c).Normalized()

	return t, n, true
}

type sheetY struct {
	y float64
}

func SheetY(y float64) Shape {
	return sheetY{y}
}

func (s sheetY) Normal(r Ray) (float64, Vec, bool) {
	t := (s.y - r.Start.Y) / r.Dir.Y
	if t < 0 {
		return 0, Vec{}, false
	}
	n := Vec{0, 1, 0}
	return t, n, true
}

//func ABox(min, max Vec) Shape {
//	return ShapeFunc(func(r Ray) (float64, Vec, bool) {
//		tmin := min.Sub(r.Start).Div3(r.Dir)
//		tmax := max.Sub(r.Start).Div3(r.Dir)
//
//		txen := Min(tmin.X, tmax.X)
//		txex := Max(tmin.X, tmax.X)
//
//		tyen := Min(tmin.Y, tmax.Y)
//		tyex := Max(tmin.Y, tmax.Y)
//
//		tzen := Min(tmin.Z, tmax.Z)
//		tzex := Max(tmin.Z, tmax.Z)
//
//		ten := Max3(txen, tyen, tzen)
//		tex := Min3(txex, tyex, tzex)
//
//		return ten, Vec{}, ten < tex
//	})
//}

func Min(x, y float64) float64 {
	return math.Min(x, y)
}

func Min3(x, y, z float64) float64 {
	return Min(Min(x, y), z)
}

func Max(x, y float64) float64 {
	return math.Max(x, y)
}

func Max3(x, y, z float64) float64 {
	return Max(Max(x, y), z)
}
