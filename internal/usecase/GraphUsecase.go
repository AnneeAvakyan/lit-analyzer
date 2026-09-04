package usecase

import (
	"context"
	"fmt"

	"github.com/AnneeAvakyan/litanalyzer/internal/cooccurence"
	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/repository"
)

type GraphUsecase struct {
	mentionRepo      repository.MentionRepository
	relationshipRepo repository.RelationshipRepository
	characterRepo    repository.CharacterRepository
}

func NewGraphUsecase(mentionRepo repository.MentionRepository,
	relationshipRepo repository.RelationshipRepository,
	characterRepo repository.CharacterRepository) *GraphUsecase {
	return &GraphUsecase{
		mentionRepo:      mentionRepo,
		relationshipRepo: relationshipRepo,
		characterRepo:    characterRepo,
	}
}

func (g *GraphUsecase) BuildGraph(ctx context.Context, bookID int, windowSize int) error {
	mentions, err := g.mentionRepo.ListByBookID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("list by book id: %w", err)
	}

	if err := g.relationshipRepo.DeleteByBookID(ctx, bookID); err != nil {
		return fmt.Errorf("delete by book id: %w", err)
	}

	graph := cooccurence.BuildGraph(mentions, windowSize)

	relationships := []entities.Relationship{}
	for pair, weight := range graph {
		rel := entities.Relationship{
			BookID:       bookID,
			CharacterAID: pair[0],
			CharacterBID: pair[1],
			Weight:       weight,
		}
		relationships = append(relationships, rel)
	}

	if _, err := g.relationshipRepo.CreateEdges(ctx, relationships); err != nil {
		return fmt.Errorf("create relationships: %w", err)
	}

	return nil
}

func (g *GraphUsecase) GetGraph(ctx context.Context, bookID int) ([]entities.Character, []entities.Relationship, error) {
	characters, err := g.characterRepo.ListByBookID(ctx, bookID)
	if err != nil {
		return nil, nil, fmt.Errorf("list characters: %w", err)
	}

	relationships, err := g.relationshipRepo.ListByBookID(ctx, bookID)
	if err != nil {
		return nil, nil, fmt.Errorf("list relationships: %w", err)
	}

	return characters, relationships, nil
}
