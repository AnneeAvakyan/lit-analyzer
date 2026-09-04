package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AnneeAvakyan/litanalyzer/internal/usecase"
	"github.com/go-chi/chi/v5"
)

const defaultWindowSize = 5

type GraphHandler struct {
	graphUsecase *usecase.GraphUsecase
}

type GraphNode struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GraphEdge struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Weight int `json:"weight"`
}

type GraphResponse struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

func NewGraphHandler(graphUsecase *usecase.GraphUsecase) *GraphHandler {
	return &GraphHandler{
		graphUsecase: graphUsecase,
	}
}

// BuildGraph godoc
// @Summary Build graph
// @Description Creates edges of co-occurrences graph by bookID and size of window
// @Tags books
// @Param bookID path int true "Book ID"
// @Param windowSize query int false "Window size in sentences (default: 5)"
// @Success 200
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal server error"
// @Router /books/{bookID}/graph [post]
func (h *GraphHandler) BuildGraph(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	windowSize := defaultWindowSize
	if ws := r.URL.Query().Get("windowSize"); ws != "" {
		parsed, err := strconv.Atoi(ws)
		if err != nil {
			http.Error(w, "invalid windowSize: "+err.Error(), http.StatusBadRequest)
			return
		}
		windowSize = parsed
	}

	if err := h.graphUsecase.BuildGraph(r.Context(), bookID, windowSize); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetGraph godoc
// @Summary Get graph by book ID
// @Description Returns a graph of co-occurrences by given bookID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal server error"
// @Router /books/{id}/graph [get]
func (h *GraphHandler) GetGraph(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	characters, relationships, err := h.graphUsecase.GetGraph(r.Context(), bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	graphNodes := make([]GraphNode, 0, len(characters))

	for _, character := range characters {
		node := GraphNode{
			ID:   character.ID,
			Name: character.CanonicalName,
		}
		graphNodes = append(graphNodes, node)
	}

	graphEdges := make([]GraphEdge, 0, len(relationships))

	for _, relationship := range relationships {
		edge := GraphEdge{
			Source: relationship.CharacterAID,
			Target: relationship.CharacterBID,
			Weight: relationship.Weight,
		}
		graphEdges = append(graphEdges, edge)
	}

	response := GraphResponse{
		Nodes: graphNodes,
		Edges: graphEdges,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
