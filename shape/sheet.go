package shape

import . "github.com/barnex/bruteray/br"

func Sheet(dir Vec, off float64, m Material) Obj {
	return &sheet{dir: dir, off: off, m: m}
}

type sheet struct {
	dir Vec
	off float64
	m   Material
}

func (s *sheet) Hit1(r *Ray, f *[]Fragment) {
	rs := r.Start.Dot(s.dir)
	rd := r.Dir().Dot(s.dir)
	t := (s.off - rs) / rd

	*f = append(*f, Fragment{T: t, Norm: s.dir, Material: s.m})
}