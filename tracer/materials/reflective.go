package materials

import (
	"math"

	"github.com/barnex/bruteray/geom"
	. "github.com/barnex/bruteray/imagef/colorf"
	. "github.com/barnex/bruteray/tracer"
	. "github.com/barnex/bruteray/util"
)

// A Reflective surface. E.g.:
// 	Reflective(WHITE)        // perfectly reflective, looks like shiny metal
// 	Reflective(WHITE.EV(-1)) // 50% reflective, looks like darker metal
// 	Reflective(RED)          // Reflects only red, looks like metal in transparent red candy-wrap.
func Reflective(c Color) Material {
	return &reflective{c}
}

type reflective struct {
	c Color
}

func (m *reflective) Shade(ctx *Ctx, s *Scene, r *Ray, h HitCoords) Color {
	pos := r.At(h.T - Tiny)
	r2 := ctx.Ray()
	defer ctx.PutRay(r2)
	r2.Start = pos
	r2.Dir = reflect(r.Dir, h.Normal)
	return s.LightField(ctx, r2).Mul3(m.c)
}

// ReflectFresnel is a transparent material with index of refraction n,
// on top of material transmitted. E.g. a wet or varnished material.
// This looks similar to simple reflection,
// but reflection is stronger under grazing incidence.
// E.g.:
//  ReflectFresnel(1.33, BLACK)           // a thin film of water on a black surface
//  ReflectFresnel(1.33, Diffuse(WHITE))  // milk
//  ReflectFresnel(20, BLACK)             // metal
func ReflectFresnel(n float64, transmitted Material) Material {
	return &reflectFresnel{n, transmitted}
}

type reflectFresnel struct {
	n     float64
	trans Material
}

func (s *reflectFresnel) Shade(ctx *Ctx, e *Scene, r *Ray, h HitCoords) Color {
	pos, norm := r.At(h.T-Tiny), h.Normal
	r2 := ctx.Ray()
	r2.Start = pos
	r2.Dir = reflect(r.Dir, norm)
	R := fresnelReflection(1, s.n, math.Abs(norm.Dot(r.Dir)))
	T := 1 - R
	trans := s.trans.Shade(ctx, e, r, h)
	refl := e.LightField(ctx, r2)
	ctx.PutRay(r2)
	return refl.Mul(R).MAdd(T, trans)
}

func fresnelReflection(n1, n2, cosθi float64) float64 {
	n12 := n1 / n2
	sin2θt := n12 * n12 * (1 - cosθi*cosθi) // sin² of transsion angle, using Snell's law.
	cosθt := math.Sqrt(1 - sin2θt)
	Rp := Sqr((n1*cosθi - n2*cosθt) / (n1*cosθi + n2*cosθt))
	Rs := Sqr((n1*cosθt - n2*cosθi) / (n1*cosθt + n2*cosθi))
	return 0.5 * (Rp + Rs)
}

// reflect reflects v along normal vector n.
//      n
//  v   |   reflected
//   \  |  /
//    \θ|θ/
//     \|/
//  --------- surface.
//
// See https://en.wikipedia.org/wiki/Ray_tracing_(graphics)#Example
func reflect(v, n geom.Vec) geom.Vec {
	return v.MAdd(-2*v.Dot(n), n)
}
