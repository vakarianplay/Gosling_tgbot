package main

import (
	"fmt"
	// "io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

var textPath string
var picDir string
var fontFile string
var fontSize int

func picEntry(textPath_, picDir_, fontFile_, fontSize_ string) {
	textPath = textPath_
	picDir = picDir_
	fontFile = fontFile_
	fontSize, _ = strconv.Atoi(fontSize_)
	// outStr()
	fmt.Println("Pic settings init   font: "+fontFile, "size: ", fontSize)
}

func generatePic(txt []string) {

	imgFile, err := selectRandomFile()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(imgFile)

	im, err := gg.LoadImage(imgFile)
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
	if err := dc.LoadFontFace(fontFile, float64(fontSize)); err != nil {
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

	dc.SavePNG(strconv.FormatInt(time.Now().UnixMilli(), 10) + "_out.png")
	// fmt.Println(textWidth)
	// fmt.Println(width, " ", height)

}

func outStr() {

	content, err := os.ReadFile(textPath)
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

func getLineTst() string {

	content, err := os.ReadFile(textPath)
	if err != nil {
		log.Fatal(err)
	}
	str := string(content)
	lines := strings.Split(str, "\n")
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(lines))
	randomLine := lines[randomIndex]

	return randomLine
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

func selectRandomFile() (string, error) {

	files, err := os.ReadDir(picDir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("No files")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(files))
	randomFile := files[randomIndex]

	return filepath.Join(picDir, randomFile.Name()), nil
}
