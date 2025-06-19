package user

import (
	"encoding/json"
	"database/sql"
)

type SignupResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func HandleSignup(req map[string]interface{}, db *sql.DB) []byte {
	// 테스트용 응답
	resp := SignupResponse{
		Success: true,
	}
	data, _ := json.Marshal(resp)
	return data
}
