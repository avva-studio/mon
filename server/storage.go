package server

import "github.com/glynternet/go-accounting-storage"

type StorageFunc func() (storage.Storage, error)
