package ascii_img

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	transforms "github.com/RohanPalivela/ascii_image_manip/transforms"
	"github.com/golang/freetype"
)

// Place your image/video into the "images/" directory. Provide the filename (ex: "hi.png" with no images/) into this function.
func Initialize(filename string, sample_size int, px_size int) [][]transforms.Pixel {
	// file_name := "portrait.jpg"
	img_file := "Images/" + filename

	fmt.Println("** BEGINNING OPERATIONS **")
	start := time.Now()

	var img image.Image
	var bounds image.Rectangle

	switch img_file[strings.Index(img_file, "."):] {
	case ".png":
		img, bounds = OpenPNGImg(img_file)
	case ".jpeg", ".jpg":
		img, bounds = OpenJPEGImg(img_file)
	}

	LogOut(fmt.Sprintf("LOGGING >> Took %s to open image", time.Since(start)))
	intermediate := time.Now()

	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("Info: Original Dimensions: %v x %v\n", width, height)

	pix_width := width / sample_size
	pix_height := height / sample_size
	fmt.Printf("Info: Shrinking to: %v x %v\n\n", pix_width, pix_height)

	arr := InitializeArray(img, sample_size, pix_height, pix_width)

	LogOut(fmt.Sprintf("LOGGING >> Took %s to make and propagate pixel info", time.Since(intermediate)))
	intermediate = time.Now()

	return arr
}

func OutputImage(arr [][]transforms.Pixel) {

}

func main() {
	// CONSTANTS FOR PROGRAM
	file_name := "portrait.jpg"
	img_file := "Images/" + file_name

	// IMPORTANT FOR IMAGE SIZES ************
	// keeping these the same value yields an image of ~ same size
	// for any image > 1920x1080, keep at 8x8, otherwise probably go down
	sample_size := 4
	px_size := 4

	fmt.Println("** BEGINNING OPERATIONS **")
	start := time.Now()

	var img image.Image
	var bounds image.Rectangle

	switch img_file[strings.Index(img_file, "."):] {
	case ".png":
		img, bounds = OpenPNGImg(img_file)
	case ".jpeg", ".jpg":
		img, bounds = OpenJPEGImg(img_file)
	}

	LogOut(fmt.Sprintf("LOGGING >> Took %s to open image", time.Since(start)))
	intermediate := time.Now()

	operations := time.Now()

	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("Info: Original Dimensions: %v x %v\n", width, height)

	pix_width := width / sample_size
	pix_height := height / sample_size
	fmt.Printf("Info: Shrinking to: %v x %v\n\n", pix_width, pix_height)

	arr := InitializeArray(img, sample_size, pix_height, pix_width)

	LogOut(fmt.Sprintf("LOGGING >> Took %s to make and propagate pixel info", time.Since(intermediate)))
	intermediate = time.Now()

	// TRANSFORMATIONS**********************************************
	// GetRunes(arr)
	// arr = transforms.XDoG(arr)
	// arr = transforms.DoG(arr)
	// transforms.AsciiFilter(arr)
	transforms.NaiveAsciiFilter(arr)
	// arr = transforms.SobelFilter(arr, false)

	// for i := range len(arr) {
	// 	for j := range len(arr[i]) {
	// 		fmt.Printf("%v ", arr[i][j])
	// 	}
	// 	fmt.Println()
	// }
	// fmt.Println("DONE")

	// TRANSFORMATIONS**********************************************

	LogOut(fmt.Sprintf("LOGGING >> CHARACTER TRANSFORMATIONS DONE: %s", time.Since(intermediate)))
	intermediate = time.Now()

	// BOUNDS FOR KEEPING THE IMAGE QUALITY PERFECT:
	out_width := pix_width * px_size
	out_height := pix_height * px_size

	newimg := image.NewRGBA(image.Rect(0, 0, out_width, out_height))
	draw.Draw(newimg, newimg.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)

	context := InitializeContext(newimg, float64(px_size))

	buffer := transforms.InitializeBuffer(0, px_size, out_width, out_height, px_size, newimg)

	LogOut(fmt.Sprintf("LOGGING >> Did pre-processing for image drawing (blank image, created image buffer, parsed font): %s", time.Since(intermediate)))
	intermediate = time.Now()

	buffer.WriteArray(context, arr)

	LogOut(fmt.Sprintf("LOGGING >> Took %s to draw pixels in buffer", time.Since(intermediate)))
	intermediate = time.Now()

	op_end := time.Since(operations)

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

	LogOut(fmt.Sprintf("LOGGING >> Time to encode output image: %s", time.Since(intermediate)))

	LogOut(fmt.Sprintf("LOGGING >> Total execution time: %s", time.Since(start)))
	LogOut(fmt.Sprintf("LOGGING >> Took %s for non-(encode/decode) operations", op_end))
	fmt.Println("Created new image")
	intermediate = time.Now()
}

// *****************
// LOGGING FXNS
// *****************
func LogOut(message string) {
	size := len(message)

	var sb strings.Builder
	for i := 0; i < size+4; i++ {
		sb.WriteRune('*')
	}
	fmt.Println(sb.String())
	fmt.Println("* " + message + " *")
	fmt.Println(sb.String() + "\n")
}

func WriteToTXT(height int, width int, mapping map[int]rune, arr [][]transforms.Pixel) {
	// .txt output
	var sb strings.Builder
	for i := range height {
		for j := range width {
			cur := arr[i][j]
			sb.WriteRune(cur.Character)
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
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

	fmt.Println("Opened " + img.Name() + "\n")

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

	fmt.Println("Opened " + img.Name() + "\n")

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
	c.SetClip(newimg.Bounds())
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
