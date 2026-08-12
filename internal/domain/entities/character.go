package entities

type Character struct {
	ID            int    `json:"id"`
	BookID        int    `json:"book_id"`
	CanonicalName string `json:"canonicalName"`
}
