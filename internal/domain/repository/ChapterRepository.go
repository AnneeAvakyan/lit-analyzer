package repository

import (
	"context"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

type ChapterRepository interface {
	CreateBatch(ctx context.Context, chapters []entities.Chapter) ([]int, error)
	ListByBookID(ctx context.Context, bookID int) ([]entities.Chapter, error)
	Delete(ctx context.Context, id int) error
}
