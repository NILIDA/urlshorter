package storage

import "errors"

var ErrNotFound = errors.New("url not found")

type Storage interface{
	Save(origURL string) (string, error)
	Get(shortURL string) (string, error)
	Close() error
}