package main

import (
	"fmt"
	// "io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

func main() {

	outStr()
}

func generatePic(txt []string) {

	im, err := gg.LoadImage("pic/2.jpg")
	if err != nil {
		log.Fatal(err)
	}
	// dc := gg.NewContext(0, 0)
	// width, height := dc.Image().Bounds().Max.X, dc.Image().Bounds().Max.Y
	width, height := im.Bounds().Max.X, im.Bounds().Max.Y
	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("font/Roboto-Bold.ttf", 30); err != nil {
		panic(err)
	}

	dc.DrawImage(im, 0, 0)

	for i := 0; i < len(txt); i++ {
		_, textHeight := dc.MeasureString(txt[i])

		dc.SetRGB255(0, 0, 0)
		dc.DrawStringAnchored(txt[i], float64(width)/2+3, (float64(height)/2+3)+textHeight+100+float64(i*30), 0.5, 0.5)

		dc.SetRGB255(255, 255, 255)
		dc.DrawStringAnchored(txt[i], float64(width)/2, (float64(height)/2)+textHeight+100+float64(i*30), 0.5, 0.5)
		// dc.Clip()

	}
	dc.Clip()
	// outName := time.Now().Weekday().String()
	dc.SavePNG("out.png")
	// fmt.Println(textWidth)
	fmt.Println(width, " ", height)

}

func outStr() {

	content, err := os.ReadFile("txt/tagged.txt")
	if err != nil {
		log.Fatal(err)
	}
	str := string(content)
	lines := strings.Split(str, "\n")
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(lines))
	randomLine := lines[randomIndex]
	fmt.Println(randomLine)

	// lines = splLine(randomLine)
	// for _, line := range lines {
	// 	fmt.Println(line)
	// }
	// fmt.Println(splLine(randomLine))
	generatePic(splLine(randomLine))
}

func splLine(randomLine string) []string {

	words := strings.Fields(randomLine)
	var result []string
	currentLine := ""

	for _, word := range words {
		currentLine += word + " "

		if len(strings.Fields(currentLine)) >= 4 {
			result = append(result, strings.TrimSpace(currentLine))
			currentLine = ""
		}
	}

	if currentLine != "" {
		result = append(result, strings.TrimSpace(currentLine))
	}

	return result
}
