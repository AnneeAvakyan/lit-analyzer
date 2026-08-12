package entities

type Chapter struct {
	ID     int    `json:"id"`
	BookID int    `json:"bookId"`
	Index  int    `json:"index"`
	Text   string `json:"text"`
}
