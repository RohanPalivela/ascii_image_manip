// Reference - https://stackoverflow.com/q
// Posted by gabriel_tiso, modified by community. See post 'Timeline' for change history

package transforms

func xDoG(img [][]Pixel, threshold float64, kernel_size int) [][]Pixel {
	// blur twice
	// subtract w/ threshold
	// return image
	// DONE!

	return nil
}

func GaussianBlur(img [][]Pixel, kernel_size int) [][]Pixel {
	kernel := createKernel(kernel_size)

	blur1 := blur(img, kernel)

	return blur1
}

func createKernel(matrix_size int) [][]float64 {
	kernel := make([][]float64, matrix_size)

	for i := range matrix_size {
		kernel[i] = make([]float64, matrix_size)
	}

	return kernel
}

func blur(img [][]Pixel, kernel [][]float64) [][]Pixel {
	img_width := len(img[0])
	img_height := len(img)

	blurred := make([][]Pixel, len(img))

	for i := range len(img) {
		blurred[i] = make([]Pixel, len(img[i]))
	}

	kernel_size := len(kernel) // we know that kernel size is always an odd number

	for i := 0; i < img_height; i++ {
		for j := 0; j < img_width; j++ {
			for i := -kernel_size; i <= kernel_size; i++ {

			}
		}
	}

	return blurred
}

func subtract(blur1 [][]Pixel, blur2 [][]Pixel, threshold int) [][]Pixel {
	return nil
}
