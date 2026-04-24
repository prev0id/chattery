package identicon

import (
	"bytes"
	"crypto/md5"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
)

func hslToRGBA(h, s, l float64) color.RGBA {
	var hue, sat, lum, a, b float64
	hue = h / 360.0
	sat = s / 100.0
	lum = l / 100.0
	if lum < 0.5 {
		b = lum * (sat + 1.0)
	} else {
		b = lum + sat - lum*sat
	}
	a = lum*2.0 - b
	return color.RGBA{
		uint8(math.Round(255 * hueToRgb(a, b, hue+1.0/3.0))),
		uint8(math.Round(255 * hueToRgb(a, b, hue))),
		uint8(math.Round(255 * hueToRgb(a, b, hue-1.0/3.0))),
		255,
	}
}

func hueToRgb(a, b, hue float64) float64 {
	var h = hue
	if hue < 0.0 {
		h = hue + 1.0
	} else if hue > 1.0 {
		h = hue - 1.0
	}
	if h < 1.0/6.0 {
		return a + (b-a)*6.0*h
	}
	if h < 1.0/2.0 {
		return b
	}
	if h < 2.0/3.0 {
		return a + (b-a)*(2.0/3.0-h)*6.0
	}
	return a
}

func remap(value, vmin, vmax, dmin, dmax uint32) float64 {
	return float64((value-vmin)*(dmax-dmin)) / float64((vmax-vmin)+dmin)
}

func foreground(hash [16]byte) color.RGBA {
	h := uint32(hash[12]&15)<<8 | uint32(hash[13])
	s := uint32(hash[14])
	l := uint32(hash[15])
	hue := remap(h, 0, 4095, 0, 360)
	sat := remap(s, 0, 255, 0, 20)
	lum := remap(l, 0, 255, 0, 20)
	return hslToRGBA(hue, 65.0-sat, 75.0-lum)
}

func matrix(hash [16]byte) [5][5]bool {
	tmp := [16]bool{}
	res := [5][5]bool{}

	for i := range 8 {
		tmp[i*2] = (hash[i]>>4)&1 == 0
		tmp[i*2+1] = hash[i]&1 == 0
	}
	for j := range 3 {
		for i := range 5 {
			res[i][2+j] = tmp[j*5+i]
			res[i][2-j] = tmp[j*5+i]
		}
	}
	return res
}

func GenerateImage(data []byte, size int) image.Image {
	hash := md5.Sum(data)
	mat := matrix(hash)
	fg := foreground(hash)
	bg := color.RGBA{240, 240, 240, 255}
	res := image.NewRGBA(image.Rect(0, 0, 6*size, 6*size))
	draw.Draw(res, res.Bounds(), &image.Uniform{bg}, image.Point{0, 0}, draw.Src)
	for j := range 5 {
		for i := range 5 {
			if mat[j][i] {
				x, y := size/2+i*size, size/2+j*size
				draw.Draw(
					res,
					image.Rect(x, y, x+size, y+size),
					&image.Uniform{fg},
					image.Point{0, 0},
					draw.Src,
				)
			}
		}
	}
	return res
}

func Bytes(data []byte, size int) ([]byte, error) {
	img := GenerateImage(data, size)
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
