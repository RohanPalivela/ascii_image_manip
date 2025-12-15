package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

type Pixel struct {
	red   uint32
	blue  uint32
	green uint32
	alpha uint32
	x     int
	y     int
}

func main() {
	img, err := os.Open("wolf.jpeg")

	if err != nil {
		log.Fatal("Error opening file: " + err.Error())
	}

	fmt.Println("Opened " + img.Name())

	defer img.Close()

	m, err := jpeg.Decode(img)

	if err != nil {
		log.Fatal("Could not decode image: " + err.Error())
	}

	// code for like something idk
	// reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	// m, _, err := image.Decode(reader)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	bounds := m.Bounds()

	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("Dimensions: %v x %v\n", width, height)

	arr := make([][]Pixel, height)

	for y := range height {
		arr[y] = make([]Pixel, width)
	}

	count := 0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := m.At(x, y).RGBA()

			// fmt.Printf("Pixel: (%v, %v, %v, %v)\n", r, g, b, a)
			arr[count/height][count%height] = Pixel{r, g, b, a, x, y}
			count++
		}
	}

	newimg := image.NewRGBA(image.Rect(bounds.Min.X, bounds.Min.X, bounds.Max.X, bounds.Max.Y))

	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			cur := arr[i][j]

			// fmt.Printf("Pixel: (%v, %v, %v, %v)\n", cur.red, cur.blue, cur.green, cur.alpha)

			// fmt.Printf("Pixel shifted: (%v, %v, %v, %v)\n\n", cur.red>>8, cur.blue>>8, cur.green>>8, cur.alpha>>8)
			pixe := PixelManip(&cur)
			newimg.Set(cur.x, cur.y, pixe)
		}
	}

	filename := "output"

	endProd, err := CreateJPEG(filename, newimg, 100)

	if err != nil {
		log.Fatal("Issue creating file: " + err.Error())
	}

	fmt.Printf("Created new image %v", endProd)
}

func PixelManip(p *Pixel) color.RGBA {
	return color.RGBA{
		// R: uint8(p.red >> 8),
		R: uint8(0),
		G: uint8(p.blue >> 8),
		B: uint8(p.green >> 8),
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
