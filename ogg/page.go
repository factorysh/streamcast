package ogg

import (
	"bytes"
	"encoding/binary"
)

var byteOrder = binary.LittleEndian

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
	raw    []byte
	header *pageHeader
}

func NewPage(raw []byte) *Page {
	return &Page{
		raw: raw,
	}
}

func (p *Page) Header() *pageHeader {
	if p.header == nil {
		var h pageHeader
		err := binary.Read(bytes.NewBuffer(p.raw[0:27]), byteOrder, &h)
		if err != nil {
			panic(err)
		}
		p.header = &h
	}
	return p.header
}
