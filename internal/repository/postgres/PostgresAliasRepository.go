package postgres

import (
	"context"
	"fmt"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
	"github.com/AnneeAvakyan/litanalyzer/internal/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAliasRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAliasRepository(pool *pgxpool.Pool) *PostgresAliasRepository {
	return &PostgresAliasRepository{pool}
}

func (r *PostgresAliasRepository) Create(ctx context.Context, alias *entities.CharacterAlias) (int, error) {
	query := `INSERT INTO character_aliases (character_id, alias) VALUES ($1, $2) RETURNING id;`

	var id int
	err := r.pool.QueryRow(ctx, query, alias.CharacterID, alias.Alias).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert character alias: %w", err)
	}

	return id, nil
}

func (r *PostgresAliasRepository) ListByCharacterID(ctx context.Context, characterID int) ([]entities.CharacterAlias, error) {
	query := `SELECT id, character_id, alias FROM character_aliases WHERE character_id = $1 ORDER BY id ASC;`

	rows, err := r.pool.Query(ctx, query, characterID)
	if err != nil {
		return nil, fmt.Errorf("list character aliases: %w", err)
	}
	defer rows.Close()

	aliases := []entities.CharacterAlias{}
	for rows.Next() {
		var alias entities.CharacterAlias
		if err := rows.Scan(&alias.ID, &alias.CharacterID, &alias.Alias); err != nil {
			return nil, fmt.Errorf("scan character alias: %w", err)
		}
		aliases = append(aliases, alias)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return aliases, nil
}

func (r *PostgresAliasRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM character_aliases WHERE id = $1;`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete character alias: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
