package hendlererrors

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
