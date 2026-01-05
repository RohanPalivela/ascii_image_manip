package transforms

import (
	"math"
	"sync"
)

type CoordPair struct {
	x, y int
}

func SobelFilterConc(arr [][]Pixel, result [][]Pixel, add_character bool, start CoordPair, end CoordPair, Gx [][]float64, Gy [][]float64) {
	var rune_insert rune
	if add_character {
		rune_insert = ' '
	} else {
		rune_insert = rune(0)
	}

	// fmt.Printf("at (%v, %v) to (%v, %v)\n", start.x, start.y, end.x, end.y)
	for i := start.y; i < end.y; i++ {
		if i >= len(arr) {
			break
		}
		for j := start.x; j < end.x; j++ {
			if j >= len(arr[i]) {
				break
			}
			if i == 0 || j == 0 || i == len(arr)-1 || j == len(arr[i])-1 {
				result[i][j] = Pixel{
					R:         0,
					G:         0,
					B:         0,
					A:         255,
					Character: rune_insert,
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

			// fmt.Printf("%v %v\n", x, y)

			angle := math.Mod(math.Atan2(y, x)+math.Pi, math.Pi) // [0, pi]
			r := rune(0)

			switch {
			case !add_character:
				r = rune_insert
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

			if add_character && !(end > 100) {
				r = ' '
			}

			result[i][j] = Pixel{
				R:         end,
				G:         end,
				B:         end,
				A:         255,
				Character: r,
			}

			// fmt.Printf("(%v %v %v %v) ", arr[i][j].R, arr[i][j].G, arr[i][j].B, arr[i][j].Character)
		}
	}

}

func SobelFilter(arr [][]Pixel, add_character bool) [][]Pixel {
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

	result := make([][]Pixel, len(arr))

	for i := range len(arr) {
		result[i] = make([]Pixel, len(arr[i]))
	}

	divisions := 10 // 8 threads ? ish
	var group sync.WaitGroup

	incr := len(arr) / divisions

	// vertical partitions
	for i := 0; i < len(arr); i += incr {
		start := CoordPair{
			x: i,
			y: 0,
		}
		end := CoordPair{
			x: min(i + incr),
			y: len(arr) - 1,
		}

		group.Go(func() {
			(SobelFilterConc(arr, result, add_character, start, end, Gx, Gy))
		})
	}

	group.Wait()

	return result
}
