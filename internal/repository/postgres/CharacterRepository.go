package postgres

import (
	"context"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

type CharacterRepository interface {
	Create(ctx context.Context, book *entities.Book) (int, error)
	GetById(ctx context.Context, id int) (*entities.Book, error)
	UpdateStatus(ctx context.Context, book *entities.Book) error
	Delete(ctx context.Context, id int) error
}
