package main

import (
	"encoding/base64"
	"image/color"
	"image"
	"bytes"
	_ "image/png"
	"image/png"
)

const width = 1000
const height = 500

type ImageInterface interface {
	ColorModel() color.Model
	Bounds() image.Rectangle
	At(x, y int) color.Color
}

type Image struct {
	w, h int
	algo algoFunc
}

type algoFunc func(img Image, x, y int) color.Color

func (img Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (img Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.w, img.h)
}

func (img Image) At(x, y int) color.Color {
	if img.algo == nil {
		return img.DefaultAlgo(x, y)
	} else {
		return img.algo(img, x, y)
	}
}

func (img Image) DefaultAlgo(x, y int) color.Color {
	return color.Gray16{
		Y: 0,
	}
}

func GenerateImage(algo algoFunc) string {
	m := Image{
		w:    width,
		h:    height,
		algo: algo,
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, m)
	if err != nil {
		panic(err)
	}
	enc := base64.StdEncoding.EncodeToString(buf.Bytes())
	return enc
}
