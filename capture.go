package bigbro

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	CapturePath = "captures"
)

type Capture struct {
	Actor string    `json:"actor"`
	Data  string    `json:"data"`
	Time  time.Time `json:"time"`
}

func WriteCapture(capture Capture) error {
	fpath := path.Join(CapturePath, capture.Actor)
	fname := strconv.Itoa(int(capture.Time.Unix())) + ".png"
	err := os.MkdirAll(fpath, os.ModePerm)
	if err != nil {
		return err
	}

	unbased, err := base64.StdEncoding.DecodeString(capture.Data[strings.IndexByte(capture.Data, ',')+1:])
	if err != nil {
		return err
	}

	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(fpath, fname))
	if err != nil {
		return err
	}
	err = png.Encode(f, im)
	return err
}
