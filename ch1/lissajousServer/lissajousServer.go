// lissajousServer.go
// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

//!-main
// Packages not needed by version in book.

//!+main

var palette = []color.Color{
	color.RGBA{0x00, 0xff, 0x00, 0xff}, // light blue
	color.Black,
	color.RGBA{0xff, 0x00, 0x00, 0xff}, // red
	color.RGBA{0x80, 0xc0, 0xff, 0xff}, // light blue
	color.RGBA{0xff, 0xff, 0x00, 0xff}, // yellow
}

const (
	whiteIndex  = 0 // first color in palette
	blackIndex  = 1 // next color in palette
	redIndex    = 2 // next color in palette
	blueIndex   = 3 // next color in palette
	yellowIndex = 4 // next color in palette (not used in this example)
)

func main() {
	//!-main
	// Use a local random source and generator (recommended way)
	rng := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lissajous(w, rng)
	})

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Update lissajous to accept *rand.Rand
func lissajous(out io.Writer, rng *rand.Rand) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rng.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				uint8(i%len(palette))) // Use i to cycle through palette colors
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

//!-main
//!-main
