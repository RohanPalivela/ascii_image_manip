package transforms

func NaiveAsciiFilter(arr [][]Pixel) {
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

	LuminFilter(arr, mapping)

	sobel := SobelFilter(arr, true)

	for i := range len(sobel) {
		for j := range len(sobel[i]) {
			char := sobel[i][j].Character
			if char != ' ' {
				arr[i][j].Character = char
			}
		}
	}
}

func AsciiFilter(arr [][]Pixel) {
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

	LuminFilter(arr, mapping)

	edged := DoG(arr, 1, 15)

	sobel := SobelFilter(edged, true)

	for i := range len(sobel) {
		for j := range len(sobel[i]) {
			char := sobel[i][j].Character
			if char != ' ' {
				arr[i][j].Character = char
			}
		}
	}
}
