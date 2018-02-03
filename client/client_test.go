package client

import "github.com/glynternet/go-accounting-storage"

// ensure that a Client can be used as a storage.Storage
var _ storage.Storage = Client("")
