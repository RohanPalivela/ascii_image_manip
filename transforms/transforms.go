package transforms

import (
	"image/color"
)

// *****************
// TRANSFORM FILTERS
// *****************

func LuminFilter(p *Pixel, mapping map[int]rune) rune {

	red, green, blue, _ := p.Color.RGBA()
	luminance := float64(0.2126*float64(red)+0.7152*(float64(green))+0.0722*float64(blue)) / 255 // Luminance from 0-255

	// fmt.Println(luminance)

	lumBuckets := min(int(luminance/10), 9) // push into buckets of 0-9 (mapping)

	// fmt.Printf("%c\n", mapping[lumBuckets])

	return mapping[lumBuckets]
}

func Normalize(p *Pixel) color.RGBA {
	red, green, blue, alpha := p.Color.RGBA()
	normalized := uint8(((red) + (blue) + (green)) / 3)
	return color.RGBA{
		R: normalized,
		G: normalized,
		B: normalized,
		A: uint8(alpha),
	}
}
