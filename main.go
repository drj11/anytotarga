package main

import (
	"encoding/binary"
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
	header := TARGAHeader{0, 0, 2,
		0, 0, 24,
		0, 0, 16, 16, 24, 1<<5}
	err := binary.Write(os.Stdout, binary.LittleEndian, header)
	for i := 0; i < 8; i++ {
		for j := 0; j < 16; j++ {
			os.Stdout.Write([]byte{255, 0, 0})
		}
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 16; j++ {
			os.Stdout.Write([]byte{0, 255, 255})
		}
	}
	if err != nil {
		log.Print("Error: ", err)
	}
}
