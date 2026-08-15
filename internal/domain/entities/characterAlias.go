package entities

type CharacterAlias struct {
	ID          int    `json:"id"`
	CharacterID int    `json:"character_id"`
	Alias       string `json:"alias"`
}
