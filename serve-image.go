package main

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"fmt"
	"time"
)

func InitServer() {
	var generator = GenerateImage
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc(
		"/VerticalGradient",
		generateImageResponse(
			generator,
			ConvertToGray16AlgoFunc(VerticalGradient),
		),
	)
	http.HandleFunc(
		"/HorizontalGradient",
		generateImageResponse(
			generator,
			ConvertToGray16AlgoFunc(HorizontalGradient),
		),
	)
	http.HandleFunc(
		"/CornerGradient",
		generateImageResponse(
			generator,
			ConvertToGray16AlgoFunc(CornerGradient),
		),
	)
	http.HandleFunc(
		"/CryptoRandom",
		generateImageResponse(
			generator,
			ConvertToGray16AlgoFunc(CryptoRandom),
		),
	)
	http.HandleFunc(
		"/CryptoRandomThreshold",
		generateImageResponse(
			generator,
			ConvertToGray16AlgoFunc(CryptoRandomThreshold(0.5)),
		),
	)
	http.HandleFunc(
		"/SimplexNoise",
		generateImageResponse(
			generator,
			ConvertToGray16AlgoFunc(SimplexNoise(0.1, 0.5, time.Now().UnixNano())),
		),
	)
	http.HandleFunc(
		"/SimplexNoiseOctaves",
		generateImageResponse(
			generator,
			ConvertToGray16AlgoFunc(SimplexNoiseOctaves(0.01, -1, time.Now().UnixNano(), 8)),
		),
	)

	err := http.ListenAndServe(":2017", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func generateImageResponse(
	generator func(algo algoFunc) string,
	algoFunc algoFunc,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("serve-image.go: %v\n", r.RequestURI)
		w.Header().Set("Content-Type", "image/png")
		content, err := base64.StdEncoding.DecodeString(generator(algoFunc))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		w.Write(content)
		return
	}
}
