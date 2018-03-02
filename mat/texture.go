package mat

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"

	. "github.com/barnex/bruteray/br"
	"github.com/barnex/bruteray/raster"
)

func Texture(img raster.Image, p0, px, py Vec) *texture {
	return &texture{img, p0, px, py}
}

type texture struct {
	img        raster.Image
	p0, px, py Vec
}

func (c *texture) Shade(ctx *Ctx, e *Env, N int, r *Ray, frag Fragment) Color {
	pos := r.At(frag.T)

	p := pos.Sub(c.p0)
	px := c.px.Sub(c.p0)
	py := c.py.Sub(c.p0)

	x := p.Dot(px) / px.Len()
	y := p.Dot(py) / py.Len()

	w := c.img.Bounds().Dx()
	h := c.img.Bounds().Dy()

	ix := int(x * float64(w))
	iy := int(y * float64(h))

	ix = clamp(ix, w)
	iy = clamp(iy, h)
	return c.img[iy][ix]
}

func clamp(v, max int) int {
	if v < 0 {
		return 0
	}
	if v >= max {
		return max - 1
	}
	return v
}

func MustLoad(name string) raster.Image {
	img, err := Load(name)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func Load(name string) (raster.Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	srgb, _, err := image.Decode(f)

	w := srgb.Bounds().Dx()
	h := srgb.Bounds().Dy()
	img := raster.MakeImage(w, h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := srgb.At(x, y).RGBA()
			Y := h - 1 - y
			img[Y][x] = Color{
				linear(float64(r) / 0xffff),
				linear(float64(g) / 0xffff),
				linear(float64(b) / 0xffff),
			}
		}
	}
	return img, err
}

// sRGB to linear conversion
// https://en.wikipedia.org/wiki/SRGB
func linear(s float64) float64 {
	if s <= 0.04045 {
		return s / 12.92
	}
	const a = 0.055
	return math.Pow((s+a)/(1+a), 2.4)
}
