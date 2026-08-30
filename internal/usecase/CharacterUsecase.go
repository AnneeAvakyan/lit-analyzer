package usecase

import (
	"context"
	"fmt"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
	"github.com/AnneeAvakyan/litanalyzer/internal/domain/repository"
)

type CharacterUsecase struct {
	characterRepo repository.CharacterRepository
}

func NewCharacterUsecase(characterRepo repository.CharacterRepository) *CharacterUsecase {
	return &CharacterUsecase{
		characterRepo: characterRepo,
	}
}

func (uc *CharacterUsecase) ListByBookID(ctx context.Context, bookID int) ([]entities.Character, error) {
	characters, err := uc.characterRepo.ListByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("usecase/CharacterUsecase.ListByBookID: %w", err)
	}
	return characters, nil
}
