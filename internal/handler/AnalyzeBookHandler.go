package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/repository"
	"github.com/AnneeAvakyan/litanalyzer/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type AnalyzeBookHandler struct {
	uc *usecase.AnalyzeBookUsecase
}

func NewAnalyzeBookHandler(uc *usecase.AnalyzeBookUsecase) *AnalyzeBookHandler {
	return &AnalyzeBookHandler{uc: uc}
}

// AnalyzeBook godoc
// @Summary Analyze book
// @Description Extracts characters from book and filters them
// @Tags books
// @Param bookID path int true "Book ID"
// @Success 200
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /books/{bookID}/analyze [post]
func (h *AnalyzeBookHandler) AnalyzeBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.uc.AnalyzeBook(r.Context(), bookID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
