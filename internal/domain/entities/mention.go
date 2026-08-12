package entities

type Mention struct {
	ID            int `json:"id"`
	CharacterID   int `json:"characterID"`
	ChapterID     int `json:"chapterID"`
	Position      int `json:"position"`
	SentenceIndex int `json:"sentenceIndex"`
}
