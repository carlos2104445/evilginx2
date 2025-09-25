package phishlet

import (
	"github.com/kgretzky/evilginx2/internal/storage"
)

type SimplePhishletRepository struct {
	storage storage.Interface
	path    string
}

func NewSimplePhishletRepository(storage storage.Interface, path string) *SimplePhishletRepository {
	return &SimplePhishletRepository{
		storage: storage,
		path:    path,
	}
}
