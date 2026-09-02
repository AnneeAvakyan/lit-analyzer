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

type occurrence struct {
	ChapterID     int
	SentenceIndex int
}

type AnalyzeBookUsecase struct {
	bookRepo      repository.BookRepository
	characterRepo repository.CharacterRepository
	chapterRepo   repository.ChapterRepository
	mentionRepo   repository.MentionRepository
}

func NewAnalyzeBookUsecase(bookRepo repository.BookRepository, characterRepo repository.CharacterRepository,
	chapterRepo repository.ChapterRepository, mentionRepo repository.MentionRepository) *AnalyzeBookUsecase {
	return &AnalyzeBookUsecase{
		bookRepo:      bookRepo,
		characterRepo: characterRepo,
		chapterRepo:   chapterRepo,
		mentionRepo:   mentionRepo,
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

	chapters := textproc.SplitIntoChapters(text)
	chapterEntities := []entities.Chapter{}

	for i := range chapters {
		chapter := entities.Chapter{
			BookID: bookID,
			Index:  i + 1,
			Text:   chapters[i],
		}
		chapterEntities = append(chapterEntities, chapter)
	}
	ids, err := uc.chapterRepo.CreateBatch(ctx, chapterEntities)
	if err != nil {
		return fmt.Errorf("create chapter: %w", err)
	}

	allOccurrences := make(map[string][]occurrence)

	for i, chapterText := range chapters {
		chapterID := ids[i]

		sentences := textproc.Segment(chapterText)
		candidatesInChapter := charextraction.ExtractCandidatesOccurrences(sentences)

		for candidate, sentenceNumbers := range candidatesInChapter {
			for _, sentenceIndex := range sentenceNumbers {
				occ := occurrence{
					ChapterID:     chapterID,
					SentenceIndex: sentenceIndex,
				}
				allOccurrences[candidate] = append(allOccurrences[candidate], occ)
			}
		}
	}

	mentions := []entities.Mention{}
	for candidate, occurrences := range allOccurrences {
		if len(occurrences) < minCharacterFrequency {
			continue
		}
		character := &entities.Character{
			BookID:        bookID,
			CanonicalName: candidate,
		}
		id, err := uc.characterRepo.Create(ctx, character)
		if err != nil {
			return fmt.Errorf("create character %q: %w", candidate, err)
		}

		for _, newOccurrence := range occurrences {
			mention := entities.Mention{
				CharacterID:   id,
				ChapterID:     newOccurrence.ChapterID,
				Position:      0, // мы не брали позицию, мб уберу это вообще
				SentenceIndex: newOccurrence.SentenceIndex,
			}
			mentions = append(mentions, mention)
		}
	}
	if err := uc.mentionRepo.CreateBatch(ctx, mentions); err != nil {
		return fmt.Errorf("create mentions: %w", err)
	}

	book.Status = entities.BookStatusDone
	err = uc.bookRepo.UpdateStatus(ctx, book)
	if err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	return nil
}
