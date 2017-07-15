package r

import (
	"image"
	"image/png"
	"os"
	"testing"
)

const (
	testW, testH = 300, 200 // test image size
	testRec      = 3        // test recursion depth
)

// Test a flat sphere
func TestSphere(tst *testing.T) {
	t := Helper(tst)

	e := NewEnv()
	e.Add(Object(Sphere(Vec{0, 0, 1}, 0.25), Flat(WHITE)))
	c := Camera(0)

	t.Compare(e, c, "001-sphere")
}

// Test a sphere behind the camera
func TestBehindCam(tst *testing.T) {
	t := Helper(tst)

	e := NewEnv()
	e.Add(Object(Sphere(Vec{0, 0, -1}, 0.25), Flat(WHITE)))

	t.Compare(e, Camera(0), "002-behindcam")
}

// Test normal vectors
func TestNormal(tst *testing.T) {
	t := Helper(tst)

	e := NewEnv()
	e.Add(Object(Sphere(Vec{0, 0, 2}, 0.25), ShadeNormal(Ez)))
	e.Add(Object(Sphere(Vec{-0.5, 0, 2}, 0.25), ShadeNormal(Ex)))
	e.Add(Object(Sphere(Vec{0.5, 0, 2}, 0.25), ShadeNormal(Ey)))

	t.Compare(e, Camera(0), "003-normals")
}

// Test camera translation
func TestCamTransl(tst *testing.T) {
	t := Helper(tst)

	e := NewEnv()
	e.Add(Object(Sphere(Vec{0, 0, 2}, 0.25), ShadeNormal(Ez)))

	t.Compare(e, Camera(0).Transl(-0.5, -0.25, 0), "004-camtransl")
}

// Test camera rotation
func TestCamRot(tst *testing.T) {
	t := Helper(tst)

	e := NewEnv()
	r := 0.5
	nz := ShadeNormal(Ez)
	e.Add(Object(Sphere(Vec{0, 0, 0}, r), nz))
	e.Add(Object(Sphere(Vec{0, 0, 2}, r), nz))
	e.Add(Object(Sphere(Vec{0, 0, 4}, r), nz))

	e.Add(Object(Sphere(Vec{2, 0, 0}, r), nz))
	e.Add(Object(Sphere(Vec{2, 0, 2}, r), nz))
	e.Add(Object(Sphere(Vec{2, 0, 4}, r), nz))

	e.Add(Object(Sphere(Vec{-2, 0, 0}, r), nz))
	e.Add(Object(Sphere(Vec{-2, 0, 2}, r), nz))
	e.Add(Object(Sphere(Vec{-2, 0, 4}, r), nz))

	t.Compare(e, Camera(1).Transl(0, 4, -4).Transf(RotX4(pi/5)), "005-camrot")
}

// Test object transform
func TestObjTransf(tst *testing.T) {
	t := Helper(tst)

	e := NewEnv()
	r := 0.25
	sx := Object(Sphere(Vec{-0.5, 0, 2}, r), ShadeNormal(Ex))
	sy := Object(Sphere(Vec{0, 0, 2}, r), ShadeNormal(Ez))
	sz := Object(Sphere(Vec{0.5, 0, 2}, r), ShadeNormal(Ey))

	rot := RotZ4(pi / 4)
	e.Add(Transf(sx, rot))
	e.Add(Transf(sy, rot))
	e.Add(Transf(sz, rot))

	t.Compare(e, Camera(0), "006-objtransf")
}

// Test intersection of two spheres
func TestObjAnd(tst *testing.T) {
	t := Helper(tst)

	e := NewEnv()
	r := 0.5
	s1 := Object(Sphere(Vec{-r / 2, 0, 2}, r), ShadeNormal(Ez))
	s2 := Object(Sphere(Vec{r / 2, 0, 2}, r), ShadeNormal(Ey))
	s := ObjAnd(s1, s2)
	e.Add(s)

	t.Compare(e, Camera(0), "007-objand")
}

// Test two partially overlapping spheres
func TestOverlap(tst *testing.T) {
	t := Helper(tst)

	e := NewEnv()
	r := 0.5
	s1 := Object(Sphere(Vec{-r / 2, 0, 2}, r), ShadeNormal(Ez))
	s2 := Object(Sphere(Vec{r / 2, 0, 2}, r), ShadeNormal(Ey))
	e.Add(s1)
	e.Add(s2)

	t.Compare(e, Camera(0), "008-overlap")
}

// Make a cube out of 3 intersecting slabs
func TestSlabIntersect(tst *testing.T) {
	t := Helper(tst)

	e := NewEnv()
	r := 1.
	s1 := Object(Slab(Ex, -r, r), Flat(RED))
	s2 := Object(Slab(Ey, -r, r), Flat(GREEN))
	s3 := Object(Slab(Ez, -r, r), Flat(BLUE))
	cube := ObjAnd(ObjAnd(s1, s2), s3)
	cube = Transf(cube, RotY4(160*deg).Mul(RotX4(20*deg)))
	e.Add(cube)

	t.Compare(e, Camera(1).Transl(0, 0, -4), "009-slabintersect")
}

//func TestSpheres(tst *testing.T) {
//	t := Helper(tst)
//
//	scene := NewEnv()
//	scene.amb = func(v Vec) Color { return Color(0.2*v.Y + 0.2) }
//	scene.Add(Sheet(-3, Ey), Diffuse1(0.5))  // floor
//	scene.Add(Sheet(8, Ey), Diffuse1(0.5))   // ceiling
//	scene.Add(Sheet(20, Ez), Diffuse1(0.8))  // back
//	scene.Add(Sheet(-5, Ez), Diffuse1(0.8))  // front
//	scene.Add(Sheet(10, Ex), Diffuse1(0.8))  // left
//	scene.Add(Sheet(-10, Ex), Diffuse1(0.8)) // right
//	scene.Add(Sphere(Vec{1, -2, 8}, 1), Reflective(0.3))
//	scene.Add(Sphere(Vec{-1, -2, 6}, 1), Diffuse1(0.95))
//	scene.AddLight(PointLight(Vec{0, 7, 1}, 100))
//
//	cam := Camera(1)
//	cam.Transf(RotX(-5 * deg))
//	t.CompareCam(scene, "009-spheres", cam)
//}
//
//func TestCheckers(tst *testing.T) {
//	t := Helper(tst)
//
//	s := NewEnv()
//	s.amb = func(dir Vec) Color { return 0.5 }
//
//	s.Add(Sheet(0, Ey), Diffuse1(0.7))                                                      // floor
//	s.Add(Box(Vec{0, 0, 8}, 5, 0.45, 5), Diffuse1(0.1))                                     // base
//	s.Add(Rect(Vec{0, 0.5, 8}, Ey, 4, inf, 4), CheckBoard(Reflective(0.05), Diffuse1(0.9))) // checkboard
//	s.Add(Sheet(20, Ez), Diffuse1(0.7))                                                     // back
//	s.Add(Sheet(20, Ey), Diffuse1(0.6))                                                     // ceiling
//	slab := Slab(0, 0.7)
//	for j := 0.; j < 2; j++ {
//		for i := j; i < 8; i += 2 {
//			cyl := Cylinder(Vec{j - 3.5, -1.5, i + 4.5}, 0.37)
//			s.Add(ShapeAnd(cyl, slab), Reflective(0.05))
//
//			cyl = Cylinder(Vec{j + 2.5, -1.5, i + 4.5}, 0.37)
//			s.Add(ShapeAnd(cyl, slab), Diffuse1(0.95))
//		}
//	}
//	s.AddLight(PointLight(Vec{3, 12, 6}, 130))
//
//	cam := Camera(1)
//	cam.Transl(Vec{0, 4, 0})
//	cam.Transf(RotX(-15 * deg))
//	t.CompareCam(s, "010-checkers", cam)
//}
//
//func TestDice1(tst *testing.T) {
//	t := Helper(tst)
//
//	s := NewEnv()
//	s.amb = func(Vec) Color { return 0.1 }
//	cube := Box(Vec{0, 0, 0}, -1, -1, -1)
//
//	die := cube
//	const r = 0.15
//	die = ShapeMinus(die, Sphere(Vec{0, 0, -0.9}, r))
//	die = ShapeMinus(die, Sphere(Vec{0.5, 0.5, -0.9}, r))
//	die = ShapeMinus(die, Sphere(Vec{-0.5, 0.5, -0.9}, r))
//	die = ShapeMinus(die, Sphere(Vec{0.5, -0.5, -0.9}, r))
//	die = ShapeMinus(die, Sphere(Vec{-0.5, -0.5, -0.9}, r))
//
//	die = ShapeMinus(die, Sphere(Vec{0.4, 1.1, -0.4}, r))
//	die = ShapeMinus(die, Sphere(Vec{-0.4, 1.1, 0.4}, r))
//
//	die = ShapeAnd(die, Sphere(Vec{}, 0.98*math.Sqrt(2.)))
//
//	s.Add(die, Diffuse1(0.9))
//
//	s.Add(Sheet(-1, Ey), Diffuse1(0.5))
//
//	s.AddLight(PointLight(Vec{2, 3, -3}, 15))
//
//	cam := Camera(1)
//	cam.Transl(Vec{0, 4, -6})
//	cam.Transf(RotX(-15 * deg))
//	t.CompareCam(s, "011-dice", cam)
//}

// Two flat-shaded spheres, partially overlapping.
//func TestOverlap(tst *testing.T) {
//	t := Helper(tst)
//
//	const r = 0.25
//	s := &Env{
//		objs: []Obj{
//			Flat(Sphere(Vec{-r / 2, 0, 3}, r), 1.0),
//			Flat(Sphere(Vec{r / 2, 0, 3}, r), 0.5),
//		},
//	}
//
//	t.Compare(s, "001-overlap")
//}
//
//// A sphere behind the camera, should not be visible
//func TestBehindCam(tst *testing.T) {
//	t := Helper(tst)
//
//	const r = 0.25
//	objects := []Obj{
//		Flat(Sphere(Vec{0, 0, -3}, r), 1),
//	}
//	s := &Env{
//		objs: objects,
//	}
//
//	t.Compare(s, "002-behindcam")
//}

// Intersection of flat-shaded spheres
//func TestIntersect(tst *testing.T) {
//	t := Helper(tst)
//
//	const r = 0.25
//	s1 := &object{Sphere(Vec{-r / 2, 0, 3}, r), Flat(1)}
//	s2 := &object{Sphere(Vec{r / 2, 0, 3}, r), Flat(0.5)}
//	s := NewEnv()
//	s.objs = append(s.objs, &objAnd{s1, s2})
//
//	t.Compare(s, "003-intersect")
//}

//// Intersection of spheres, as shapes (not objects)
//func TestIntersectShape(tst *testing.T) {
//	t := Helper(tst)
//
//	const r = 0.25
//	s1 := Sphere(Vec{-r / 2, 0, 3}, r)
//	s2 := Sphere(Vec{r / 2, 0, 3}, r)
//	sh := ShapeAnd{s1, s2}
//
//	s := &Env{
//		objs: []Obj{
//			Flat(sh, 1),
//		},
//	}
//
//	t.Compare(s, "004-intersectshape")
//}
//
//// Minus of spheres, as shapes (not objects)
//func TestMinusShape(tst *testing.T) {
//	t := Helper(tst)
//
//	const r = 0.5
//	s1 := Sphere(Vec{-r / 2, 0, 3}, r)
//	s2 := Sphere(Vec{r / 2, 0, 3}, r)
//	sh := ShapeMinus{s1, s2}
//	s := &Env{
//		objs: []Obj{
//			Flat(sh, 1),
//		},
//	}
//
//	t.Compare(s, "005-minusshape")
//}
//
//// Intersection of normal.z-shaded spheres
//func TestSphereNormals(tst *testing.T) {
//	t := Helper(tst)
//
//	const r = 0.25
//	s3 := ShadeNormal(Sphere(Vec{0, -0.5, 3}, 2*r))
//	s1 := ShadeNormal(Sphere(Vec{-r / 2, 0, 3}, r))
//	s2 := ShadeNormal(Sphere(Vec{r / 2, 0, 3}, r))
//	s := &Env{
//		objs: []Obj{
//			ObjAnd{s3, s1},
//			ObjAnd{s3, s2},
//		},
//	}
//
//	t.Compare(s, "006-spherenormals")
//}
//
//func TestBox(tst *testing.T) {
//	t := Helper(tst)
//
//	b1 := ABox(Vec{-2, -2, 4}, Vec{-1, -1, 3})
//	b2 := ABox(Vec{2, -2, 4}, Vec{1, -1, 3})
//	b3 := ABox(Vec{-2, 2, 4}, Vec{-1, 1, 3})
//	b4 := ABox(Vec{2, 2, 4}, Vec{1, 1, 3})
//	s := &Env{}
//	s.objs = []Obj{
//		Diffuse1(s, b1, 1),
//		Diffuse1(s, b2, 1),
//		Diffuse1(s, b3, 1),
//		Diffuse1(s, b4, 1),
//	}
//	s.sources = []Source{
//		&PointSource{Vec{0.5, 0.3, -5}, 60},
//	}
//
//	t.CompareCam(s, "007-box", Camera(testW, testH, 0.7))
//}
//
//func TestReflection(tst *testing.T) {
//	t := Helper(tst)
//
//	s := &Env{}
//
//	const h = 2
//	ground := Diffuse1(s, Slab(-h, -h-100), 0.5)
//	sp := Sphere(Vec{-0.5, -1, 8}, 2)
//	die := &ShapeAnd{sp, Slab(-h+.2, -.2)}
//	dice := Diffuse1(s, die, 0.95)
//	s.objs = []Obj{
//		ground,
//		dice,
//		Reflective(s, Sphere(Vec{3, -1, 10}, 1), 0.9),
//	}
//	s.sources = []Source{
//		&PointSource{Vec{6, 10, 2}, 180},
//	}
//	s.amb = func(Vec) Color { return 1 }
//
//	cam := Camera(testW, testH, 1)
//	cam.Transf = RotX(-10 * deg)
//
//	t.CompareCam(s, "008-reflections", cam)
//}

//func (t helper) Compare(s *Env, name string) {
//	t.Helper()
//	cam := Camera(0)
//	t.CompareCam(s, name, cam)
//}

func (t helper) Compare(s *Env, cam *Cam, name string) {
	t.Helper()
	os.Mkdir("out", 0777)

	name = name + ".png"
	have := "out/" + name
	want := "testdata/" + name

	img := MakeImage(testW, testH)
	cam.Render(s, testRec, img)
	Encode(img, have)
	deviation, err := imgComp(have, want)

	if err != nil {
		t.Fatal(err)
	}
	const tolerance = 10
	if deviation > tolerance {
		t.Errorf("%v: differs from reference by %v", name, deviation)
	}
}

func imgComp(a, b string) (float64, error) {
	A, err := imgRead(a)
	if err != nil {
		return 0, err
	}
	B, err := imgRead(b)
	if err != nil {
		return 0, err
	}

	delta := 0
	for y := 0; y < A.Bounds().Max.Y; y++ {
		for x := 0; x < A.Bounds().Max.X; x++ {
			r1, g1, b1, _ := A.At(x, y).RGBA()
			r2, g2, b2, _ := B.At(x, y).RGBA()
			delta += diff(r1, r2) + diff(g1, g2) + diff(b1, b2)
		}
	}
	return float64(delta) / (3 * 255), nil
}

func diff(a, b uint32) int {
	d := int(a) - int(b)
	if d < 0 {
		return -d
	}
	return d
}

func imgRead(fname string) (image.Image, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}
