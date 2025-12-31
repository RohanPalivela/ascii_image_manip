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
	tmp := make([][]Pixel, len(arr))

	for i := range len(arr) {
		result[i] = make([]Pixel, len(arr[i]))
		tmp[i] = make([]Pixel, len(arr[i]))
	}

	// horizontal pass
	for i := range len(arr) {
		for j := range len(arr[i]) {
			var r, g, b, a float64 = 0, 0, 0, 0
			for k := -radius; k <= radius; k++ {
				var pix Pixel
				if j+k < 0 {
					pix = arr[i][0]
				} else if j+k > len(result[i])-1 {
					pix = arr[i][len(result[i])-1]
				} else {
					pix = arr[i][j+k]
				}
				weight := kernel[k+radius]
				rc, rg, rb, ra := pix.R, pix.G, pix.B, pix.A
				r += (float64(rc) * weight)
				g += (float64(rg) * weight)
				b += (float64(rb) * weight)
				a += (float64(ra) * weight)
			}
			tmp[i][j] = Pixel{
				R:         uint8(r),
				G:         uint8(g),
				B:         uint8(b),
				A:         uint8(a),
				Character: arr[i][j].Character,
			}
		}
	}

	// vertical pass
	for i := range len(tmp) {
		for j := range len(tmp[i]) {
			var r, g, b, a float64 = 0, 0, 0, 0
			for k := -radius; k <= radius; k++ {
				var pix Pixel
				if i+k < 0 {
					pix = tmp[0][j]
				} else if i+k > len(tmp)-1 {
					pix = tmp[len(tmp)-1][j]
				} else {
					pix = tmp[i+k][j]
				}
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
		}
	}

	return result
}
