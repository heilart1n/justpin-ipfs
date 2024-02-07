package pinners

import "io"

type Pinner interface {
	Name() string
	PinFile(fp string) (Result, error)
	PinWithReader(rd io.Reader) (Result, error)
	PinWithBytes(buf []byte) (Result, error)
	PinHash(hash string) (bool, error)
	PinDir(name string) (Result, error)
	Pin(path interface{}) (Result, error)
}

type Result interface {
	GetHash() string
	GetLink() string
}
