package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
	"github.com/AnneeAvakyan/litanalyzer/internal/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresChapterRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresChapterRepository(pool *pgxpool.Pool) *PostgresChapterRepository {
	return &PostgresChapterRepository{pool: pool}
}

func (r *PostgresChapterRepository) CreateBatch(ctx context.Context, chapters []entities.Chapter) ([]int, error) {
	if len(chapters) == 0 {
		return nil, nil // нечего вставлять
	}

	valueStrings := make([]string, 0, len(chapters))
	args := make([]interface{}, 0, len(chapters)*3)

	argIdx := 1
	for _, ch := range chapters {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", argIdx, argIdx+1, argIdx+2))
		args = append(args, ch.BookID, ch.Index, ch.Text)
		argIdx += 3
	}

	query := fmt.Sprintf(
		"INSERT INTO chapters (book_id, index, text) VALUES %s RETURNING id",
		strings.Join(valueStrings, ", "),
	)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("insert chapters batch: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan chapter id: %w", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return ids, nil
}

func (r *PostgresChapterRepository) ListByBookID(ctx context.Context, id int) ([]entities.Chapter, error) {
	query := `SELECT id, book_id, index, text FROM chapters WHERE book_id = $1 ORDER BY index;`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("select chapters by book_id: %w", err)
	}

	defer rows.Close()
	var chapters []entities.Chapter
	for rows.Next() {
		var chapter entities.Chapter
		if err := rows.Scan(&chapter.ID,
			&chapter.BookID, &chapter.Index, &chapter.Text); err != nil {
			return nil, fmt.Errorf("scan chapters by book_id: %w", err)
		}

		chapters = append(chapters, chapter)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return chapters, nil
}

func (r *PostgresChapterRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM chapters WHERE id = $1;`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete chapters: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
