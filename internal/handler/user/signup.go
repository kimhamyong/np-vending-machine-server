package user

import (
    "database/sql"
    "encoding/json"
    userRepo "vending-system/internal/repository/user"
    userService "vending-system/internal/service/user"
	"log"
)

type SignupResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error,omitempty"`
}

func HandleSignup(req map[string]interface{}, db *sql.DB) []byte {
	log.Printf("[Signup] 요청 값: %+v\n", req)

    userid, ok1 := req["userid"].(string)
    password, ok2 := req["password"].(string)
    if !ok1 || !ok2 {
        resp := SignupResponse{Success: false, Error: "필드 누락 또는 형식 오류"}
        data, _ := json.Marshal(resp)
        return data
    }

    repo := userRepo.NewRepository(db)
    err := userService.Signup(repo, userid, password)

    if err != nil {
        resp := SignupResponse{Success: false, Error: err.Error()}
        data, _ := json.Marshal(resp)
        return data
    }

    resp := SignupResponse{Success: true}
    data, _ := json.Marshal(resp)
    return data
}
