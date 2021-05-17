package sketch

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
)

type UserParams struct {
	DestWidth                int
	DestHeight               int
	StrokeJitter             int
	MinEdgeCount             int
	MaxEdgeCount             int
	StrokeRatio              float64
	StrokeReduction          float64
	StrokeInversionThreshold float64
	InitialAlpha             float64
	AlphaIncrease            float64
}

type Sketch struct {
	params            UserParams
	source            image.Image
	dc                *gg.Context // drawing context
	sourceWidth       int
	sourceHeight      int
	strokeSize        float64
	initialStrokeSize float64
}

func NewSketch(source image.Image, userParams UserParams) *Sketch {
	s := &Sketch{params: userParams}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.initialStrokeSize = s.params.StrokeRatio * float64(s.params.DestWidth)
	s.strokeSize = s.initialStrokeSize

	canvas := gg.NewContext(s.params.DestWidth, s.params.DestHeight)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.params.DestWidth), float64(s.params.DestHeight))
	canvas.FillPreserve()

	s.source = source
	s.dc = canvas
	return s
}

func (s *Sketch) Update() {
	// core drawing logic
	// 1. Obtain color information from the source
	rndX := rand.Float64() * float64(s.sourceWidth)
	rndY := rand.Float64() * float64(s.sourceHeight)
	r, g, b := rgb255(s.source.At(int(rndX), int(rndY)))

	// 2. Determine a destination in the output space
	destX := rndX * float64(s.params.DestWidth) / float64(s.sourceWidth)
	destX += float64(randRange(s.params.StrokeJitter))
	destY := rndY * float64(s.params.DestHeight) / float64(s.sourceHeight)
	destY += float64(randRange(s.params.StrokeJitter))

	// 3. Draw a "stroke" using the desired parameters
	edges := s.params.MinEdgeCount + rand.Intn(s.params.MaxEdgeCount-s.params.MinEdgeCount+1)

	s.dc.SetRGBA255(r, g, b, int(s.params.InitialAlpha))
	s.dc.DrawRegularPolygon(edges, destX, destY, s.strokeSize, rand.Float64())
	s.dc.FillPreserve()

	if s.strokeSize <= s.params.StrokeInversionThreshold*s.initialStrokeSize {
		if (r+g+b)/3 < 128 {
			s.dc.SetRGBA255(255, 255, 255, int(s.params.InitialAlpha*2))
		} else {
			s.dc.SetRGBA255(0, 0, 0, int(s.params.InitialAlpha*2))
		}
	}
	s.dc.Stroke()

	// 4. Update the parameter state for the next execution
	s.strokeSize -= s.params.StrokeReduction * s.strokeSize
	s.params.InitialAlpha += s.params.AlphaIncrease
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 257), int(g0 / 257), int(b0 / 257)
}

func randRange(max int) int {
	return -max + rand.Intn(2*max)
}

func (s *Sketch) Output() image.Image {
	return s.dc.Image()
}
