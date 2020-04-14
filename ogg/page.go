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
		_ = binary.Read(bytes.NewBuffer(p.raw[0:27]), byteOrder, p.header)
	}
	return p.header
}

func (p Page) Type() byte {
	return 0
}

func (p Page) Serial() uint32 {
	return 0
}

func (p Page) Granule() int64 {
	return 0
}
