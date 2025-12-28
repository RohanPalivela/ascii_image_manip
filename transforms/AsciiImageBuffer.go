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
}

// Initializes a new AsciiImageBuffer
func InitializeBuffer(x int, y int, width int, height int, letter_size int) (buffer *AsciiImageBuffer) {
	return &AsciiImageBuffer{x, y, width, height, letter_size}
}

/* Writes rune to the Context provided. AsciiImageBuffer keeps track of the current position, does wrapping for you.
 */
func (buffer *AsciiImageBuffer) WriteRune(context *freetype.Context, c color.Color, r rune) (point fixed.Point26_6, err error) {
	if buffer.x >= buffer.width {
		buffer.x = 0
		buffer.y += buffer.letter_size
	}

	if buffer.y >= buffer.height {
		return fixed.Point26_6{}, fmt.Errorf("draw string overflow, y height is %v", buffer.y)
	}

	// fmt.Println("Drawing " + (string(r)))
	context.SetSrc(image.NewUniform(c))
	pt, err := context.DrawString(string(r), fixed.P(buffer.x, buffer.y))

	if err != nil {
		return fixed.Point26_6{}, err
	}

	buffer.x += buffer.letter_size

	return pt, nil
}
