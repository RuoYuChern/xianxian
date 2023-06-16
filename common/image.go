package common

import (
	"bytes"
	"image"
	"image/png"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
)

func ImgResize(data []byte, h int, w int) ([]byte, error) {
	img, name, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		Logger.Infof("Decode failed:%s", err.Error())
		return nil, err
	}
	Logger.Infof("Decode format:%s", name)
	inverted := effect.Invert(img)
	resized := transform.Resize(inverted, w, h, transform.Linear)
	buff := new(bytes.Buffer)
	err = png.Encode(buff, resized)
	if err != nil {
		Logger.Infof("Decode failed:%s", err.Error())
		return nil, err
	}
	return buff.Bytes(), nil
}
