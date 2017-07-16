package bruteray

import (
	"math"
)

type Shape interface {
	// Ray-shape intersection.
	// t values may be < 0 (behind camera), but must be sorted (min <= max)
	Inters(r *Ray) Interval

	// Normal vector at position.
	// Does not necessarily need to point outwards.
	Normal(pos Vec) Vec
}

// -- sphere

func Sphere(center Vec, radius float64) *sphere {
	return &sphere{center, Sqr(radius)}
}

type sphere struct {
	c  Vec
	r2 float64
}

func (s *sphere) Normal(pos Vec) Vec {
	n := pos.Sub(s.c).Normalized()
	return n
}

func (s *sphere) Inters(r *Ray) Interval {
	v := r.Start.Sub(s.c)
	d := r.Dir
	vd := v.Dot(d)
	D := Sqr(vd) - (v.Len2() - s.r2)
	if D < 0 {
		return Interval{}
	}
	t1 := (-vd - math.Sqrt(D))
	t2 := (-vd + math.Sqrt(D))
	return Interv(t1, t2)
}

// -- sheet (infinite)

func Sheet(dir Vec, off float64) *sheet {
	return &sheet{dir, off}
}

type sheet struct {
	dir Vec
	off float64
}

func (s *sheet) Normal(pos Vec) Vec {
	return s.dir
}

func (s *sheet) Inters(r *Ray) Interval {
	rs := r.Start.Dot(s.dir)
	rd := r.Dir.Dot(s.dir)
	t := (s.off - rs) / rd
	return Interval{t, t}
}

// --rectangle (finite sheet)

// A rectangle (i.e. finite sheet) at given position,
// with normal vector dir and half-axes rx, ry, rz.
func Rect(pos, dir Vec, rx, ry, rz float64) Shape {
	return &rect{pos, dir, rx, ry, rz}
}

type rect struct {
	pos, dir   Vec
	rx, ry, rz float64
}

func (s *rect) Inters(r *Ray) Interval {
	rs := r.Start.Dot(s.dir)
	rd := r.Dir.Dot(s.dir)
	t := (s.pos.Dot(s.dir) - rs) / rd
	p := r.At(t).Sub(s.pos)
	if p[X] < -s.rx || p[X] > s.rx ||
		p[Y] < -s.ry || p[Y] > s.ry ||
		p[Z] < -s.rz || p[Z] > s.rz {
		return Interval{}
	}
	return Interval{t, t}
}

func (s *rect) Normal(p Vec) Vec {
	return s.dir
}

// -- slab

func Slab(dir Vec, off1, off2 float64) *slab {
	return &slab{dir, off1, off2}
}

type slab struct {
	dir        Vec
	off1, off2 float64
}

func (s *slab) Normal(pos Vec) Vec {
	return s.dir
}

func (s *slab) Inters(r *Ray) Interval {
	rs := r.Start.Dot(s.dir)
	rd := r.Dir.Dot(s.dir)
	t1 := (s.off1 - rs) / rd
	t2 := (s.off2 - rs) / rd
	t1, t2 = Sort(t1, t2)
	return Interval{t1, t2}
}
