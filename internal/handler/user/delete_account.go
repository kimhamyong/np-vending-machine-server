package user

import (
	"database/sql"
	"encoding/json"
	"log"
	userRepo "vending-system/internal/repository/user"
	userService "vending-system/internal/service/user"
)

type DeleteAccountResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func HandleDeleteAccount(req map[string]interface{}, db *sql.DB) []byte {
	log.Printf("[DeleteAccount] 요청 값: %+v\n", req)

	userid, ok1 := req["user_id"].(string)
	password, ok2 := req["password"].(string)
	if !ok1 || !ok2 {
		resp := DeleteAccountResponse{Success: false, Error: "필드 누락 또는 형식 오류"}
		data, _ := json.Marshal(resp)
		return data
	}

	repo := userRepo.NewRepository(db)
	err := userService.DeleteAccount(repo, userid, password)
	if err != nil {
		resp := DeleteAccountResponse{Success: false, Error: err.Error()}
		data, _ := json.Marshal(resp)
		return data
	}

	resp := DeleteAccountResponse{Success: true}
	data, _ := json.Marshal(resp)
	return data
}
