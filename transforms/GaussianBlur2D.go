// Reference - https://stackoverflow.com/q
// Posted by gabriel_tiso, modified by community. See post 'Timeline' for change history

package transforms

import (
	"fmt"
	"math"
)

// Creates a Gaussian Blur effect with a kernel_size x kernel_size convolution matrix
func GaussianBlur2D(img [][]Pixel, kernel_size int) [][]Pixel {
	if kernel_size%2 != 1 {
		panic("Enter odd number for kernel_size")
	}
	kernel := createKernel2D(kernel_size)

	blur1 := blur2D(img, kernel)

	fmt.Printf("(%v, %v) \n", len(blur1), len(blur1[0]))
	// for i := range len(blur1) {
	// 	for j := range len(blur1[i]) {
	// 		r, g, b, a := blur1[i][j].Color.RGBA()
	// 		fmt.Printf("(%v, %v, %v, %v) ", r>>8, g>>8, b>>8, a>>8)
	// 	}
	// 	fmt.Println()
	// }

	return blur1
}

func createKernel2D(matrix_size int) [][]float64 {
	kernel := make([][]float64, matrix_size)

	for i := range matrix_size {
		kernel[i] = make([]float64, matrix_size)
	}

	radius := matrix_size / 2
	sigma := float64(matrix_size) / 6.0

	sum := float64(0)
	for i := range matrix_size {
		for j := range matrix_size {
			true_x := i - radius
			true_y := j - radius

			kernel[i][j] = gaussianFunction2D(float64(true_x), float64(true_y), float64(sigma))
			sum += kernel[i][j]
		}
	}

	for i := range matrix_size {
		for j := range matrix_size {
			kernel[i][j] /= sum
		}
	}
	// fmt.Println("kernel:")
	// for i := range matrix_size {
	// 	fmt.Println(kernel[i])
	// }

	return kernel
}

func gaussianFunction2D(x float64, y float64, sigma float64) float64 {
	result := 1 / (2 * math.Pi * sigma * sigma)
	exponent := -(x*x + y*y) / (2 * sigma * sigma)
	result *= math.Exp(exponent)

	return result
}

func blur2D(img [][]Pixel, kernel [][]float64) [][]Pixel {
	img_width := len(img[0])
	img_height := len(img)

	blurred := make([][]Pixel, len(img))

	for i := range len(img) {
		blurred[i] = make([]Pixel, len(img[i]))
	}

	// note:
	// we know that kernel size is always an odd number,
	// so our halves can just be truncated kernel_size / 2
	kernel_size := len(kernel)
	radius := kernel_size / 2

	for i := range img_height {
		min_y_val := -min(i, (radius))
		max_y_val := min(img_height-1-i, (radius))

		for j := range img_width {
			min_x_val := -min(j, radius)
			max_x_val := min(img_width-1-j, radius)

			var r, g, b, a float64 = 0, 0, 0, 0

			// fmt.Printf("Cur coord [%v][%v] -- Image width = %v\n", i, j, img_width)
			// fmt.Printf("Going from [%v][%v] to [%v][%v] -- radius = %v\n", min_x_val, min_y_val, max_x_val, max_y_val, radius)

			for p := min_y_val; p <= max_y_val; p++ {
				for q := min_x_val; q <= max_x_val; q++ {
					pix := img[i+p][j+q]

					rc, rg, rb, ra := pix.R, pix.G, pix.B, pix.A

					// fmt.Println((float64(rc) / raa))
					r += (float64(rc) * kernel[p+radius][q+radius])
					g += (float64(rg) * kernel[p+radius][q+radius])
					b += (float64(rb) * kernel[p+radius][q+radius])
					a += (float64(ra) * kernel[p+radius][q+radius])
				}
			}

			blurred[i][j] = Pixel{
				R:         uint8(r),
				G:         uint8(g),
				B:         uint8(b),
				A:         uint8(a),
				Character: img[i][j].Character,
			}
			// fmt.Printf("New val: %v %v %v %v\n", uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}

	return blurred
}
