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
	img, bounds := OpenJPEGIMG("wolf.jpeg")

	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("Dimensions: %v x %v\n", width, height)

	arr := make([][]Pixel, height)

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

	for y := range height {
		arr[y] = make([]Pixel, width)
	}

	count := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// fmt.Printf("Pixel: (%v, %v, %v, %v)\n", r, g, b, a)
			arr[count/width][count%width] =
				Pixel{
					red:   uint8(r >> 8),
					green: uint8(g >> 8),
					blue:  uint8(b >> 8),
					alpha: uint8(a >> 8),
					x:     x,
					y:     y,
				}
			count++
		}
	}

	// IMAGE CREATION CODE
	// {
	//
	// 	for i := 0; i < len(arr); i++ {
	// 		for j := 0; j < len(arr[i]); j++ {
	// 			cur := arr[i][j]
	// 			pixe := Normalize(&cur)
	// 			newimg1.Set(cur.x, cur.y, pixe)
	// 		}
	// 		fmt.Printf("\n")
	// 	}
	//
	// 	endProd, err := CreateJPEG("output1", newimg1, 100)
	//
	// 	if err != nil {
	// 		log.Fatal("Issue creating file: " + err.Error())
	// 	}
	//
	// 	fmt.Printf("Created new image %v\n", endProd)
	// }

	// .txt output
	// var sb strings.Builder
	// for i := range height {
	// 	for j := range width {
	// 		cur := arr[i][j]
	// 		sb.WriteRune(LuminFilter(&cur, mapping))
	// 		sb.WriteString(" ")
	// 	}
	// 	sb.WriteString("\n")
	// }
	//
	// file, err := os.Create("output.txt")
	//
	// if err != nil {
	// 	log.Fatal("Couldn't create output file " + err.Error())
	// }
	//
	// defer file.Close()
	//
	// file.WriteString(sb.String())

	// BOUNDS FOR KEEPING THE IMAGE QUALITY PERFECT:
	px_size := 16
	newimg := image.NewRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X*px_size, bounds.Max.Y*px_size))
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
	c.SetFont(f)                    // NEED FONT
	c.SetFontSize(float64(px_size)) // 3pt font apparently translates to 4pixels
	c.SetClip(newimg.Bounds())
	c.SetDst(newimg)
	c.SetSrc(image.White)

	buffer := AsciiImageBuffer{x: 0, y: px_size, width: width * px_size, height: height * px_size, letter_size: px_size}

	for i := range height {
		for j := range width {
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
