package internal

import "github.com/glynternet/go-accounting-storage"

var NewStorage StorageFunc

type StorageFunc func() (storage.Storage, error)
