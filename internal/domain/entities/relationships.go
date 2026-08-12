package entities

type Relationship struct {
	ID           int `json:"id"`
	BookID       int `json:"bookId"`
	CharacterAID int `json:"characterAId"`
	CharacterBID int `json:"characterBId"`
	Weight       int `json:"weight"`
}
