package ogg

/*

https://en.wikipedia.org/wiki/Ogg

Large part of this code came from https://github.com/mccoyst/ogg © 2016 Steve McCoy

*/
import (
	"bytes"
	"encoding/binary"
)

var byteOrder = binary.LittleEndian

const (
	Continuation = uint8(1)
	BOS          = uint8(2)
	EOS          = uint8(4)
)

type pageHeader struct {
	OggS          [4]byte // 0-3, always == "OggS"
	StreamVersion byte    // 4, always == 0
	HeaderType    byte    // 5
	Granule       int64   // 6-13, codec-specific
	Serial        uint32  // 14-17, associated with a logical stream
	Page          uint32  // 18-21, sequence number of page in packet
	Crc           uint32  // 22-25
	Nsegs         byte    // 26
}

type Page struct {
	Raw    []byte
	header *pageHeader
}

func NewPage(raw []byte) *Page {
	return &Page{
		Raw: raw,
	}
}

func (p *Page) Header() *pageHeader {
	if p.header == nil {
		var h pageHeader
		err := binary.Read(bytes.NewBuffer(p.Raw[0:27]), byteOrder, &h)
		if err != nil {
			panic(err)
		}
		p.header = &h
	}
	return p.header
}
