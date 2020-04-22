package ogg

type WriterFlusher interface {
	Flush()
	Write([]byte)
}
