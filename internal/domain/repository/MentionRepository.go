package repository

import (
	"context"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

type MentionRepository interface {
	ListByBookID(ctx context.Context, bookID int) ([]entities.Mention, error)
	CreateBatch(ctx context.Context, mentions []entities.Mention) error
	ReassignCharacter(ctx context.Context, oldCharacterID, newCharacterID int) error
}
