package repository

import (
	"context"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

type CharacterRepository interface {
	Create(ctx context.Context, character *entities.Character) (int, error)
	ListByBookId(ctx context.Context, id int) ([]entities.Character, error)
	Delete(ctx context.Context, id int) error
}
