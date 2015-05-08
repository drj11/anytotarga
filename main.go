package main

import (
	"encoding/binary"
	"image"
	"image/png"
	"io"
	"log"
	"os"
)

type TARGAHeader struct {
	idLen             byte
	colorMapType      byte
	imageTypeCode     byte
	colorMapOrigin    uint16
	colorMapLength    uint16
	colorMapEntrySize byte
	xOrigin           uint16
	yOrigin           uint16
	width             uint16
	height            uint16
	imagePixelSize    byte
	imageDescriptor   byte
}

func main() {
	image, err := png.Decode(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	err = ToTARGA(os.Stdout, image)
	if err != nil {
		log.Fatal(err)
	}
}

func ToTARGA(w io.Writer, image image.Image) error {
	W := image.Bounds().Max.X
	H := image.Bounds().Max.Y

	header := TARGAHeader{0, 0, 2,
		0, 0, 24,
		0, 0, uint16(W), uint16(H), 24, 1 << 5}
	err := binary.Write(os.Stdout, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	for y := 0; y < W; y++ {
		for x := 0; x < H; x++ {
			r, g, b, _ := image.At(x, y).RGBA()
			rb := byte(float64(r)/257 + 0.5)
			gb := byte(float64(g)/257 + 0.5)
			bb := byte(float64(b)/257 + 0.5)
			_, err = os.Stdout.Write([]byte{bb, gb, rb})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
