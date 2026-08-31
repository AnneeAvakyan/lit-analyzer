package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/AnneeAvakyan/litanalyzer/internal/charextraction"
	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
	"github.com/AnneeAvakyan/litanalyzer/internal/domain/repository"
	"github.com/AnneeAvakyan/litanalyzer/internal/textproc"
)

const minCharacterFrequency int = 1

type AnalyzeBookUsecase struct {
	bookRepo      repository.BookRepository
	characterRepo repository.CharacterRepository
}

func NewAnalyzeBookUsecase(bookRepo repository.BookRepository, characterRepo repository.CharacterRepository) *AnalyzeBookUsecase {
	return &AnalyzeBookUsecase{
		bookRepo:      bookRepo,
		characterRepo: characterRepo,
	}
}

func (uc *AnalyzeBookUsecase) AnalyzeBook(ctx context.Context, bookID int) error {
	book, err := uc.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("get book: %w", err)
	}

	fileBytes, err := os.ReadFile(book.RawTextPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	text := string(fileBytes)
	//log.Printf("DEBUG: text length = %d", len(text))

	sentences := textproc.Segment(text)
	//log.Printf("DEBUG: sentences count = %d", len(sentences))
	//for i, s := range sentences {
	//	log.Printf("DEBUG: sentence[%d] = %q", i, s)
	//}

	candidates := charextraction.ExtractCandidatesWithFrequency(sentences)
	//log.Printf("DEBUG: candidates = %v", candidates)

	filtered := charextraction.FilterByFrequency(candidates, minCharacterFrequency)
	//log.Printf("DEBUG: filtered = %v", filtered)

	for name := range filtered {
		character := &entities.Character{
			BookID:        bookID,
			CanonicalName: name,
		}

		if _, err := uc.characterRepo.Create(ctx, character); err != nil {
			return fmt.Errorf("create character %q: %w", name, err)
		}
	}

	book.Status = entities.BookStatusDone
	err = uc.bookRepo.UpdateStatus(ctx, book)
	if err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	return nil
}
