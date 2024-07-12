package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"net/http"
	"os"
)

var asciiChars = []string{
	"$", "@", "B", "%", "8", "&", "W", "M", "#", "*", "o", "a", "h", "k", "b", "d",
	"p", "q", "w", "m", "Z", "O", "0", "Q", "L", "C", "J", "U", "Y", "X", "z", "c",
	"v", "u", "n", "x", "r", "j", "f", "t", "/", "\\", "|", "(", ")", "1", "{", "}",
	"[", "]", "?", "-", "_", "+", "~", "i", "!", "l", "I", ";", ":", ",", "\"", "^", "`", ".",
}

func getImageFromFile(filePath string) image.Image {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Faled to load the image: %v", err)
	}
	defer file.Close()
	img, fileName, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Faled to decode %s : %v", fileName, err)
	}
	return img
}

func grayScale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			grayImg.Set(x, y, grayColor)
		}
	}

	// A part of code to test if gray scale works
	// logImg, err := os.Create("log_img.jpg")
	// if err != nil {
	// 	fmt.Printf("Can't create log image file, error: %v \n", err)
	// } else {
	// 	defer logImg.Close()
	// 	if err := jpeg.Encode(logImg, grayImg, nil); err != nil {
	// 		fmt.Printf("Can't save log image, error: %v \n", err)
	// 	}
	// }

	return grayImg
}

func grayToASCII(gray *image.Gray, scale float64) [][]string {
	bounds := gray.Bounds()
	asciiWidth := int(float64(bounds.Dx()) * scale)
	asciiHeight := int(float64(bounds.Dy()) * scale)

	// Calculate aspect ratio to maintain proportions
	aspectRatio := float64(bounds.Dx()) / float64(bounds.Dy())
	if aspectRatio < 1 {
		asciiWidth = int(float64(asciiHeight) * aspectRatio)
	} else {
		asciiHeight = int(float64(asciiWidth) / aspectRatio)
	}

	ascii := make([][]string, asciiHeight)

	for y := 0; y < asciiHeight; y++ {
		ascii[y] = make([]string, asciiWidth)
		for x := 0; x < asciiWidth; x++ {
			origX := int(float64(x) / scale)
			origY := int(float64(y) / scale)

			grayColor := gray.GrayAt(origX, origY).Y

			index := int(grayColor) * len(asciiChars) / 256
			ascii[y][x] = asciiChars[index]
		}
	}

	return ascii
}

func createArtString(asciiArt [][]string) string {
	asciiStr := ""
	for _, row := range asciiArt {
		for _, char := range row {
			asciiStr += string(char) + " "
		}
		asciiStr += "\n"
	}
	return asciiStr
}

func handleASCIIArt(w http.ResponseWriter, r *http.Request) {
	img := getImageFromFile("/home/user/Pictures/ascii_test_image.jpg")
	grayImg := grayScale(img)

	asciiArtScale := 0.5

	asciiArt := grayToASCII(grayImg, asciiArtScale)
	asciiArtString := createArtString(asciiArt)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(asciiArtString))
}

func main() {
	http.HandleFunc("/ascii-art", handleASCIIArt)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
