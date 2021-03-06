package objects

import (
	"github.com/barnex/bruteray/tracer"
	. "github.com/barnex/bruteray/tracer/types"
	"github.com/barnex/bruteray/util"
)

// TODO: change API: Or, And accept multiple args,
// build balanced tree

// And returns the intersection (boolean AND) of two objects.
func And(a, b Interface) Interface {
	return &and{a, b}
}

type and struct {
	a, b Interface
}

func (o *and) Intersect(r *Ray) HitRecord {
	a := andMarch(o.a, o.b, r, inf)
	max := inf
	if a.T > 0 {
		max = a.T // potential early return: solution a may occlude solution b
	}
	b := andMarch(o.b, o.a, r, max)
	return tracer.Frontmost(&a, &b)
}

// March along the ray until we find an intersection with a that is inside b.
// Do not march further than maxT. (maxT is the position of an earlier solution.
// It would occlude, so there is no point in searching beyond it.)
func andMarch(a, b Interface, r *Ray, maxT float64) HitRecord {
	backup := r.Start

	tOff := 0.0
	h := a.Intersect(r)
	ttl := 100 // hack: time-to-life neede in case a ray hits parallel to a surface and takes forever to escape
	for h.T > 0 && ttl != 0 && h.T < maxT {
		ttl--
		if b.Inside(r.At(h.T)) {
			h.T += tOff
			setStart(r, backup)
			return h
		}
		deltaT := h.T + Tiny
		tOff += deltaT
		r.Start = r.At(deltaT)
		h = a.Intersect(r)
	}
	setStart(r, backup)
	return HitRecord{}
}

func setStart(r *Ray, start Vec) {
	r.Start = start
}

func (o *and) Bounds() BoundingBox {
	return intersectBounds(o.a.Bounds(), o.b.Bounds())
}

func intersectBounds(a, b BoundingBox) BoundingBox {
	return BoundingBox{
		Min: Vec{
			util.Max(a.Min[0], b.Min[0]),
			util.Max(a.Min[1], b.Min[1]),
			util.Max(a.Min[2], b.Min[2]),
		},
		Max: Vec{
			util.Min(a.Max[0], b.Max[0]),
			util.Min(a.Max[1], b.Max[1]),
			util.Min(a.Max[2], b.Max[2]),
		},
	}
}

func (o *and) Inside(p Vec) bool {
	return o.a.Inside(p) && o.b.Inside(p)
}

func Or(a, b Interface) Interface {
	return &or{a, b}
}

type or struct {
	a, b Interface
}

func (o *or) Intersect(r *Ray) HitRecord {
	a := orMarch(o.a, o.b, r, inf)
	max := inf
	if a.T > 0 {
		max = a.T // potentinal early return: solution a may occlude solution b.
	}
	b := orMarch(o.b, o.a, r, max)
	return tracer.Frontmost(&a, &b)
}

func (o *or) Inside(p Vec) bool {
	return o.a.Inside(p) || o.b.Inside(p)
}

func (o *or) Bounds() BoundingBox {
	a := o.a.Bounds()
	b := o.b.Bounds()
	return BoundingBox{
		Min: Vec{
			util.Min(a.Min[0], b.Min[0]),
			util.Min(a.Min[1], b.Min[1]),
			util.Min(a.Min[2], b.Min[2]),
		},
		Max: Vec{
			util.Max(a.Max[0], b.Max[0]),
			util.Max(a.Max[1], b.Max[1]),
			util.Max(a.Max[2], b.Max[2]),
		},
	}
}

// March along the ray until we find an intersection with a that is not inside b.
// Do not march further than maxT. (maxT is the position of an earlier solution.
// It would occlude, so there is no point in searching beyond it.)
func orMarch(a, b Interface, r *Ray, maxT float64) HitRecord {
	backup := r.Start

	tOff := 0.0
	h := a.Intersect(r)
	ttl := 100 // hack: time-to-life neede in case a ray hits parallel to a surface and takes forever to escape
	for h.T > 0 && ttl != 0 && h.T < maxT {
		ttl--
		if !b.Inside(r.At(h.T)) {
			h.T += tOff
			setStart(r, backup)
			return h
		}
		deltaT := h.T + Tiny
		tOff += deltaT
		r.Start = r.At(deltaT)
		h = a.Intersect(r)
	}
	setStart(r, backup)
	return HitRecord{}
}

func Restrict(a, b Interface) Interface {
	return &restrict{a, b}
}

type restrict struct {
	orig   Interface
	inside Interface
}

func (o *restrict) Intersect(r *Ray) HitRecord {
	return andMarch(o.orig, o.inside, r, inf)
}

func (o *restrict) Bounds() BoundingBox {
	return intersectBounds(o.orig.Bounds(), o.inside.Bounds())
}

func (o *restrict) Inside(p Vec) bool {
	return o.orig.Inside(p) && o.inside.Inside(p)
}

type not struct {
	orig Interface
}

// Not returns the "inverse" of an object:
// The insideness and normal vectors are reversed with respect to the original.
//
// This is useful for boolean operations, E.g.:
// 	object1.And(Not(object2))
// removes from object1 all points that are inside object2.
func Not(object Interface) Interface {
	return &not{object}
}

func (b *not) Intersect(r *Ray) HitRecord {
	h := b.orig.Intersect(r)
	h.Normal = h.Normal.Mul(-1) // flip normal so it points outwards again.
	return h
}

func (b *not) Bounds() BoundingBox {
	return infBox
}

func (b *not) Inside(p Vec) bool {
	return !b.orig.Inside(p)
}

func Difference(a, b Interface) Interface {
	return And(a, Not(b))
}

// Hollow returns the object with a modified Inside method that always returns false.
// This causes the object to become hollow inside (only its surface remains).
//
// This has a visible effect only when parts of an object are cut away
// (e.g. with And(Not(...))), so that the inside is revealed.
func Hollow(orig Interface) Interface {
	return &hollow{orig: orig}
}

type hollow struct {
	hollowSurface
	orig Interface
}

func (o *hollow) Intersect(r *Ray) HitRecord {
	return o.orig.Intersect(r)
}

func (o *hollow) Bounds() BoundingBox {
	return o.orig.Bounds()
}
