package storage

import "io"

type Provider interface {
	Store(name string, body io.Reader) (string, error)
	Load()
}
