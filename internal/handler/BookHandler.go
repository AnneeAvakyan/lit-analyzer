package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
	"github.com/AnneeAvakyan/litanalyzer/internal/domain/repository"
	"github.com/AnneeAvakyan/litanalyzer/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type UpdateBookStatusRequest struct {
	Status string `json:"status"`
}

type BookHandler struct {
	bookUsecase *usecase.BookUsecase
}

func NewBookHandler(bookUsecase *usecase.BookUsecase) *BookHandler {
	return &BookHandler{bookUsecase: bookUsecase}
}

// CreateBook godoc
// @Summary Create a new book
// @Description Upload a book file for analysis
// @Tags books
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Book title"
// @Param author formData string false "Book author"
// @Param file formData file true "Book text file"
// @Success 201 {object} entities.Book
// @Failure 400 {string} string "bad request"
// @Router /books [post]
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

// GetBookByID godoc
// @Summary Get book by ID
// @Description Returns a book by given ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} entities.Book
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "book not found"
// @Failure 500 {string} string "internal server error"
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id: "+err.Error(), http.StatusBadRequest)
		return
	}
	book, err := h.bookUsecase.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// UpdateBook godoc
// @Summary Update book by ID
// @Description Update a book by given ID
// @Tags books
// @Accept json
// @Produce json
// @Param request body UpdateBookStatusRequest true "New status"
// @Param id path int true "Book ID"
// @Success 200
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /books/{id} [patch]
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req UpdateBookStatusRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedBook := &entities.Book{ID: id, Status: req.Status}
	err = h.bookUsecase.UpdateBook(r.Context(), updatedBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

// DeleteBook godoc
// @Summary Delete book by ID
// @Description Deletes a book by given ID
// @Tags books
// @Param id path int true "Book ID"
// @Success 204
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "book not found"
// @Failure 500 {string} string "internal server error"
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.bookUsecase.DeleteBook(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
