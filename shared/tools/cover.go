package tools

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image/jpeg"
	"os"
)

func GetVideoCover(url string) ([]byte, error) {
	reader := bytes.NewBuffer(nil)
	err := ffmpeg.Input(url).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(reader, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}

	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
