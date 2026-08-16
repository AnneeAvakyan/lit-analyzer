package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AnneeAvakyan/litanalyzer/internal/usecase"
)

type BookHandler struct {
	bookUsecase *usecase.BookUsecase
}

func NewBookHandler(bookUsecase *usecase.BookUsecase) *BookHandler {
	return &BookHandler{bookUsecase: bookUsecase}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	// ограничиваем размер тела запроса, чтобы не словить OOM на гигантском файле
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // 20 MB

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "invalid multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	book, err := h.bookUsecase.CreateBook(r.Context(), title, author, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}
