package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"
)

type Pixel struct {
	red   uint8
	blue  uint8
	green uint8
	alpha uint8
	x     int
	y     int
}

type AsciiImageBuffer struct {
	x           int
	y           int
	width       int
	height      int
	letter_size int // a letter will take up a letter_size x letter_size amount of space (i.e. 4x4 space for each character)
}

func (buffer *AsciiImageBuffer) WriteRune(r rune, context *freetype.Context) (point fixed.Point26_6, err error) {
	if buffer.x >= buffer.width {
		buffer.x = 0
		buffer.y += buffer.letter_size
	}

	if buffer.y >= buffer.height {
		return fixed.Point26_6{}, fmt.Errorf("draw string overflow, y height is %v", buffer.y)
	}

	pt, err := context.DrawString(string(r), fixed.P(buffer.x, buffer.y))

	if err != nil {
		return fixed.Point26_6{}, err
	}

	buffer.x += buffer.letter_size

	// fmt.Println("Drew " + string(r) + " at " + pt.X.String() + ", " + pt.Y.String())

	return pt, nil
}

func main() {
	img, bounds := OpenJPEGIMG("bg.jpg")

	width := bounds.Dx()
	height := bounds.Dy()

	// keeping these the same value yields an image of ~ same size
	sample_size := 1
	px_size := 16

	fmt.Printf("Dimensions: %v x %v\n", width, height)

	pix_width := width / sample_size
	pix_height := height / sample_size
	fmt.Printf("Shrinking to: %v x %v\n", pix_width, pix_height)

	arr := make([][]Pixel, pix_height)

	// luminescence to ascii mapping
	mapping := map[int]rune{
		0: ' ',
		1: '.',
		2: ':',
		3: 'c',
		4: 'o',
		5: 'C',
		6: 'O',
		7: '0',
		8: '@',
		9: '■',
	}

	for y := range pix_height {
		arr[y] = make([]Pixel, pix_width)
	}

	// consolidate a pixel grid of size sample_size x sample_size into one pixel
	for by := range pix_height {
		for bx := range pix_width {
			x := bounds.Min.X + bx*sample_size
			y := bounds.Min.Y + by*sample_size

			red := uint32(0)
			green := uint32(0)
			blue := uint32(0)
			alpha := uint32(0)

			sample_count := 0

			for offset_x := range sample_size {
				if x+offset_x >= bounds.Max.X {
					break
				}
				for offset_y := 0; offset_y < sample_size; offset_y++ {
					if y+offset_y >= bounds.Max.Y {
						break
					}
					r, g, b, a := img.At(x+offset_x, y+offset_y).RGBA()
					red += r >> 8
					green += g >> 8
					blue += b >> 8
					alpha += a >> 8
					sample_count++
				}
			}

			red /= uint32(sample_count)
			green /= uint32(sample_count)
			blue /= uint32(sample_count)
			alpha /= uint32(sample_count)

			// fmt.Printf("Pixel: (%v, %v, %v, %v)\n", r, g, b, a)
			arr[by][bx] =
				Pixel{
					red:   uint8(red),
					green: uint8(green),
					blue:  uint8(blue),
					alpha: uint8(alpha),
					x:     x,
					y:     y,
				}
		}
	}

	// BOUNDS FOR KEEPING THE IMAGE QUALITY PERFECT:
	out_width := pix_width * px_size
	out_height := pix_height * px_size

	newimg := image.NewRGBA(image.Rect(0, 0, out_width, out_height))
	draw.Draw(newimg, newimg.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)

	fontBytes, err := os.ReadFile("BoldPixelsFont.ttf")
	if err != nil {
		log.Println(err)
		return
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	c := freetype.NewContext()

	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(float64(px_size)) // 3pt font apparently translates to 4pixels
	c.SetClip(newimg.Bounds())
	c.SetDst(newimg)
	c.SetSrc(image.White)

	buffer := AsciiImageBuffer{x: 0, y: px_size, width: out_width, height: out_height, letter_size: px_size}

	for i := range pix_height {
		for j := range pix_width {
			cur := arr[i][j]
			// sb.WriteRune(LuminFilter(&cur, mapping))
			if _, err := buffer.WriteRune(LuminFilter(&cur, mapping), c); err != nil {
				// fmt.Println(err.Error())
				break
			}
		}
	}

	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()

	err = png.Encode(outFile, newimg)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("Created new image")
}

// *****************
// TRANSFORM FILTERS
// *****************

func LuminFilter(p *Pixel, mapping map[int]rune) rune {

	luminance := float64(0.2126*float64(p.red)+0.7152*(float64(p.green))+0.0722*float64(p.blue)) / 255 // Luminance from 0-255

	lumBuckets := min(int(luminance*10), 9) // push into buckets of 0-9 (mapping)

	return mapping[lumBuckets]
}

func Normalize(p *Pixel) color.RGBA {
	normalized := uint8(((p.red) + (p.blue) + (p.green)) / 3)
	return color.RGBA{
		R: normalized,
		G: normalized,
		B: normalized,
		A: uint8(p.alpha),
	}
}

// func WriteToTXT(height int, width int) {
// 	// .txt output
// 	var sb strings.Builder
// 	for i := range height {
// 		for j := range width {
// 			cur := arr[i][j]
// 			sb.WriteRune(LuminFilter(&cur, mapping))
// 			sb.WriteString(" ")
// 		}
// 		sb.WriteString("\n")
// 	}

// 	file, err := os.Create("output.txt")

// 	if err != nil {
// 		log.Fatal("Couldn't create output file " + err.Error())
// 	}

// 	defer file.Close()

// 	file.WriteString(sb.String())
// }

// *****************
// IO OPERATIONS
// *****************

func OpenJPEGIMG(filename string) (image image.Image, bounding image.Rectangle) {
	img, err := os.Open(filename)

	if err != nil {
		log.Fatal("Error opening file: " + err.Error())
	}

	fmt.Println("Opened " + img.Name())

	defer img.Close()

	m, err := jpeg.Decode(img)

	if err != nil {
		log.Fatal("Could not decode image: " + err.Error())
	}

	bounds := m.Bounds()

	return m, bounds
}

func CreatePNG(filename string, newimg image.Image) (output string, err error) {
	name := filename + ".png"

	f, err := os.Create(name)

	if err != nil {
		return "", err
	}

	defer f.Close()

	err = png.Encode(f, newimg)

	if err != nil {
		return "", err
	}

	return name, nil
}

func CreateJPEG(filename string, newimg image.Image, quality int) (output string, err error) {
	name := filename + ".jpeg"

	f, err := os.Create(name)

	if err != nil {
		return "", err
	}

	defer f.Close()

	err = jpeg.Encode(f, newimg, &jpeg.Options{Quality: quality})

	if err != nil {
		return "", err
	}

	return name, nil
}
