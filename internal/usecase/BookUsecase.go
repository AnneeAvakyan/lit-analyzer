package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
	"github.com/AnneeAvakyan/litanalyzer/internal/repository"
	"github.com/google/uuid"
)

type BookUsecase struct {
	bookRepo    repository.BookRepository
	storagePath string // куда на диске сохраняем загруженные тексты
}

func NewBookUsecase(bookRepo repository.BookRepository, storagePath string) *BookUsecase {
	return &BookUsecase{
		bookRepo:    bookRepo,
		storagePath: storagePath,
	}
}

func (uc *BookUsecase) CreateBook(ctx context.Context, title, author string, fileContent io.Reader) (*entities.Book, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	filename := fmt.Sprintf("%s.txt", uuid.NewString())
	fullPath := filepath.Join(uc.storagePath, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, fileContent); err != nil {
		return nil, fmt.Errorf("write file: %w", err)
	}

	book := &entities.Book{
		Title:       title,
		Author:      author,
		Status:      "pending",
		RawTextPath: fullPath,
	}

	id, err := uc.bookRepo.Create(ctx, book)
	if err != nil {
		return nil, fmt.Errorf("save book: %w", err)
	}
	book.ID = id

	return book, nil
}
