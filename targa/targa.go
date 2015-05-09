package targa

import (
	"encoding/binary"
	"image"
	"io"
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

type Packet struct {
	w   io.Writer
	H   byte
	buf []byte
}

func NewPacket(w io.Writer) Packet {
	var p Packet
	p.w = w
	return p
}

func (p *Packet) Add(b, g, r byte) error {
	if p.H&0x7f == 127 {
		err := p.Flush()
		if err != nil {
			return err
		}
	}
	if len(p.buf) == 0 {
		p.H = 0
		p.buf = append(p.buf, b, g, r)
		return nil
	}
	if p.H&128 == 0 {
		lb := p.buf[len(p.buf)-3]
		lg := p.buf[len(p.buf)-2]
		lr := p.buf[len(p.buf)-1]
		if b == lb && g == lg && r == lr {
			// Flush and convert to RLE packet.
			p.buf = p.buf[:len(p.buf)-3]
			err := p.Flush()
			if err != nil {
				return err
			}
			p.H = 129
			p.buf = []byte{b, g, r}
			return nil
		}
		p.H += 1
		p.buf = append(p.buf, b, g, r)
	} else {
		lb := p.buf[len(p.buf)-3]
		lg := p.buf[len(p.buf)-2]
		lr := p.buf[len(p.buf)-1]
		if b == lb && g == lg && r == lr {
			p.H += 1
			return nil
		}
		// Convert to raw packet.
		err := p.Flush()
		if err != nil {
			return err
		}
		p.H += 1
		p.buf = append(p.buf, b, g, r)
	}
	return nil
}

func (p *Packet) Flush() error {
	if len(p.buf) == 0 {
		return nil
	}
	_, err := p.w.Write([]byte{p.H})
	if err != nil {
		return err
	}
	_, err = p.w.Write(p.buf)
	if err != nil {
		return err
	}
	p.H = 0
	p.buf = []byte{}
	return nil
}

func EncodeUncompressed(w io.Writer, image image.Image) error {
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

func Encode(w io.Writer, image image.Image) error {
	W := image.Bounds().Max.X
	H := image.Bounds().Max.Y

	header := TARGAHeader{0, 0, 10,
		0, 0, 24,
		0, 0, uint16(W), uint16(H), 24, 1 << 5}
	err := binary.Write(os.Stdout, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	packet := NewPacket(os.Stdout)

	for y := 0; y < W; y++ {
		for x := 0; x < H; x++ {
			r, g, b, _ := image.At(x, y).RGBA()
			rb := byte(float64(r)/257 + 0.5)
			gb := byte(float64(g)/257 + 0.5)
			bb := byte(float64(b)/257 + 0.5)
			err = packet.Add(bb, gb, rb)
			if err != nil {
				return err
			}
		}
	}
	err = packet.Flush()
	return err
}
