package main

import (
	"image/color"
	"github.com/ojrac/opensimplex-go"
	"time"
)

const MAX_COLOR_VALUE = 65535

type algoFuncBasic func(w, h, x, y int) uint16

func ConvertToGray16AlgoFunc(algo algoFuncBasic) algoFunc {
	return func(img Image, x, y int) color.Color {
		return color.Gray16{
			Y: algo(img.w, img.h, x, y),
		}
	}
}

var VericalGradient algoFuncBasic = func(w, h, x, y int) uint16 {
	return uint16(float64(y) / float64(h-1) * MAX_COLOR_VALUE)
}

var HorizontalGradient algoFuncBasic = func(w, h, x, y int) uint16 {
	return uint16(float64(x) / float64(w-1) * MAX_COLOR_VALUE)
}

var CornerGradient algoFuncBasic = func(w, h, x, y int) uint16 {
	return uint16(((float64(x) + float64(y)) / float64(w+h)) * MAX_COLOR_VALUE)
}

var CryptoRandom algoFuncBasic = func(w, h, x, y int) uint16 {
	return uint16(
		GenerateRandomUint64(MAX_COLOR_VALUE) % MAX_COLOR_VALUE,
	)
}

// Pass value from 0.0 to 1.0
var CryptoRandomThreshold = func(threshold float64) algoFuncBasic {
	calculatedThreshold := uint16(MAX_COLOR_VALUE * threshold)
	return func(w, h, x, y int) uint16 {
		random := CryptoRandom(w, h, x, y)
		if random > calculatedThreshold {
			return MAX_COLOR_VALUE
		} else {
			return 0
		}
	}
}

var SimplexNoise = func() algoFuncBasic {
	var noise = opensimplex.NewWithSeed(time.Now().UnixNano())
	return func(w, h, x, y int) uint16 {
		return uint16((noise.Eval2(float64(x), float64(y)) + 1) / 2 * MAX_COLOR_VALUE)
	}
}
