package ascii_img

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	transforms "github.com/RohanPalivela/ascii_image_manip/transforms"
	"github.com/golang/freetype"
)

// Place your image (future: supporting video) into the "images/" directory. Provide the filename (ex: "hi.png" with no images/) into this function. A sample size N averages every NxN space, downscaling the image by Nx.
func Initialize(filename string, sample_size int) [][]transforms.Pixel {
	var img image.Image
	var bounds image.Rectangle

	switch filename[strings.Index(filename, "."):] {
	case ".png":
		img, bounds = OpenPNGImg(filename)
	case ".jpeg", ".jpg":
		img, bounds = OpenJPEGImg(filename)
	}

	width := bounds.Dx()
	height := bounds.Dy()

	pix_width := width / sample_size
	pix_height := height / sample_size

	arr := InitializeArray(img, sample_size, pix_height, pix_width)

	return arr
}

func OutputImage(arr [][]transforms.Pixel, px_size int, color_image bool) *image.RGBA {
	pix_width := len(arr[0])
	pix_height := len(arr)

	// BOUNDS FOR KEEPING THE IMAGE QUALITY PERFECT:
	out_width := pix_width * px_size
	out_height := pix_height * px_size

	newimg := image.NewRGBA(image.Rect(0, 0, out_width, out_height))
	draw.Draw(newimg, newimg.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)

	context := InitializeContext(newimg, float64(px_size))

	buffer := transforms.InitializeBuffer(0, px_size, out_width, out_height, px_size, newimg)

	buffer.WriteArray(context, arr, color_image)

	return newimg
}

func WriteToTXT(arr [][]transforms.Pixel) {
	// .txt output
	var sb strings.Builder
	for i := range len(arr) {
		for j := range len(arr[i]) {
			cur := arr[i][j]
			sb.WriteRune(cur.Character)
			sb.WriteRune(' ')
		}
		sb.WriteRune('\n')
	}

	file, err := os.Create("output.txt")

	if err != nil {
		log.Fatal("Couldn't create output file " + err.Error())
	}

	defer file.Close()

	file.WriteString(sb.String())
}

// *****************
// IO OPERATIONS
// *****************

func OpenPNGImg(filename string) (image image.Image, bounding image.Rectangle) {
	img, err := os.Open(filename)

	if err != nil {
		log.Fatal("Error opening file: " + err.Error())
	}

	defer img.Close()

	m, err := png.Decode(img)

	if err != nil {
		log.Fatal("Could not decode image: " + err.Error())
	}

	bounds := m.Bounds()

	return m, bounds
}

func OpenJPEGImg(filename string) (image image.Image, bounding image.Rectangle) {
	img, err := os.Open(filename)

	if err != nil {
		log.Fatal("Error opening file: " + err.Error())
	}

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

func InitializeContext(newimg draw.Image, px_size float64) (cont *freetype.Context) {
	fontBytes, err := os.ReadFile("Fonts/MC.ttf")
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
	c.SetFontSize(px_size)
	c.SetClip((newimg).Bounds())
	c.SetDst(newimg)
	c.SetSrc(image.White) // default value ig

	return c
}

func InitializeArray(img image.Image, sample_size int, pix_height int, pix_width int) (pixels [][]transforms.Pixel) {
	arr := make([][]transforms.Pixel, pix_height)
	for y := range pix_height {
		arr[y] = make([]transforms.Pixel, pix_width)
	}

	bounds := img.Bounds()
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
			// fmt.Printf("Pixel: (%v, %v, %v, %v)\n", (red), (green), (blue), alpha)
			red /= uint32(sample_count)
			green /= uint32(sample_count)
			blue /= uint32(sample_count)
			alpha /= uint32(sample_count)
			// fmt.Printf("Pixel: (%v, %v, %v, %v)\n", uint8(red), uint8(green), uint8(blue), alpha)
			arr[by][bx] =
				transforms.Pixel{
					R: uint8(red),
					G: uint8(green),
					B: uint8(blue),
					A: uint8(alpha),
				}
		}
	}

	return arr
}

func GetRunes(arr [][]transforms.Pixel) {
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

	transforms.LuminFilter(arr, mapping)
}
