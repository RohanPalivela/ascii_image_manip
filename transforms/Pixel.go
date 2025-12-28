package transforms

import "image/color"

type Pixel struct {
	Color     color.Color
	X, Y      int
	Character rune
}
