package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

type PostgresRelationshipRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRelationshipRepository(pool *pgxpool.Pool) *PostgresRelationshipRepository {
	return &PostgresRelationshipRepository{pool: pool}
}

func (r *PostgresRelationshipRepository) CreateEdges(ctx context.Context, relationships []entities.Relationship) ([]int, error) {
	if len(relationships) == 0 {
		return nil, nil
	}

	valueStrings := make([]string, 0, len(relationships))
	args := make([]interface{}, 0, len(relationships)*4)

	argIdx := 1
	for _, rel := range relationships {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", argIdx, argIdx+1, argIdx+2, argIdx+3))
		args = append(args, rel.BookID, rel.CharacterAID, rel.CharacterBID, rel.Weight)
		argIdx += 4
	}

	query := fmt.Sprintf(`
		INSERT INTO relationships (book_id, character_a_id, character_b_id, weight)
		VALUES %s
		ON CONFLICT (character_a_id, character_b_id)
		DO UPDATE SET weight = EXCLUDED.weight
		RETURNING id`,
		strings.Join(valueStrings, ", "),
	)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("upsert relationships batch: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan relationship id: %w", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return ids, nil
}

func (r *PostgresRelationshipRepository) ListByBookID(ctx context.Context, bookID int) ([]entities.Relationship, error) {
	query := `
		SELECT id, book_id, character_a_id, character_b_id, weight
		FROM relationships
		WHERE book_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("list relationships by book: %w", err)
	}
	defer rows.Close()

	var relationships []entities.Relationship
	for rows.Next() {
		var rel entities.Relationship
		if err := rows.Scan(&rel.ID, &rel.BookID, &rel.CharacterAID, &rel.CharacterBID, &rel.Weight); err != nil {
			return nil, fmt.Errorf("scan relationship: %w", err)
		}
		relationships = append(relationships, rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return relationships, nil
}

func (r *PostgresRelationshipRepository) GetByCharacterID(ctx context.Context, characterID int) ([]entities.Relationship, error) {
	query := `
		SELECT id, book_id, character_a_id, character_b_id, weight
		FROM relationships
		WHERE character_a_id = $1 OR character_b_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, characterID)
	if err != nil {
		return nil, fmt.Errorf("get relationships by character: %w", err)
	}
	defer rows.Close()

	var relationships []entities.Relationship
	for rows.Next() {
		var rel entities.Relationship
		if err := rows.Scan(&rel.ID, &rel.BookID, &rel.CharacterAID, &rel.CharacterBID, &rel.Weight); err != nil {
			return nil, fmt.Errorf("scan relationship: %w", err)
		}
		relationships = append(relationships, rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return relationships, nil
}

func (r *PostgresRelationshipRepository) DeleteByBookID(ctx context.Context, bookID int) error {
	query := `DELETE FROM relationships WHERE book_id = $1;`
	_, err := r.pool.Exec(ctx, query, bookID)
	if err != nil {
		return fmt.Errorf("delete relationship: %w", err)
	}

	return nil
}
