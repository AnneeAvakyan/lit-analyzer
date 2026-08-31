package usecase

import (
	"context"
	"fmt"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
	"github.com/AnneeAvakyan/litanalyzer/internal/domain/repository"
)

type CharacterUsecase struct {
	characterRepo repository.CharacterRepository
	aliasRepo     repository.AliasRepository
	mentionRepo   repository.MentionRepository
}

func NewCharacterUsecase(characterRepo repository.CharacterRepository,
	aliasRepo repository.AliasRepository, mentionRepo repository.MentionRepository) *CharacterUsecase {
	return &CharacterUsecase{
		characterRepo: characterRepo,
		aliasRepo:     aliasRepo,
		mentionRepo:   mentionRepo,
	}
}

func (uc *CharacterUsecase) ListByBookID(ctx context.Context, bookID int) ([]entities.Character, error) {
	characters, err := uc.characterRepo.ListByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("usecase/CharacterUsecase.ListByBookID: %w", err)
	}
	return characters, nil
}

func (uc *CharacterUsecase) MergeCharacters(ctx context.Context, targetID int, sourceIDs []int) error {
	for _, sourceID := range sourceIDs {
		ch, err := uc.characterRepo.GetByID(ctx, sourceID)
		if err != nil {
			return fmt.Errorf("usecase/CharacterUsecase.GetByID: %w", err)
		}

		alias := &entities.CharacterAlias{
			CharacterID: targetID,
			Alias:       ch.CanonicalName,
		}

		if _, err := uc.aliasRepo.Create(ctx, alias); err != nil {
			return fmt.Errorf("usecase/CharacterUsecase.Create: %w", err)
		}

		if err := uc.mentionRepo.ReassignCharacter(ctx, ch.ID, targetID); err != nil {
			return fmt.Errorf("usecase/CharacterUsecase.ReassignCharacter: %w", err)
		}

		if err := uc.characterRepo.Delete(ctx, sourceID); err != nil {
			return fmt.Errorf("usecase/CharacterUsecase.Delete: %w", err)
		}
	}
	return nil
}
