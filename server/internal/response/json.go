package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func jsonResponse(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	return encoder.Encode(payload)
}

func Error(w http.ResponseWriter, status int, message string, err any) error {
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
	return jsonResponse(w, status, response)
}

func Success(w http.ResponseWriter, status int, message string, data any) error {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	return jsonResponse(w, status, response)
}

type PaginationMeta struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int32 `json:"total_pages"`
}

type PaginatedData struct {
	Items any            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

func NewPaginatedData(items any, page, limit int32, total int64) PaginatedData {
	totalPages := int32((total + int64(limit) - 1) / int64(limit))
	return PaginatedData{
		Items: items,
		Meta: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
