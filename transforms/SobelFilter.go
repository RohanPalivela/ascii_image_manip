package transforms

import (
	"math"
)

func SobelFilter(arr [][]Pixel) [][]Pixel {
	result := make([][]Pixel, len(arr))

	for i := range len(arr) {
		result[i] = make([]Pixel, len(arr[i]))
	}

	Gx := [][]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	Gy := [][]float64{
		{1, 2, 1},
		{0, 0, 0},
		{-1, -2, -1},
	}

	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			if i == 0 || j == 0 || i == len(arr)-1 || j == len(arr[i])-1 {
				result[i][j] = Pixel{
					R:         0,
					G:         0,
					B:         0,
					A:         255,
					Character: ' ',
				}
				continue
			}
			x := float64(0.0)
			y := float64(0.0)
			for k := -1; k <= 1; k++ {
				for l := -1; l <= 1; l++ {
					x += Gx[k+1][l+1] * Luminance(&arr[i+k][j+l])
					y += Gy[k+1][l+1] * Luminance(&arr[i+k][j+l])
				}
			}

			angle := math.Mod(math.Atan2(y, x)+math.Pi, math.Pi) // [0, pi]
			r := rune(0)

			switch {
			case angle < math.Pi/8 || angle >= 7*math.Pi/8:
				r = '|'
			case angle < 3*math.Pi/8:
				r = '\\'
			case angle < 5*math.Pi/8:
				r = '_'
			default:
				r = '/'
			}

			end := uint8(min(255, math.Abs(x)+math.Abs(y)))

			if !(end > 100) {
				r = ' '
			}

			result[i][j] = Pixel{
				R:         end,
				G:         end,
				B:         end,
				A:         255,
				Character: r,
			}
		}
	}

	return result
}
