package transforms

import "math"

func DoG(img [][]Pixel) [][]Pixel {
	blur1 := GaussianBlur1D(img, 13)
	blur2 := GaussianBlur1D(img, 15)

	result := make([][]Pixel, len(blur1))
	for i := range len(blur1) {
		result[i] = make([]Pixel, len(blur1[i]))
	}

	for i := range len(blur1) {
		for j := range len(blur1[i]) {
			pix1 := &blur1[i][j]
			pix2 := &blur2[i][j]

			finRes := max(0, uint8(Luminance(pix1)-Luminance(pix2)))
			if finRes < 200 {
				finRes = 0
			}
			result[i][j] = Pixel{
				R: finRes,
				G: finRes,
				B: finRes,
				A: (pix1.A + pix2.A) / 2,
			}
		}
	}

	return result
}

func XDoG(img [][]Pixel) [][]Pixel {
	blur1 := GaussianBlur1D(img, 7)
	blur2 := GaussianBlur1D(img, 11)

	result := make([][]Pixel, len(blur1))
	for i := range len(blur1) {
		result[i] = make([]Pixel, len(blur1[i]))
	}

	tau := 0.98     // blur constant
	epsilon := 0.08 // threshold value
	phi := 40.0     // edge hardness
	for i := range len(blur1) {
		for j := range len(blur1[i]) {
			pix1 := &blur1[i][j]
			pix2 := &blur2[i][j]
			pix1Lum := Luminance(pix1) / 255.0
			pix2Lum := Luminance(pix2) / 255.0

			finRes := max(0, pix1Lum-tau*pix2Lum) // 0-1 range

			if finRes >= epsilon {
				finRes = 1
			} else {
				finRes = 0.5 * (1 + math.Tanh(phi*(finRes-epsilon)))
			}

			pix_val := (uint8)(finRes * 255)
			result[i][j] = Pixel{
				R: pix_val,
				G: pix_val,
				B: pix_val,
				A: 255,
			}
		}
	}

	return result
}
