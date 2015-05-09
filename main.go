package main

import (
	"image/png"
	"log"
	"os"

	"unpublished/anytotarga/targa"
)

func main() {
	image, err := png.Decode(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	err = targa.Encode(os.Stdout, image)
	if err != nil {
		log.Fatal(err)
	}
}
