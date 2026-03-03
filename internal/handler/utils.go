package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
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

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 72 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}
