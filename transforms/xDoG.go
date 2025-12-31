package transforms

func XDoG(img [][]Pixel) [][]Pixel {
	// blur twice
	// subtract w/ threshold
	// return image
	// DONE!

	blur1 := GaussianBlur1D(img, 5)
	blur2 := GaussianBlur1D(img, 11)

	result := make([][]Pixel, len(blur1))
	for i := range len(blur1) {
		result[i] = make([]Pixel, len(blur1[i]))
	}

	for i := range len(blur1) {
		for j := range len(blur1[i]) {
			pix1 := &blur1[i][j]
			pix2 := &blur2[i][j]
			// fmt.Printf("Getting value for new pix: %v - %v", Luminance(pix1), Luminance(pix2))

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
