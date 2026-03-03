package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type PaginatedResponse struct {
	Data       any   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func ParseIDParam(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid argument: \"ID\"", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func Encode(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(v) != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
func GetPaginationValues(r *http.Request) (limit, offset int) {
	page := min(parseIntOrDefault(r.URL.Query().Get("page"), 1), 10000)
	limit = min(parseIntOrDefault(r.URL.Query().Get("limit"), 20), 100)
	offset = (page - 1) * limit
	return
}
func parseIntOrDefault(value string, def int) int {
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 {
		return def
	}
	return n
}
