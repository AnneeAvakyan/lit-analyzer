CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    status TEXT NOT NULL DEFAULT 'pending', -- pending | processing | done | failed
    raw_text_path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE chapters (
    id SERIAL PRIMARY KEY,
    book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    index INTEGER NOT NULL,
    text TEXT NOT NULL
);

CREATE INDEX idx_chapters_book_id ON chapters(book_id);

CREATE TABLE characters (
    id SERIAL PRIMARY KEY,
    book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    canonical_name TEXT NOT NULL
);

CREATE INDEX idx_characters_book_id ON characters(book_id);

CREATE TABLE character_aliases (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    alias TEXT NOT NULL
);

CREATE INDEX idx_character_aliases_character_id ON character_aliases(character_id);

CREATE TABLE mentions (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    chapter_id INTEGER NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,       -- позиция символа в тексте главы
    sentence_index INTEGER NOT NULL  -- глобальный индекс предложения в книге
);

CREATE INDEX idx_mentions_character_id ON mentions(character_id);
CREATE INDEX idx_mentions_chapter_id ON mentions(chapter_id);
-- составной индекс под co-occurrence: выборка упоминаний, отсортированных по sentence_index
CREATE INDEX idx_mentions_sentence_idx ON mentions(sentence_index);

CREATE TABLE relationships (
    id SERIAL PRIMARY KEY,
    book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    character_a_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    character_b_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    weight INTEGER NOT NULL DEFAULT 0,
    UNIQUE (character_a_id, character_b_id)
);

CREATE INDEX idx_relationships_book_id ON relationships(book_id);