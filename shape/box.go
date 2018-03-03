package shape

import (
	"fmt"
	. "github.com/barnex/bruteray/br"
	"math"
)

// NBox constructs a box with given width, depth and height.
func NBox(w, h, d float64, m Material) *box {
	rx, ry, rz := w/2, h/2, d/2
	return &box{
		Min: Vec{rx, ry, rz}.Mul(-1),
		Max: Vec{rx, ry, rz},
		Mat: m,
	}
}

type box struct {
	Min, Max Vec
	Mat      Material
}

func (s *box) Center() Vec {
	return (s.Min.Add(s.Max)).Mul(0.5)
}

func (s *box) Transl(d Vec) *box {
	s.Min.Transl(d)
	s.Max.Transl(d)
	return s
}

// Corner returns one of the box's corners:
// 	Corner( 1, 1, 1) -> right top  back
// 	Corner(-1,-1,-1) -> left bottom front
// 	Corner( 1,-1,-1) -> right bottom front
// 	...
func (s *box) Corner(x, y, z int) Vec {
	which := Vec{float64(x), float64(y), float64(z)}
	return s.Center().Add(s.Max.Sub(s.Center()).Mul3(which))
}

func Box(center Vec, rx, ry, rz float64, m Material) CSGObj {
	return &box{
		Min: center.Sub(Vec{rx, ry, rz}),
		Max: center.Add(Vec{rx, ry, rz}),
		Mat: m,
	}
}

func Cube(center Vec, r float64, m Material) CSGObj {
	return Box(center, r, r, r, m)
}

func (s *box) Hit1(r *Ray, f *[]Fragment) { s.HitAll(r, f) }

func (s *box) HitAll(r *Ray, f *[]Fragment) {
	min_ := s.Min
	max_ := s.Max

	tmin := min_.Sub(r.Start).Mul3(r.InvDir)
	tmax := max_.Sub(r.Start).Mul3(r.InvDir)

	txen := min(tmin[X], tmax[X])
	txex := max(tmin[X], tmax[X])

	tyen := min(tmin[Y], tmax[Y])
	tyex := max(tmin[Y], tmax[Y])

	tzen := min(tmin[Z], tmax[Z])
	tzex := max(tmin[Z], tmax[Z])

	ten := max3(txen, tyen, tzen)
	tex := min3(txex, tyex, tzex)

	if ten > tex {
		return
	}

	if math.IsNaN(ten) || math.IsNaN(tex) {
		return
	}

	*f = append(*f,
		Fragment{T: ten, Norm: s.Normal(r.At(ten)), Material: s.Mat},
		Fragment{T: tex, Norm: s.Normal(r.At(tex)), Material: s.Mat},
	)

}

func (s *box) Inside(v Vec) bool {
	return v[X] > s.Min[X] && v[X] < s.Max[X] &&
		v[Y] > s.Min[Y] && v[Y] < s.Max[Y] &&
		v[Z] > s.Min[Z] && v[Z] < s.Max[Z]
}

func (s *box) Normal(p Vec) Vec {
	//p.check()
	for i := range p {
		if approx(p[i], s.Min[i]) || approx(p[i], s.Max[i]) {
			return Unit[i]
		}
	}

	panic(fmt.Sprint("box.normal", p, s.Min, s.Max))
}

func approx(a, b float64) bool {
	return math.Abs(a-b) < 1e-4
}