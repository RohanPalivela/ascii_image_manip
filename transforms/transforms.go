package transforms

import (
	"image/color"
)

// *****************
// TRANSFORM FILTERS
// *****************

func luminize(p *Pixel, mapping map[int]rune) rune {
	red, green, blue := p.R, p.G, p.B
	luminance := float64(0.2126*float64(red)+0.7152*(float64(green))+0.0722*float64(blue)) / 255 // Luminance from 0-255

	// fmt.Println(luminance)

	lumBuckets := min(int(luminance/10), 9) // push into buckets of 0-9 (mapping)

	// fmt.Printf("%c\n", mapping[lumBuckets])

	return mapping[lumBuckets]
}

func LuminFilter(arr [][]Pixel, mapping map[int]rune) {
	for i := range len(arr) {
		for j := range len(arr[i]) {
			cur := &arr[i][j]
			run := luminize(cur, mapping)
			cur.Character = run
			// fmt.Printf("Character: %c, rgb: %s\n", cur.Character, cur.Color)
		}
	}
}

func SobelFilter() {

}

func Normalize(p *Pixel) color.RGBA {
	red, green, blue, alpha := p.R, p.G, p.B, p.A
	normalized := uint8(((red) + (blue) + (green)) / 3)
	return color.RGBA{
		R: normalized,
		G: normalized,
		B: normalized,
		A: uint8(alpha),
	}
}
