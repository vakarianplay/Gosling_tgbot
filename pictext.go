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
	// generatePic("aaaaaaaииииии")
	outStr()
}

func generatePic(txt string) {

	// txt := "Хто я?"
	im, err := gg.LoadImage("pic/nav.jpg")
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
	if err := dc.LoadFontFace("font/Roboto-Bold.ttf", 36); err != nil {
		panic(err)
	}

	textWidth, textHeight := dc.MeasureString(txt)

	dc.DrawImage(im, 0, 0)

	dc.SetRGB255(0, 0, 0)
	dc.DrawStringAnchored(txt, float64(width)/2+3, (float64(height)/2+3)+textHeight+150, 0.5, 0.5)

	dc.SetRGB255(255, 255, 255)
	dc.DrawStringAnchored(txt, float64(width)/2, (float64(height)/2)+textHeight+150, 0.5, 0.5)

	dc.Clip()
	dc.SavePNG("outnav.png")

	fmt.Println(textWidth)
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

	generatePic(randomLine)
}
