package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AnneeAvakyan/litanalyzer/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type CharacterHandler struct {
	characterUsecase *usecase.CharacterUsecase
}

type MergeCharactersRequest struct {
	TargetID  int   `json:"targetId"`
	SourceIDs []int `json:"sourceIds"`
}

func NewCharacterHandler(characterUsecase *usecase.CharacterUsecase) *CharacterHandler {
	return &CharacterHandler{characterUsecase: characterUsecase}
}

// ListByBookID godoc
// @Summary Get characters list by book ID
// @Description Returns a characters list by given book ID
// @Tags characters
// @Produce json
// @Param bookID path int true "Book ID"
// @Success 200 {object} []entities.Character
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /books/{bookID}/characters [get]
func (h *CharacterHandler) ListByBookID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "bookID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.characterUsecase.ListByBookID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// MergeCharacters godoc
// @Summary Merge duplicate characters
// @Description Merges source characters into a target character, preserving their names as aliases and reassigning their mentions
// @Tags characters
// @Accept json
// @Param request body MergeCharactersRequest true "Target and source character IDs"
// @Success 200
// @Failure 400 {string} string "invalid request body"
// @Failure 500 {string} string "internal server error"
// @Router /characters/merge [post]
func (h *CharacterHandler) MergeCharacters(w http.ResponseWriter, r *http.Request) {
	var req MergeCharactersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.characterUsecase.MergeCharacters(r.Context(), req.TargetID, req.SourceIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
