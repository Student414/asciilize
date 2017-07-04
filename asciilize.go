package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"math"
	"os"
)

// max block size = 4096 X 4096

// gray to ascii map
var grayMap = [...]string{"@", "w", "#", "$", "k", "d", "t", "j", "i", ".", " "}

func main() {
	ptr_filename := flag.String("f", "./a.jpg", "Define the source jpeg file")
	ptr_blockX := flag.Int("bx", 5, "Define the block size in X axis")
	ptr_blockY := flag.Int("by", 10, "Define the block size in Y axis")
	ptr_step := flag.Int("j", 25, "sampling step")
	flag.Parse()

	_, asciiData, err := Asciilize(*ptr_filename, *ptr_blockX, *ptr_blockY, *ptr_step)
	if nil != err || *ptr_blockX < 1 || *ptr_blockY < 1 {
		fmt.Println(err)
		flag.Usage()
		fmt.Println("Example: ./program -f ./hello.jpg -bx 4 -by 6")
		fmt.Println("可以使用-j n参数来控制采用的细腻度，表示每n个点取样一次，如果n等于0，则关闭此功能,n=0等价于n=1")
	}
	OutputAsciilizedData(asciiData)

}

func Asciilize(jpegFile string, blockSizeX int, blockSizeY int, jump int) (grayData [][]uint8, asciiData [][]string, err error) {
	blockX := blockSizeX
	blockY := blockSizeY
	// open file
	m_file, err := os.Open(jpegFile)
	defer m_file.Close()
	if err != nil {
		// fmt.Print(err)
		return nil, nil, err
	}

	// decode jpeg
	m_img, err := jpeg.Decode(m_file)
	if err != nil {
		// fmt.Print("===")
		// fmt.Print(err)
		return nil, nil, err
	}

	// pre
	m_bounds := m_img.Bounds()

	// split blocks
	nBlockX := (int)(math.Ceil((float64)(m_bounds.Dx()) / float64(blockX)))
	nBlockY := (int)(math.Ceil((float64)(m_bounds.Dy()) / float64(blockY)))

	// gen 2d array
	result := make([][]uint8, nBlockY)
	for row := 0; row < nBlockY; row++ {
		subresult := make([]uint8, nBlockX, nBlockX)
		for col := 0; col < nBlockX; col++ {
			subresult[col] = 0
		}
		result[row] = subresult
	}

	// calc grey value
	var pickColorCount uint = 0
	for x := 0; x < nBlockX; x++ {
		for y := 0; y < nBlockY; y++ {
			// fmt.Print("1")
			// every block
			var graySum uint32 = 0
			pixelX := x * blockX
			pixelY := y * blockY
			breakPixel := 0
			stepCount := 0
			for i := 0; i < blockX; i++ {
				for j := 0; j < blockY; j++ {
					if 0 != jump {
						stepCount %= jump
						if 0 != stepCount {
							breakPixel++
							stepCount++
							continue
						} else {
							stepCount++
						}
					}
					targetX := pixelX + i
					targetY := pixelY + j
					if targetX >= m_bounds.Dx() || targetY >= m_bounds.Dy() {
						breakPixel++
						continue
					}
					// m_color := m_img.At(targetX, targetY)
					// gray := image.NewGray(image.Rect(0, 0, m_bounds.Dx(), m_bounds.Dy())).GrayAt(targetX, targetY)
					pickColorCount++
					r, g, b, a := m_img.At(targetX, targetY).RGBA()
					r /= 256
					g /= 256
					b /= 256
					a /= 256
					graySum += (uint32)((float64(r)*float64(0.3) + float64(g)*float64(0.59) + float64(b)*float64(0.11)) * (float64(a) / float64(256.0)))
				}
			}
			graySum /= uint32(blockX*blockY - breakPixel)
			// fmt.Println(graySum)
			result[y][x] = uint8(graySum)
		}
	}
	// output pick time
	fmt.Printf("Pick Color %d times\n", pickColorCount)

	// to ascii 2d array
	step := float64(256) / float64(len(grayMap))
	// var ascii [][]string
	ascii := make([][]string, nBlockY)
	for row := 0; row < nBlockY; row++ {
		subascii := make([]string, nBlockX)
		for col := 0; col < nBlockX; col++ {
			// fmt.Print(int(result[row][col]) / step)
			index := int(float64(result[row][col]) / step)
			subascii[col] = grayMap[index]
		}
		ascii[row] = subascii
	}
	grayData = result
	asciiData = ascii
	return
}

func OutputAsciilizedData(asciiData [][]string) {
	// format output
	nBlockY := len(asciiData)
	if 0 == nBlockY {
		return
	}
	nBlockX := len(asciiData[0])
	for row := 0; row < nBlockY; row++ {
		for col := 0; col < nBlockX; col++ {
			fmt.Print(asciiData[row][col])
		}
		fmt.Println()
	}
}

/*




fmt.Print("Test Output !\n")
	// load jpeg file
	m_file, err := os.Open("photo.jpeg")
	defer m_file.Close()
	if err != nil {
		fmt.Print(err)
	}
	// decode jpeg
	m_img, err := jpeg.Decode(m_file)
	if err != nil {
		fmt.Print("===")
		fmt.Print(err)
	}
	// fmt.Print("info" + info)
	// output some information
	// clr := image.At(100, 100)
	m_bounds := m_img.Bounds()
	// fmt.Println(color)
	fmt.Println(m_bounds.Dx(), " ", m_bounds.Dy())
	// split blocks
	nBlockX := (int)(math.Ceil((float64)(m_bounds.Dx()) / blockX))
	nBlockY := (int)(math.Ceil((float64)(m_bounds.Dy()) / blockY))
	fmt.Printf("horizontal blocks : %d\nvertical blocks : %d\n", nBlockX, nBlockY)

	// gray to ascii map
	var grayMap = [...]string{"@", "w", "#", "$", "k", "d", "t", "j", "i", ".", " "}
	// gen 2d array
	var result [][]uint8
	for row := 0; row < nBlockY; row++ {
		subresult := make([]uint8, 0, nBlockX)
		for col := 0; col < nBlockX; col++ {
			subresult = append(subresult, 0)
		}
		result = append(result, subresult)
	}

	// calc grey value
	for x := 0; x < nBlockX; x++ {
		for y := 0; y < nBlockY; y++ {
			// fmt.Print("1")
			// every block
			var graySum uint32 = 0
			pixelX := x * blockX
			pixelY := y * blockY
			var breakPixel uint32
			breakPixel = 0
			for i := 0; i < blockX; i++ {
				for j := 0; j < blockY; j++ {
					targetX := pixelX + i
					targetY := pixelY + j
					if targetX >= m_bounds.Dx() || targetY >= m_bounds.Dy() {
						breakPixel++
						continue
					}
					// m_color := m_img.At(targetX, targetY)
					// gray := image.NewGray(image.Rect(0, 0, m_bounds.Dx(), m_bounds.Dy())).GrayAt(targetX, targetY)
					r, g, b, a := m_img.At(targetX, targetY).RGBA()
					r /= 256
					g /= 256
					b /= 256
					a /= 256
					graySum += (uint32)((float64(r)*float64(0.3) + float64(g)*float64(0.59) + float64(b)*float64(0.11)) * (float64(a) / float64(256.0)))
				}
			}
			graySum /= blockX*blockY - breakPixel
			// fmt.Println(graySum)
			result[y][x] = uint8(graySum)
		}
	}
	// fmt.Println(result)
	// to ascii 2d array
	step := 256 / len(grayMap)
	var ascii [][]string
	for row := 0; row < nBlockY; row++ {
		subascii := make([]string, 0, nBlockX)
		for col := 0; col < nBlockX; col++ {
			subascii = append(subascii, grayMap[int(result[row][col])/step-1])
		}
		ascii = append(ascii, subascii)
	}
	// fmt.Println(ascii)



*/
