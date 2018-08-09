package main

import (
	"github.com/alexflint/go-arg"
	"github.com/dustin/go-heatmap"
	"os"
	"github.com/hscells/bigbro"
	"image"
	"github.com/dustin/go-heatmap/schemes"
	"image/png"
	"encoding/csv"
	"io"
	"time"
	"strconv"
	"image/draw"
	"regexp"
	"fmt"
	"image/gif"
	"image/color/palette"
	"github.com/andybons/gogif"
	"image/color"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type args struct {
	Log      string `arg:"help:path to log file produced by bigbro,required"`
	Heatmap  string `arg:"help:path to output heatmap to,required"`
	Image    string `arg:"help:path to image file for analysis"`
	Start    string `arg:"help:time to start analysis"`
	End      string `arg:"help:time to end analysis"`
	Interval int64  `arg:"help:create an animation every n seconds"`
	Location string `arg:"help:URL to limit analysis (regex)"`
	Actor    string `arg:"help:actor to limit analysis"`
	Method   string `arg:"help:method to limit analysis"`
	Width    int    `arg:"help:width to limit analysis"`
	Height   int    `arg:"help:height to limit analysis"`
}

func (args) Version() string {
	return "bbheat 09.Aug.2018"
}

func (args) Descriptions() string {
	return "heatmap generator for bigbro data"
}

func main() {
	var args args
	arg.MustParse(&args)

	location, err := regexp.Compile(args.Location)
	if err != nil {
		panic(err)
	}

	// Get the start and end times (if any).
	var start, end time.Time
	if len(args.Start) > 0 {
		start, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", args.Start)
		if err != nil {
			panic(err)
		}
	}
	if len(args.End) > 0 {
		end, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", args.End)
		if err != nil {
			panic(err)
		}
	}

	// Open the log file.
	logFile, err := os.OpenFile(args.Log, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Open the image file to output.
	imageFile, err := os.OpenFile(args.Heatmap, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	var events []bigbro.Event
	r := csv.NewReader(logFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		// Parse fields of the log into Go types.
		t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", record[0])
		if err != nil {
			panic(err)
		}
		x, err := strconv.ParseInt(record[7], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(record[8], 10, 64)
		if err != nil {
			panic(err)
		}
		sw, err := strconv.ParseInt(record[9], 10, 64)
		if err != nil {
			panic(err)
		}
		sh, err := strconv.ParseInt(record[10], 10, 64)
		if err != nil {
			panic(err)
		}

		// Create an event from the log record.
		event := bigbro.Event{
			Time:         t,
			Actor:        record[1],
			Method:       record[2],
			Target:       record[3],
			Name:         record[4],
			ID:           record[5],
			Location:     record[6],
			X:            int(x),
			Y:            int(y),
			ScreenWidth:  int(sw),
			ScreenHeight: int(sh),
			Comment:      record[11],
		}

		// Filter heatmap based on command arguments.
		if len(args.Method) > 0 && event.Method != args.Method {
			continue
		}
		if len(args.Actor) > 0 && event.Actor != args.Actor {
			continue
		}
		if len(args.Location) > 0 && !location.MatchString(event.Location) {
			continue
		}
		if args.Width > 0 && event.ScreenWidth != args.Width {
			continue
		}
		if args.Height > 0 && event.ScreenHeight != args.Height {
			continue
		}
		if len(args.Start) > 0 && event.Time.Before(start) {
			continue
		}
		if len(args.End) > 0 && event.Time.After(end) {
			break // We can stop processing early here actually.
		}

		events = append(events, event)
	}

	// Make the heatmap points.
	w, h, points := makePoints(events)

	if len(events) == 0 {
		fmt.Println("no events to process")
		return
	}

	// If no interval is specified, output a static image.
	if args.Interval == 0 {
		img := drawHeatmap(args.Image, w, h, points)
		// Write the heatmap to file.
		err = png.Encode(imageFile, img)
		if err != nil {
			panic(err)
		}
	} else { // Otherwise, output an animation as a GIF.
		// Create the interval and first amount of time.
		interval := time.Duration(args.Interval) * time.Second
		t := events[0].Time.Add(interval)
		var intervalEvents []bigbro.Event

		col := color.RGBA{R: 200, G: 100, A: 255}
		point := fixed.Point26_6{X: fixed.Int26_6(12 * 64), Y: fixed.Int26_6(12 * 64)}

		var images []*image.Paletted
		var delays []int
		for i, event := range events {
			intervalEvents = append(intervalEvents, event)
			if event.Time.After(t) {
				fmt.Println(event.Time)
				t = event.Time.Add(interval)
				w, h, points := makePoints(intervalEvents)
				intervalEvents = make([]bigbro.Event, 0)
				img := drawHeatmap(args.Image, w, h, points)

				// https://stackoverflow.com/questions/38299930/how-to-add-a-simple-text-label-to-an-image-in-go
				d := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(col),
					Face: basicfont.Face7x13,
					Dot:  point,
				}
				d.DrawString(fmt.Sprintf("[%d/%d] %s - %s", i, len(events), t, event.Time.String()))

				// https://gist.github.com/rbwendt/29b46678600e019800154652d5f5b054
				pimg := image.NewPaletted(img.Bounds(), palette.Plan9)

				// https://stackoverflow.com/questions/35850753/how-to-convert-image-rgba-image-image-to-image-paletted
				quantizer := gogif.MedianCutQuantizer{NumColor: 64}
				quantizer.Quantize(pimg, img.Bounds(), img, image.ZP)
				images = append(images, pimg)
				delays = append(delays, 100)
			}
		}
		err := gif.EncodeAll(imageFile, &gif.GIF{
			Image: images,
			Delay: delays,
		})
		if err != nil {
			panic(err)
		}
	}

	return
}

// makePoints creates a series of points out of some events.
func makePoints(events []bigbro.Event) (int, int, []heatmap.DataPoint) {
	var w, h int
	var points []heatmap.DataPoint
	for _, e := range events {
		points = append(points, heatmap.P(float64(e.X), float64(e.ScreenHeight-e.Y)))
		if w > 0 && e.ScreenWidth != w {
			panic(fmt.Errorf("record contains differing widths: found %d (expected %d)", e.ScreenWidth, w))
		}
		if h > 0 && e.ScreenHeight != h {
			panic(fmt.Errorf("record contains differing heights: found %d (expected %d)", e.ScreenHeight, h))
		}
		w = e.ScreenWidth
		h = e.ScreenHeight
	}
	return w, h, points
}

// drawHeatmap creates a heatmap image from a series of points.
// If the length of fname is 0, the image will only contain the heatmap.
func drawHeatmap(fname string, w, h int, p []heatmap.DataPoint) *image.RGBA {
	// Make the heatmap.
	imgHeatmap := heatmap.Heatmap(image.Rect(0, 0, w, h), p, 32, 128, schemes.AlphaFire)
	bounds := imgHeatmap.Bounds()
	img := image.NewRGBA(bounds)
	if len(fname) > 0 {
		websiteFile, err := os.OpenFile(fname, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer websiteFile.Close()
		imgWebsite, err := png.Decode(websiteFile)
		if err != nil {
			panic(err)
		}
		draw.Draw(img, bounds, imgWebsite, image.ZP, draw.Src)
		draw.Draw(img, bounds, imgHeatmap, image.ZP, draw.Over)
	} else {
		draw.Draw(img, imgHeatmap.Bounds(), imgHeatmap, image.ZP, draw.Src)
	}
	return img
}
