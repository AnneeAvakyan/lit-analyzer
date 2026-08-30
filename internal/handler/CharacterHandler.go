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

func NewCharacterHandler(characterUsecase *usecase.CharacterUsecase) *CharacterHandler {
	return &CharacterHandler{characterUsecase: characterUsecase}
}

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

func (h *CharacterHandler) MergeCharacters(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TargetID  int   `json:"targetId"`
		SourceIDs []int `json:"sourceIds"`
	}

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
