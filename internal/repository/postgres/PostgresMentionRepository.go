package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

type PostgresMentionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresMentionRepository(pool *pgxpool.Pool) *PostgresMentionRepository {
	return &PostgresMentionRepository{pool: pool}
}

func (r *PostgresMentionRepository) CreateBatch(ctx context.Context, mentions []entities.Mention) error {
	if len(mentions) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(mentions))
	args := make([]interface{}, 0, len(mentions)*4)

	argIdx := 1
	for _, m := range mentions {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", argIdx, argIdx+1, argIdx+2, argIdx+3))
		args = append(args, m.CharacterID, m.ChapterID, m.GlobalSentenceIndex, m.SentenceIndex)
		argIdx += 4
	}

	query := fmt.Sprintf(
		"INSERT INTO mentions (character_id, chapter_id, global_sentence_index, sentence_index) VALUES %s",
		strings.Join(valueStrings, ", "),
	)

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert mentions batch: %w", err)
	}

	return nil
}

func (r *PostgresMentionRepository) ListByBookID(ctx context.Context, bookID int) ([]entities.Mention, error) {
	// mentions не хранит book_id напрямую — join через chapters
	query := `
		SELECT m.id, m.character_id, m.chapter_id, m.global_sentence_index, m.sentence_index
		FROM mentions m
		JOIN chapters c ON c.id = m.chapter_id
		WHERE c.book_id = $1
		ORDER BY m.sentence_index ASC;
	`

	rows, err := r.pool.Query(ctx, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("list mentions: %w", err)
	}
	defer rows.Close()

	var mentions []entities.Mention
	for rows.Next() {
		var m entities.Mention
		if err := rows.Scan(&m.ID, &m.CharacterID, &m.ChapterID, &m.GlobalSentenceIndex, &m.SentenceIndex); err != nil {
			return nil, fmt.Errorf("scan mention: %w", err)
		}
		mentions = append(mentions, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return mentions, nil
}

func (r *PostgresMentionRepository) ReassignCharacter(ctx context.Context, oldCharacterID, newCharacterID int) error {
	query := `UPDATE mentions SET character_id = $1 WHERE character_id = $2;`

	_, err := r.pool.Exec(ctx, query, newCharacterID, oldCharacterID)
	if err != nil {
		return fmt.Errorf("reassign mentions: %w", err)
	}

	return nil
}
