package main

import (
	"image/color"
	"github.com/ojrac/opensimplex-go"
	"math"
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

var VerticalGradient algoFuncBasic = func(w, h, x, y int) uint16 {
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

var SimplexNoise = func(frequency, threshold float64, seed int64) algoFuncBasic {
	var noise = opensimplex.NewWithSeed(seed)
	return func(w, h, x, y int) uint16 {
		var random = (noise.Eval2(frequency*float64(x), frequency*float64(y)) + 1) / 2

		if threshold < 0 {
			return uint16(random * MAX_COLOR_VALUE) % MAX_COLOR_VALUE
		} else {
			if random > threshold {
				return MAX_COLOR_VALUE
			} else {
				return 0
			}
		}
	}
}

var SimplexNoiseOctaves = func(frequency, threshold float64, seed int64, octaves int) algoFuncBasic {
	var simplex = SimplexNoise(frequency, threshold, seed)
	return func(w, h, x, y int) uint16 {
		var res float64
		var step float64 = 1.0
		for i := octaves; i > 0; i-- {
			res += step * float64(simplex(w, h, int(1/step)*x, int(1/step)*y))
			step /= 2
		}
		return uint16(res)
	}
}

var SimplexNoiseRedistribution = func(
	frequency,
	threshold float64,
	seed int64,
	octaves int,
	redistribution float64,
) algoFuncBasic {
	var simplexOctaves = SimplexNoiseOctaves(frequency, threshold, seed, octaves)
	return func(w, h, x, y int) uint16 {
		var result = float64(math.Pow(float64(simplexOctaves(w, h, x, y)), redistribution))
		if result > MAX_COLOR_VALUE {
			return MAX_COLOR_VALUE
		} else {
			return uint16(result)
		}
	}
}
