package main

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/drj11/anytotarga/targa"
)

func main() {
	image, _, err := image.Decode(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	err = targa.Encode(os.Stdout, image)
	if err != nil {
		log.Fatal(err)
	}
}
