package transforms

import (
	"math"
)

func GaussianBlur1D(arr [][]Pixel, kernel_size int) [][]Pixel {
	newarr := blur(arr, kernel_size)

	return newarr
}

func gaussianFunction1D(x float64, sigma float64) float64 {
	result := 1 / math.Sqrt(2*math.Pi*sigma*sigma)
	exponent := -(x * x) / (2 * sigma * sigma)
	result *= math.Exp(exponent)

	return result
}

func gausKernel1D(kernel_size int) []float64 {
	if kernel_size%2 == 0 || kernel_size < 1 {
		panic("Enter valid kernel_size")
	}

	sigma := float64(kernel_size) / 6
	radius := kernel_size / 2
	kernel := make([]float64, kernel_size)

	sum := float64(0)
	for i := range kernel_size {
		kernel[i] = gaussianFunction1D(float64(i-radius), sigma)
		sum += kernel[i]
	}

	for i := range kernel_size {
		kernel[i] /= sum
	}

	return kernel
}

func blur(arr [][]Pixel, kernel_size int) [][]Pixel {
	kernel := gausKernel1D(kernel_size)
	radius := kernel_size / 2

	result := make([][]Pixel, len(arr))

	for i := range len(arr) {
		result[i] = make([]Pixel, len(arr[i]))
	}

	// vertical pass
	for i := range len(arr) {
		for j := range len(arr[i]) {
			var r, g, b, a float64 = 0, 0, 0, 0
			for k := -radius; k <= radius; k++ {
				if i+k < 0 || i+k > len(arr)-1 {
					break
				}
				pix := arr[i+k][j]

				weight := kernel[k+radius]

				rc, rg, rb, ra := pix.R, pix.G, pix.B, pix.A

				r += (float64(rc) * weight)
				g += (float64(rg) * weight)
				b += (float64(rb) * weight)
				a += (float64(ra) * weight)
			}

			result[i][j] = Pixel{
				R:         uint8(r),
				G:         uint8(g),
				B:         uint8(b),
				A:         uint8(a),
				Character: arr[i][j].Character,
			}

			// fmt.Printf("New val: %v %v %v %v\n", uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}

	// horizontal pass
	for i := range len(result) {
		for j := range len(result[i]) {
			var r, g, b, a float64 = 0, 0, 0, 0
			for k := -radius; k <= radius; k++ {
				if j+k < 0 || j+k > len(result[i])-1 {
					break
				}
				pix := result[i][j+k]

				rc, rg, rb, ra := pix.R, pix.G, pix.B, pix.A

				r += (float64(rc) * kernel[k+radius])
				g += (float64(rg) * kernel[k+radius])
				b += (float64(rb) * kernel[k+radius])
				a += (float64(ra) * kernel[k+radius])
			}

			result[i][j].R = uint8(r)
			result[i][j].G = uint8(g)
			result[i][j].B = uint8(b)
			result[i][j].A = uint8(a)

			// fmt.Printf("New val: %v %v %v %v\n", uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}

	return result
}
