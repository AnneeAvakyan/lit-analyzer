package repository

import (
	"context"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

type RelationshipRepository interface {
	CreateEdges(ctx context.Context, relationships []entities.Relationship) ([]int, error)
	ListByBookID(ctx context.Context, bookID int) ([]entities.Relationship, error)
	GetByCharacterID(ctx context.Context, id int) ([]entities.Relationship, error)
	DeleteByBookID(ctx context.Context, id int) error
}
