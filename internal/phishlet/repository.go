package phishlet

import (
	"github.com/kgretzky/evilginx2/internal/storage"
)

type PhishletRepository struct {
	storage storage.Interface
	path    string
}

func NewPhishletRepository(storage storage.Interface, path string) *PhishletRepository {
	return &PhishletRepository{
		storage: storage,
		path:    path,
	}
}
