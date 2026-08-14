package repository

import (
	"context"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

type AliasRepository interface {
	Create(ctx context.Context, alias *entities.CharactersAlias) (int, error)
	ListByCharacterId(ctx context.Context, id int) ([]entities.CharactersAlias, error)
	Delete(ctx context.Context, id int) error
}
