package entities

type Mention struct {
	ID                  int `json:"id"`
	CharacterID         int `json:"characterID"`
	ChapterID           int `json:"chapterID"`
	GlobalSentenceIndex int `json:"global_sentence_index"`
	SentenceIndex       int `json:"sentenceIndex"`
}
