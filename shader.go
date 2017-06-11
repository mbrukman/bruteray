package main

type Shader func(r Ray, t float64, normal Vec, N int) float64

func (s Shader) Intensity(r Ray, t float64, n Vec, N int) float64 {
	return s(r, t, n, N)
}

func Flat(v float64) Shader {
	return func(r Ray, t float64, normal Vec, N int) float64 {
		return v
	}
}

// Diffuse shading with shadows, but no interreflection
func Diffuse1(reflect float64) Shader {
	return func(r Ray, t float64, n Vec, N int) float64 {
		return diffuse1(reflect, r, t, n, N)
	}
}

func diffuse1(reflect float64, r Ray, t float64, n Vec, N int) float64 {
	p := r.At(t).MAdd(off, n)
	acc := 0.
	for _, light := range sources {
		lightPos, flux := light.Sample()
		d := lightPos.Sub(p)
		if !intersectsAny(Ray{p, d.Normalized()}) {
			acc += reflect * flux * n.Dot(d) / (d.Len2())
		}
	}
	return acc
}

// Diffuse shading with shadows and interreflection
func Diffuse2(reflect float64) Shader {
	return func(r Ray, t float64, n Vec, N int) float64 {
		acc := diffuse1(reflect, r, t, n, N)
		p := r.At(t).MAdd(off, n)
		d := randVec(n)
		sec := Ray{p, d}
		acc += reflect * Intensity(sec, N+1, false) * n.Dot(d.Normalized())
		return acc
	}
}

func intersectsAny(r Ray) bool {
	_, _, obj := FirstIntersect(r, false)
	return obj != nil
}

func Reflective(reflect float64) Shader {
	return func(r Ray, t float64, n Vec, N int) float64 {
		p := r.At(t).MAdd(off, n)
		dir2 := reflectVec(r.Dir, n)
		return reflect * Intensity(Ray{p, dir2}, N+1, true)
	}
}

func ReflectiveMate(reflect float64, jitter float64) Shader {
	return func(r Ray, t float64, n Vec, N int) float64 {
		p := r.At(t).MAdd(off, n)
		dir2 := reflectVec(r.Dir, n).MAdd(jitter, randVec(n))
		return reflect * Intensity(Ray{p, dir2}, N+1, true)
	}
}

func ShaderAdd(a, b Shader) Shader {
	return func(r Ray, t float64, n Vec, N int) float64 {
		return a.Intensity(r, t, n, N) + b.Intensity(r, t, n, N)
	}
}

// reflects v of the surface with normal n.
func reflectVec(v, n Vec) Vec {
	return v.MAdd(-2*v.Dot(n), n)
}
