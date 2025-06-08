package repository

import (
	"database/sql"

	"wallpaper-system/internal/domain/repository"
)

// RepositoryRegistry содержит все репозитории приложения
type RepositoryRegistry struct {
	Product  repository.ProductRepository
	Material repository.MaterialRepository
}

// NewRepositoryRegistry создает новый реестр репозиториев
func NewRepositoryRegistry(db *sql.DB) *RepositoryRegistry {
	return &RepositoryRegistry{
		Product:  NewProductRepository(db),
		Material: NewMaterialRepository(db),
	}
} 