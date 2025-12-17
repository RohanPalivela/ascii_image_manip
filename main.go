package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

type Pixel struct {
	red   uint8
	blue  uint8
	green uint8
	alpha uint8
	x     int
	y     int
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

	mapping[1] = 'd'

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
	// newimg1 := image.NewRGBA(image.Rect(bounds.Min.X, bounds.Min.X, bounds.Max.X, bounds.Max.Y))

	// for i := 0; i < len(arr); i++ {
	// 	for j := 0; j < len(arr[i]); j++ {
	// 		cur := arr[i][j]
	// 		pixe := Normalize(&cur)
	// 		sb.WriteRune(LuminFilter(&cur))
	// 		newimg1.Set(cur.x, cur.y, pixe)
	// 	}
	// 	fmt.Printf("\n")
	// }

	// endProd, err := CreateJPEG("output1", newimg1, 100)

	// if err != nil {
	// 	log.Fatal("Issue creating file: " + err.Error())
	// }

	// fmt.Printf("Created new image %v\n", endProd)
	var sb strings.Builder

	for i := range height {
		for j := range width {
			cur := arr[i][j]
			sb.WriteRune(LuminFilter(&cur, mapping))
		}
		sb.WriteString("\n")
	}

	file, err := os.Create("output.txt")

	if err != nil {
		log.Fatal("Couldn't create output file " + err.Error())
	}

	defer file.Close()

	file.WriteString(sb.String())

	fmt.Println("Created new image")
}

func LuminFilter(p *Pixel, mapping map[int]rune) rune {

	luminance := float64(0.2126*float64(p.red)+0.7152*(float64(p.green))+0.0722*float64(p.blue)) / 255 // Luminance from 0-255

	lumBuckets := min(int(luminance*10), 9) // push into buckets of 0-9 (mapping)

	return mapping[lumBuckets]
}

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

func Normalize(p *Pixel) color.RGBA {
	normalized := uint8(((p.red >> 8) + (p.blue >> 8) + (p.green >> 8)) / 3)
	return color.RGBA{
		R: normalized,
		G: normalized,
		B: normalized,
		A: uint8(p.alpha >> 8),
	}
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
