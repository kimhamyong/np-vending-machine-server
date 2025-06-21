package user

import (
	"database/sql"
	"encoding/json"
	"log"
	userRepo "vending-system/internal/repository/user"
	userService "vending-system/internal/service/user"
)

type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func HandleChangePassword(req map[string]interface{}, db *sql.DB) []byte {
	log.Printf("[ChangePassword] 요청 값: %+v\n", req)

	userid, ok1 := req["user_id"].(string)
	oldPassword, ok2 := req["old_password"].(string)
	newPassword, ok3 := req["new_password"].(string)

	log.Printf("user_id: %v (%T), old_password: %v (%T), new_password: %v (%T)",
	userid, req["user_id"], oldPassword, req["old_password"], newPassword, req["new_password"])


	if !ok1 || !ok2 || !ok3 {
		resp := ChangePasswordResponse{Success: false, Error: "필드 누락 또는 형식 오류"}
		data, _ := json.Marshal(resp)
		return data
	}

	repo := userRepo.NewRepository(db)
	err := userService.ChangePassword(repo, userid, oldPassword, newPassword)

	if err != nil {
		resp := ChangePasswordResponse{Success: false, Error: err.Error()}
		data, _ := json.Marshal(resp)
		return data
	}

	resp := ChangePasswordResponse{Success: true}
	data, _ := json.Marshal(resp)
	return data
}
