package postgres

import (
	"context"
	"fmt"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
	"github.com/AnneeAvakyan/litanalyzer/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCharacterRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresCharacterRepository(pool *pgxpool.Pool) *PostgresCharacterRepository {
	return &PostgresCharacterRepository{pool}
}

func (r *PostgresCharacterRepository) Create(ctx context.Context, character *entities.Character) (int, error) {
	query := `INSERT INTO characters (book_id, canonical_name) values ($1, $2) RETURNING id;`

	var id int
	err := r.pool.QueryRow(ctx, query, character.BookID, character.CanonicalName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert character: %w", err)
	}

	return id, nil
}

func (r *PostgresCharacterRepository) ListByBookID(ctx context.Context, id int) ([]entities.Character, error) {
	query := `SELECT id, book_id, canonical_name FROM characters WHERE book_id = $1;`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("select characters: %w", err)
	}

	defer rows.Close()

	var characters []entities.Character
	for rows.Next() {
		var c entities.Character
		if err := rows.Scan(&c.ID, &c.BookID, &c.CanonicalName); err != nil {
			return nil, fmt.Errorf("select characters: %w", err)
		}

		characters = append(characters, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("select characters: %w", err)
	}

	return characters, nil
}

func (r *PostgresCharacterRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM characters WHERE id = $1;`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete character: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
