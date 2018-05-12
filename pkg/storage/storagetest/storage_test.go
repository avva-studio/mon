package storagetest

import "github.com/glynternet/accounting-rest/pkg/storage"

// This line ensures that a Storage pointer can be assigned to a storage.Storage
// variable and, therefore, ensures that *Storage satisfies the storage.Storage
// interface
var _ storage.Storage = &Storage{}
