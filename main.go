package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"github.com/nfnt/resize"
)

func main() {
	// open "test.jpg"
	if _, err := os.Stat("./"); os.IsNotExist(err) {
		err = os.Mkdir("input", 0777)
		if err != nil {
			log.Println(err)
		}
	}
	if _, err := os.Stat("./"); os.IsNotExist(err) {
		err = os.Mkdir("output", 0777)
		if err != nil {
			log.Println(err)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("- OPTIMAZE - Image optimizer by Sharif")
	fmt.Println("Put your JPG image/s in 'input' folder")
	fmt.Println("width,height,quality,")
	text, _ := reader.ReadString('\n')

	p := strings.Split(text, ",")

	log.Printf("p = %#v", p)

	w64, _ := strconv.ParseUint(string(p[0]), 10, 64)
	h64, _ := strconv.ParseUint(string(p[1]), 10, 64)
	q64, _ := strconv.ParseInt(string(p[2]), 10, 64)
	op := p[3]

	width := uint(w64)
	height := uint(h64)
	quality := int(q64)

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Println(err)
	}

	for _, f := range files {
		//fmt.Println(f.Name())
		input := "./"
		input += f.Name()
		if strings.Contains(input, ".jpg") {
			fmt.Println(input)

			imwInt, imhInt, size := getFileInfo(input)
			imw := uint(imwInt)
			imh := uint(imhInt)

			switch op {
			case "h":
				width = imw / 2
				height = imh / 2
			case "q":
				width = imw / 4
				height = imh / 4
			default:
				width = imw
				height = imh
			}

			file, err := os.Open(input)
			if err != nil {
				log.Println(err)
			}

			// decode jpeg into image.Image
			img, err := jpeg.Decode(file)
			if err != nil {
				log.Println(err)
				file.Close()
			} else {
				file.Close()

				log.Println("Input:", input,"Width:", width, "Height:", height, "Size:", ByteCountSI(size))

				var opt jpeg.Options

				opt.Quality = quality

				// resize to width 1000 using Lanczos resampling
				// and preserve aspect ratio
				m := resize.Resize(width, height, img, resize.Lanczos3)

				output := "./"
				output += f.Name()

				out, err := os.Create(output)
				if err != nil {
					log.Println(err)
				}
				defer out.Close()
				// write new image to file
				jpeg.Encode(out, m, &opt)
				var outsize int64 = 0
				outstat, outstatErr := out.Stat()
				if outstatErr == nil {
					outsize = outstat.Size()
				}
				log.Println("Output:", output, "Width:", width, "Height:", height, "Size:", ByteCountSI(outsize))
			}
		} else {
			log.Printf("Filename is: %s. Skipping", input)
		}
	}

}

func getFileInfo(imgPath string) (int, int, int64) {
	file, err := os.Open(imgPath)
	if err != nil {
		log.Println(err)
	}

	stats, serr := file.Stat()
	var size int64 = 0
	if serr == nil {
		size = stats.Size()
	}

	im, _, err := image.DecodeConfig(file) // Image Struct
	if err != nil {
		log.Println(err)
	}
	return im.Width, im.Height, size
}

// See https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

// See https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}