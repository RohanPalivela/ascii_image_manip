package transforms

import (
	"fmt"
	"image"
	"image/color"

	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"
)

type AsciiImageBuffer struct {
	x, y          int
	width, height int
	letter_size   int // a letter will take up a letter_size x letter_size amount of space (i.e. 4x4 space for each character)
	img           *image.RGBA
}

// Initializes a new AsciiImageBuffer
func InitializeBuffer(x int, y int, width int, height int, letter_size int, img *image.RGBA) (buffer *AsciiImageBuffer) {
	return &AsciiImageBuffer{x, y, width, height, letter_size, img}
}

/* Writes rune to the Context provided. AsciiImageBuffer keeps track of the current position, does wrapping for you.
 */
func (buffer *AsciiImageBuffer) WriteRune(context *freetype.Context, c color.Color, r rune) error {
	if buffer.x >= buffer.width {
		buffer.x = 0
		buffer.y += buffer.letter_size
	}

	if buffer.y >= buffer.height {
		return fmt.Errorf("draw string overflow, y height is %v", buffer.y)
	}

	// fmt.Println("Drawing " + (string(r)))
	context.SetSrc(image.NewUniform(c))
	if r == 0 {
		for i := buffer.x; i < buffer.x+buffer.letter_size; i++ {
			// we draw from bottom left, so translate up one letter size
			for j := buffer.y - buffer.letter_size; j < buffer.y+buffer.letter_size; j++ {
				buffer.img.Set(i, j, c)
			}
		}
		buffer.x += buffer.letter_size
		return nil
	}

	_, err := context.DrawString(string(r), fixed.P(buffer.x, buffer.y))

	if err != nil {
		return err
	}

	buffer.x += buffer.letter_size

	return nil
}

/* Writes array to the Context provided. AsciiImageBuffer keeps track of the current position, does wrapping for you.
 */
func (buffer *AsciiImageBuffer) WriteArray(context *freetype.Context, arr [][]Pixel) {
	for i := range len(arr) {
		for j := range len(arr[i]) {
			cur := &arr[i][j]
			// fmt.Printf("Character: %c, rgb: %s\n", cur.Character, cur.Color)
			if err := buffer.WriteRune(context, color.RGBA{cur.R, cur.G, cur.B, cur.A}, cur.Character); err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}
