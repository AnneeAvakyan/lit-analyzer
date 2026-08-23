package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
	"github.com/AnneeAvakyan/litanalyzer/internal/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresBookRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresBookRepository(pool *pgxpool.Pool) *PostgresBookRepository {
	return &PostgresBookRepository{pool: pool}
}

func (r *PostgresBookRepository) Create(ctx context.Context, book *entities.Book) (int, error) {
	query := `
	INSERT INTO books (title, author, status, raw_text_path, created_at)
	VALUES ($1, $2, $3, $4, now())
	RETURNING id, created_at;
	`

	var id int
	err := r.pool.QueryRow(ctx, query,
		book.Title,
		book.Author,
		book.Status,
		book.RawTextPath,
	).Scan(&id, &book.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("insert book: %w", err)
	}

	return id, nil
}

func (r *PostgresBookRepository) GetByID(ctx context.Context, id int) (*entities.Book, error) {
	query := `
SELECT id, title, author, status, raw_text_path, created_at
FROM books
WHERE id = $1;`

	var book entities.Book
	err := r.pool.QueryRow(ctx, query, id).Scan(&book.ID,
		&book.Title,
		&book.Author,
		&book.Status,
		&book.RawTextPath,
		&book.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get book: %w", err)
	}

	return &book, nil
}

func (r *PostgresBookRepository) UpdateStatus(ctx context.Context, book *entities.Book) error {
	query := `UPDATE books SET status = $1 WHERE id = $2;`

	tag, err := r.pool.Exec(ctx, query, book.Status, book.ID)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *PostgresBookRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1;`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
