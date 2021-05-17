package main

import (
	"art/sketch"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/fogleman/gg"
)

var (
	srcImageName = "srcImage3.jpg"
	outImgName   = "outImage3.png"
	cycleCount   = 5000
)

func main() {
	rand.Seed(time.Now().Unix())

	srcImg, err := gg.LoadImage(srcImageName)
	if err != nil {
		log.Panicln(err)
	}

	destWidth := 2000
	sketch := sketch.NewSketch(srcImg, sketch.UserParams{

		StrokeRatio:              0.75,
		DestWidth:                destWidth,
		DestHeight:               2000,
		InitialAlpha:             0.1,
		StrokeReduction:          0.002,
		AlphaIncrease:            0.06,
		StrokeInversionThreshold: 0.05,
		StrokeJitter:             int(0.1 * float64(destWidth)),
		MinEdgeCount:             3,
		MaxEdgeCount:             4,
	})
	for i := 0; i < cycleCount; i++ {
		sketch.Update()
	}

	saveOutput(sketch.Output(), outImgName)

}

func saveOutput(img image.Image, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode to `PNG` with `DefaultCompression` level
	// then save to file
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}
